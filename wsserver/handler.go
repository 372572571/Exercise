package wsserver

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// 2. 实现一个http handler 结构（或说是类，需要实现httpserver方法）
// a. 限定最大链接数量
// b. 写入数据最大长度
// c. 最大信息长度
// d. 新链接回调函数（新的一个链接升级成conn之后的回调（函数））
// e. 升级链接第三方库实现
// 	i. websocket.Upgrader
// 		1) HandshakeTimeout 握手超时时间
// 		2) ReadBufferSize, WriteBufferSize 读写缓冲大小 （如果不定义则使用 默认 生成的大小）
// 		3) WriteBufferPool  写数据缓冲池 如果没有则用默认 生成的writebuffersize（或者自己设置的 writebuffersize）长度创建
// 		4) Subprotocols 子协议 字符串数组
// 		5) Error 异常函数
// 		6) CheckOrigin 检测源函数 可以自己实现 用于屏蔽某些来源请求
// 		7) EnableCompression 是否启用压缩（消息压缩）
// f. WebsocketConnSet 链接队列用于存放链接的map数据结构（websocket库实现的 *websocket.conn 当作key）
// g. mutexConns 同步锁 操作 WebsocketConnSet  中的数据时用于同步
// h. wg sync.WaitGroup 同步等待组

// WsHandler ...
type WsHandler struct {
	maxConnNum         int                 // 最大的链接数量
	writeContentLength int                 // 写入内容长度
	maxMessageLength   int64               // 消息最大长度
	mutexConns         sync.Mutex          // 同步锁
	upHTTPToConn       websocket.Upgrader  // http 升级到 conn(长链接)
	NewAgent           func(*WsConn) Agent // 根据包装的链接创建一个代理人
	tempUserConns      TempUserConns       // 临时客户队列
	userConns          UserConns           // 认证用户队列
}

// NewHandler ...
func NewHandler(agentfunc func(*WsConn) Agent) *WsHandler {
	var h = new(WsHandler)
	h.NewAgent = agentfunc
	h.tempUserConns = make(TempUserConns)
	h.userConns = make(UserConns)
	return h
}

// 实现接口
// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }
func (wsh *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // 如果是GET 请求直接拒绝
		return
	}
	con, err := wsh.upHTTPToConn.Upgrade(w, r, nil) // 升级链接
	if err != nil {
		return // 如果升级链接发生错误
	}
	if wsh.maxConnNum > len(wsh.tempUserConns) { // 判断当前最大链接数量是否超出
		con.Close()
		return
	}
	wsh.SetTempConn(con)                                                   // 把链接丢入零食用户列表
	var wsc = NewWsConn(wsh.writeContentLength, con, wsh.maxMessageLength) // 创建一个WsConn并创建一个线程监听写 chan
	var agent = wsh.NewAgent(wsc)                                          // 创建一个用户代理人
	agent.StartRead()                                                      // 开始阻塞

	//  阻塞结束 清理链接
	wsh.CleanAgent(agent)
}

// SetTempConn ...
// 把一个临时链接放入 临时用户列表
func (wsh *WsHandler) SetTempConn(con *websocket.Conn) {
	wsh.mutexConns.Lock()
	wsh.tempUserConns[con] = struct{}{}
	wsh.mutexConns.Unlock()
}

// CleanAgent ... 清理客户链接
func (wsh *WsHandler) CleanAgent(agent Agent) {
	wsh.mutexConns.Lock()
	if agent.GetID() == 0 { // 移除链接
		delete(wsh.tempUserConns, agent.GetConn().conn)
	} else {
		delete(wsh.userConns, agent.GetID())
	}
	wsh.mutexConns.Unlock()
	agent.OnClose()
}
