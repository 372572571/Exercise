package chanrpc

import "errors"

// NAME 包名
var NAME = "chanrpc"

// ErrNotFindFunc 找不到注册的方法
var ErrNotFindFunc = errors.New(NAME + ": Can't find registration function.")

// ErrAsyncPush 异步返回管道推入失败
var ErrAsyncPush = errors.New(NAME + ": Async chan push error.")

// ErrCallBackType 错误的回调类型
var ErrCallBackType = errors.New(NAME + ": error callBack type.")

// PipeReturnError 管道推入错误
func PipeReturnError(info *CallInfo, err error) {
	errRes := &Result{call: info.call}
	errRes.Err = err
	info.result <- errRes
}
