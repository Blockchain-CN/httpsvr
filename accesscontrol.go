// Copyright 2018 Lothar . All rights reserved.
// https://github.com/Blockchain-CN

package httpsvr

import (
	"errors"
	"sync"
	"time"
	"sync/atomic"
)

// Access 接入控制
type Access struct {
	maxAccess int
	// 桶,用来限制最大并发数
	bucket chan struct{}
	// 用来判断何时桶内已空，否则需要循环判断
	wg sync.WaitGroup
	closed int32
}

// NewAccessor ...
func NewAccessor(i int) *Access {
	var ac Access
	ac.closed = 0
	ac.maxAccess = i
	// 10ms 准入超时计时器，时间窗口
	ac.bucket = make(chan struct{}, i)
	return &ac
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
