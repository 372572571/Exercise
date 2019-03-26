package wsserver

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// UserConns 验证通过的链接（注册用户）
type UserConns = map[int]interface{}

// TempUserConns 临时用户 （没用通过验证的）
type TempUserConns = map[*websocket.Conn]interface{}

// 未验证的链接（没有验证的用户）

// 3. 实现一个表示客户链接的结构（对象）（并发安全）
// a. 添加线程安全（并发安全）
// b. 存放一个客户长链接 *conn
// c. 写入管道 chan []byte
// d. 最大信息长度限制
// 链接关闭标识符

// WsConn ... 用于表示一个用户链接的结构(并发安全)
type WsConn struct {
	sync.Mutex
	conn             *websocket.Conn // 链接指针
	write            chan []byte     // 写入管道
	maxMessageLength int64           // 信息最大长度
	isClose          bool            // 是否关闭(关闭为true)
	isRead           bool            // 是否已经开启线程读管道(true 为开启)
	isWrite          bool            // 是否已经开始线程写管道(写管道是否就绪)(true 为开启)
}

// NewWsConn 包装websocket链接
func NewWsConn(chanWriteLength int, conn *websocket.Conn, maxMessageLength int64) *WsConn {
	var wsc = new(WsConn) // 包装
	wsc.conn = conn
	wsc.write = make(chan []byte, chanWriteLength)
	wsc.maxMessageLength = maxMessageLength
	go func() { // 开启一个携程监听写管道内容
		fmt.Printf("开始监听写入")
		for item := range wsc.write { // 循环管道内容（如果没有内容阻塞）
			if item == nil {
				break // 如果读取到空结束循环
			}

			err := conn.WriteMessage(websocket.BinaryMessage, item) // 二进制写入数据 websocket.BinaryMessage

			if err != nil { // 如果写入信息错误
				break
			}
		}
		conn.Close()        // 如果循环监听写入管道因为失败被打断尝试关闭链接
		wsc.Lock()          // 同步锁
		wsc.isClose = true  // 链接废弃
		wsc.isWrite = false //写管道监听结束
		wsc.Unlock()        // 解锁
	}()
	return wsc
}

// RunWrite ...(线程阻塞循环)
// 实现创建携程监听链接写入管道,有数据就发送给客户 发生意外后关闭链接
func (wsc *WsConn) runWrite() bool {
	wsc.Lock()       // 同步锁
	if wsc.isWrite { // 如果写入管道已经被开启过
		return false
	}
	wsc.isWrite = true // 表示写入监听已开启
	wsc.Unlock()       // 解锁
	return true
}

// ReadMsg ... 读取客户发来的信息
func (wsc *WsConn) ReadMsg() ([]byte, error) {
	_, b, err := wsc.conn.ReadMessage()
	return b, err
}

// Close ...
func (wsc *WsConn) Close() {
	wsc.conn.Close()
	// wsc.Lock()
	// wsc.isClose = true
	// wsc.Unlock()
}
