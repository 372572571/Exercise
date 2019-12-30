package wsserver

// Agent ... 链接代理人（存放链接和其他客户数据代表着一个客户）
type Agent interface {
	GetID() int    // 获取用户代理人的id 如果为“0”说明是临时客户
	GetConn() *WsConn // 获取代理人的链接
	StartRead()       // 循环监听代理人发送的数据(如果是临时代理人可以做一些限制)
	OnClose()         // 关闭链接
}