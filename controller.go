// Copyright 2018 Blockchain-CN . All rights reserved.
// https://github.com/Blockchain-CN

package httpsvr

// IController ...
type IController interface {
	GenIdl() interface{}
	Do(interface{}) interface{}
}
