package pgmfx

import "fmt"

var ErrQueryPoolGet = fmt.Errorf("get free connection fail")
var ErrQueryCall = fmt.Errorf("call request fail")
var ErrRowsProcess = fmt.Errorf("rows process fail")

var ErrGetJsonFail = fmt.Errorf("fail to get json")
var ErrUnmarshalJsonFail = fmt.Errorf("fail json unmarshal")
