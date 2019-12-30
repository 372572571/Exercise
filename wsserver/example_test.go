package wsserver_test

import (
	"fmt"
	"time"

	"github.com/372572571/Exercise/wsserver"
)

func Testwsserver() {
	ws := wsserver.NewWsServer("127.0.0.1:8091", time.Duration(5)*time.Second, "", "", wsserver.NewHandler(func(ws *wsserver.WsConn) wsserver.Agent {
		a := &TempConn{wsConn: ws}
		fmt.Println("一个链接代理创建成功")
		return a
	}))
	ws.Start()
	for {
	}
}

// TempConn 临时用户
type TempConn struct {
	wsConn *wsserver.WsConn // 用户链接
}

// GetID  ...如果当前是临时链接返回0
func (t *TempConn) GetID() int {
	_ = t
	return 0
}

// GetConn ... 获取当前代理人链接
func (t *TempConn) GetConn() *wsserver.WsConn {
	_ = t
	return t.wsConn
}

// StartRead 开始监听数据
func (t *TempConn) StartRead() {
	for {
		data, err := t.wsConn.ReadMsg()
		if err != nil {
			break
		}
		if data == nil {
			break
		}
		fmt.Println("read", data)
	}
}

// OnClose 关闭代理人
func (t *TempConn) OnClose() {
	t.wsConn.Close()
}
