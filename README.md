# httpsvr
HttpServer.
简单的http框架，只实现了option、router、accesscontrol、idl功能
目前仅支持POSTJSON

## 使用
step1:创建httpsvr对象。
step2:添加路由。
step3:开始服务。
optional: 优雅退出。
``` go
func main() {
    // step1:创建httpsvr对象。
    s := httpsvr.New(
    "127.0.0.1:10024",
    httpsvr.SetReadTimeout(time.Millisecond*200),
    httpsvr.SetWriteTimeout(time.Millisecond*200),
    httpsvr.SetMaxAccess(2),
    )
    // 优雅退出
    go GracefulExit(s)
    // step2:添加路由。
    s.AddRoute("POST", "/test/api", &ctrls.DemoCtrl{})
    // step3:开始服务。
    s.Serve()
}

// GracefulExit 优雅退出
func GracefulExit(svr *httpsvr.Server) {
        sigc := make(chan os.Signal, 0)
        signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
        <-sigc
        println("closing agent...")
        svr.GracefulExit()
        println("agent closed.")
        os.Exit(0)
}
```

## 基础功能
### 路由框架 router
1、使用"github.com/julienschmidt/httprouter"作为核心路由。
2、使用原生ListenAndServe方法进行监听(内部调用net/http.router接口)。
3、通过实现ServeHTTP()方法来实现net/http.router接口，并在其中实现路由。
``` go
// AddRoute ...
func (s *Server) AddRoute(method, path string, ctrl IController) {
    s.router.Handle(method, path, handle)

}
```

``` go
// ServeHTTP implement net/http.router
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    s.router.ServeHTTP(w, req)
}
```

``` go
// http listen
    oriSvr *http.Server
    s.oriSvr = &http.Server{Addr: addr, Handler: s}
```

``` go
// Serve ...
func (s *Server) Serve() error {
    return s.oriSvr.ListenAndServe()
}
```

### 服务器选项 option
通过提供回调函数的形式来实现变参设置参数。
``` go
// ServerOption option赋值回调函数
type ServerOption func(o *option)
```

### 接口定义 idl
规定了IController接口，用户定义controller必须实现此接口。
``` go
// IController ...
type IController interface {
    // 此方法获得接口参数的结构体，框架会直接将req.body直接unmarshal到这里面
    GenIdl() interface{}
    // 实际处理的handle函数
    Do(interface{}) interface{}
}
```

### 流量控制／优雅退出 accesscontrol
- 流量控制
通过路由之后，在处理入口进行入口控制:
1、首先判断是否允许进入(标志为会在优雅退出时打开，平时关闭)。
2、再向chan中写入标志，如果chan已满(到达最大并发数)则无法写入，经过等待时间窗口(100ms)后如果还写不进去就退出。
3、请求执行完毕时从chan中读一个标志为，代表当前人数减1。

- 优雅退出
1、在允许进入后，每次放入一个信号，就WaitGroup.Add(1)，代表当前又有1个在执行。
2、请求执行完毕时，就WaitGroup.Done()，代表有一个请求执行完了。
3、退出时首先先将标志位用原子操作进行赋值，防止后续请求接着进入(参看流量控制.1)。
4、再等待所有请求执行完毕，用WaitGroup.Wait()

``` go
// InControl 入口控制
func (a *Access) InControl() error {
    if atomic.LoadInt32(&a.closed) == 1 {
        return errors.New("server is closing")
    }
    select {
    case a.bucket <- struct{}{}:
        a.wg.Add(1)
    case <-time.After(time.Millisecond * 100):
        return errors.New("server is busy please try later")
    }
    return nil
}

// OutControl 出口注销
func (a *Access) OutControl() {
    <-a.bucket
    a.wg.Done()
}

// Stop 优雅退出
func (a *Access) Stop() {
	if !atomic.CompareAndSwapInt32(&a.closed, 0, 1) {
		return
	}
	// 第一种判断桶内为空
	a.wg.Wait()
	/*第二种判断桶内为空
	for {
		if len(a.bucket)== 0 {
			return
		}
	}*/
}
```

## 测试

### 测试方案
1、路由测试：直接发送请求，遍历路由。
测试结果：PASS
2、并发控制测试：设置并发数量为1，并设置句柄中sleep(10s)，检测第二个请求是否可以进入。
测试结果：PASS
3、panic测试：设置并发数量为1，在处理句柄中Panic，检测后续请求是否可以进入。
测试结果：PASS
注意：这里一定要自己实现panic-recover，因为上层是无法在接入计数上减去1的，会导致"panic请求"WaitGroup.Add(1)，recover之后没有WaitGroup.Done()，因此WaitGroup.Wait()死锁，同理流量控制的chan
4、优雅退出：设置并发数>1 ,并设置句柄中sleep(10s)，第一个正常请求进入后，快速关闭server，测试第二个请求访问接口的返回值。
测试结果：PASS
由于server没有立刻关闭(第一个请求还在处理)，因此后续请求还是能访问，但是没有进行处理，直接返回了。

### 测试代码
见[使用]

### 测试请求
```
curl -H 'content-type: application/json' -X POST -d '{"Name":"luda","Age":111}' http://127.0.0.1:10024/test/api //正常
curl -H 'content-type: application/json' -X POST -d '{"Name":"luda","Age":11}' http://127.0.0.1:10024/test/api  //Panic
```
