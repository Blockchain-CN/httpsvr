// Copyright 2018 Blockchain-CN . All rights reserved.
// https://github.com/Blockchain-CN

package httpsvr

import (
	"fmt"
)

func getErrMsg(err error) []byte {
	return []byte(fmt.Sprintf(`{"errno":-1,"errmsg":"%s"}`, err.Error()))
}
