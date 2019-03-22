package chanrpc

import (
	"fmt"
	"sync"
)

// Client 客户
type Client struct {
	server      *Server      // 请求的服务对象（向哪个服务请求）
	Result      chan *Result // 返回的结果结构
	AsyncResult chan *Result // 存放异步管道
	asyncWait   int          // 异步发起未处理数据条数
	asyncStart  bool         // 当前是否在处理异步管道
	mutex       *sync.Mutex  //互斥锁
}

// NewClient 创建一个客户(需求生产者)
// asyncLen 异步结果管道长度
func NewClient(asyncLen int) *Client {
	client := new(Client)
	client.Result = make(chan *Result, 1)
	client.AsyncResult = make(chan *Result, asyncLen)
	client.asyncWait = 0
	client.mutex = new(sync.Mutex)
	return client
}

// Link 链接服务
func (client *Client) Link(server *Server) {
	client.server = server
}

// callBack 调用绑定的回调信息
func (client *Client) callBack(res *Result) {
	switch res.call.(type) {
	case func(*Result):
		res.call.(func(*Result))(res)
	default:
		fmt.Println(ErrCallBackType)
		panic("回调类型错误")
	}
}

// Request 给服务发送请求(同步处理)
func (client *Client) Request(id interface{}, args []interface{}, callBack func(*Result)) {
	go func() {
		callInfo := &CallInfo{id: id, args: args, result: client.Result, call: callBack, isAsync: false}
		client.server.CallInfo <- callInfo
		res := <-client.Result
		client.callBack(res)
	}()
}

// AsyncRequest 异步
func (client *Client) AsyncRequest(id interface{}, args []interface{}, callBack func(*Result)) {
	go func() {
		callInfo := &CallInfo{id: id, args: args, result: client.AsyncResult, call: callBack, isAsync: true}
		select {
		case client.server.CallInfo <- callInfo:
			client.asyncWait++ // 期待返回的异步管道数量
		default:
			// 意外处理
			print("异步 意外处理")
		}
	}()
}

// AsyncRun 开始处理异步
func (client *Client) AsyncRun() {
	go func() {
		if client.asyncStart {
			return // 已经启动过
		}
		client.mutex.Lock()
		client.asyncStart = true
		client.mutex.Unlock()
		for client.asyncWait > 0 {
			client.callBack(<-client.AsyncResult)
			client.asyncWait--
		}
		client.asyncStart = false
	}()
}
