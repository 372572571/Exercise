package chanrpc

import "fmt"

// Server 一个服务类型(需求消费者)
type Server struct {
	isStart  bool                        // 是否已经启动(Run后)
	function map[interface{}]interface{} // 存放注册方法
	CallInfo chan *CallInfo              // 存放调用信息
}

// NewServer rpc服务
// rpcLen 服务管道长度
func NewServer(rpcLen int) *Server {
	server := new(Server)                               // 创建一个服务对象
	server.isStart = false                              // 默认没有启动服务
	server.function = make(map[interface{}]interface{}) // 初始化服务对象方法容器
	server.CallInfo = make(chan *CallInfo, rpcLen)      // 初始化回调管道长度(一次能够放存放的数据条数)
	return server                                       // 返回一个服务对象
}

// Registered ...注册函数
func (server *Server) Registered(id, function interface{}) bool {
	// 注册时进行类型鉴定
	switch function.(type) {
	case func([]interface{}) error:
		server.function[id] = function // 注册服务方法有参数 返回error
	case func([]interface{}) ([]interface{}, error):
		server.function[id] = function // 注册服务方法有参数  返回 []interface{} error
	default:
		fmt.Println("注册类型错误")
		return false // 没有对应类型返回false
	}
	return true
}

// callSelf 调用自己的方法并调用,调用结果包装成
func (server *Server) callSelf(info *CallInfo) {
	defer func(info *CallInfo) {
		if err := recover(); err != nil { // 错误处理 崩溃处理 ()
			PipeReturnError(info, ErrNotFindFunc)
		}
	}(info)
	res := &Result{call: info.call} // 创建结果信息结构
	f := server.function[info.id]   // 根据id查找注册的方法
	switch f.(type) {               // 类型推断并执行
	case func([]interface{}) error:
		res.Err = f.(func([]interface{}) error)(info.args)
	case func([]interface{}) ([]interface{}, error):
		res.data, res.Err = f.(func([]interface{}) ([]interface{}, error))(info.args)
	default:
		res.Err = ErrNotFindFunc // 没有对应类型输出对应错误,并处理
	}
	// 判断同步流程，异步流程
	if info.isAsync {
		select {
		case info.result <- res: // 异步分支
		default: // 意外流程处理
			res = &Result{call: info.call}
			res.Err = ErrAsyncPush
			info.result <- res
		}
	} else {
		info.result <- res // 结果结构管道中存入数据,交付客户处理(同步流程)
	}
}

// Run ....开始处理 CallInfo (处理业务)
func (server *Server) Run() {
	// 防止重复启动
	if server.isStart {
		return
	}
	server.isStart = true
	go func() {
		for {
			data := <-server.CallInfo
			go func(info *CallInfo) {
				server.callSelf(info)
			}(data)
		}
	}()
}

// FastCallBack 快速调用服务（回调模式调用）
func (server *Server) FastCallBack(id interface{}, args []interface{}, callback func(i map[string]interface{})) {
	fun := server.function[id]
	switch fun.(type) {
	case func([]interface{}, func(i map[string]interface{})): // 正常类型调用
		fun.(func([]interface{}, func(i map[string]interface{})))(args, callback)
	default: // 类型错误处理
		callback(nil)
	}
}

// Fast ... 快速调用同步返回
func (server *Server) Fast(id interface{}, args []interface{}) *Result {
	client := NewClient(1)
	go func() {
		client.Link(server)
		client.CallRequest(id, args)
	}()

	res := <-client.Result
	fmt.Println(res)
	return res
}
