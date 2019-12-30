package chanrpc

// CallInfo 一个回调结构
type CallInfo struct {
	id      interface{}   // 需要调用方法的标识符
	args    []interface{} // 携带的参数
	result  chan *Result  //结构存放结构
	call    interface{}   // 回调函数
	isAsync bool          // 是否需要异步
}

// Result 用于存放调用服务后,返回的数据
type Result struct {
	data []interface{} // 存放服务返回数据
	Err  error         // 错误信息存放
	call interface{}   // 回调函数
}

// Get 根据index 获取数据
func (r Result) Get(index int) interface{} {
	if index >= 0 && (len(r.data)-1) >= index {
		return r.data[index]
	}
	return nil
}
