package wsserver

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// 实现websocket服务

// 1. 实现一个 WebsocketServer 结构
// a. 服务的地址
// b. 链接超时时间
// c. Https 证书文件
// d. Https key文件
// e. 监听链接
// http handler

// WsServer ...
type WsServer struct {
	Addr        string        // 服务地址
	HTTPTimeout time.Duration // http 超时时间
	CertFile    string        // 证书文件
	KeyFile     string        // key文件
	listener    net.Listener  // 监听链接
	handler     *WsHandler    // ws处理方法
}

// Start ...
// 启动服务
func (server *WsServer) Start(handler *WsHandler) {
	server.handler = handler
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		fmt.Println("创建链接失败")
		return
	}
	// 默认值处理
	server.listener = ln
	server.setDefault(server.listener)

	var httpServer = &http.Server{
		Addr:           server.Addr,
		Handler:        server.handler,
		ReadTimeout:    server.HTTPTimeout,
		WriteTimeout:   server.HTTPTimeout,
		MaxHeaderBytes: 1024,
	}
	go httpServer.Serve(server.listener)
}

// 配置默认设置
func (server *WsServer) setDefault(ln net.Listener) {
	if server.handler == nil {
		panic("handler 不能为空.")
	}
	if server.handler.maxConnNum <= 0 {
		server.handler.maxConnNum = 100 // 最大链接数量
	}
	if server.handler.writeContentLength <= 0 {
		server.handler.writeContentLength = 4096 // 写入内容长度
	}
	if server.handler.maxMessageLength <= 0 {
		server.handler.maxMessageLength = 4096 // 最大消息长度
	}
	if server.HTTPTimeout <= 0 {
		server.HTTPTimeout = 10 * time.Second // 超时时间
	}
	if server.handler.NewAgent == nil {
		panic("server.handler.NewAgent 不能为空.")
	}
	// 如果有文件配置安全链接 https wss
	if server.CertFile != "" || server.KeyFile != "" {
		config := &tls.Config{}
		config.NextProtos = []string{"http/1.1"}
		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(server.CertFile, server.KeyFile)
		if err != nil {
			log.Fatal("%v", err)
		}
		ln = tls.NewListener(ln, config)
	}
}

// Close 关闭服务
func (server *WsServer) Close() {
	server.listener.Close()
	os.Exit(1)
}
