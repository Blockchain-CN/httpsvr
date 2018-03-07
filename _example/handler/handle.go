// Copyright 2018 Blockchain-CN . All rights reserved.
// https://github.com/Blockchain-CN

package handler

import (
	"time"

	"github.com/Blockchain-CN/httpsvr/_example/idl"
)

func DemoHandle(req *idl.DemoReq) *idl.DemoResp {

	resp := idl.NewDemoResp()
	if req.Age <= 18 {
		// 测试下panic recover
		panic(1)
		resp.Errno = 0
		resp.Msg = "Success"
		resp.Result = "免费"
		return resp
	}
	time.Sleep(time.Second * 10)
	resp.Errno = 0
	resp.Msg = "Success"
	resp.Result = "付费"
	return resp
}
