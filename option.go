// Copyright 2018 Blockchain-CN . All rights reserved.
// https://github.com/Blockchain-CN

package httpsvr

import "time"

type option struct {
	maxAccess     int
	dumpResponse  bool
	enableElasped bool
	dumpAccess    bool
	validate      bool
	readTimeout   time.Duration
	writeTimeout  time.Duration
}

// ServerOption option赋值回调函数
type ServerOption func(o *option)

// SetReadTimeout 设置读超时
func SetReadTimeout(rt time.Duration) ServerOption {
	return func(o *option) {
		o.readTimeout = rt
	}
}

// SetWriteTimeout 设置写超时
func SetWriteTimeout(wt time.Duration) ServerOption {
	return func(o *option) {
		o.writeTimeout = wt
	}
}

// SetMaxAccess 最大接入数
func SetMaxAccess(i int) ServerOption {
	return func(o *option) {
		o.maxAccess = i
	}
}
