package httpservice

import (
	"net/http"
	"log"
)

// port 监听端口

func ServerRun_Http(port string){
	// http.Server 管理httpserver行为
	// 创建一个自定义的serve
	serve:=http.Server{
		Addr: port, // 端口
		Handler: newHandler(), // webapp/app/handler 实现
	}
	log.Fatal(serve.ListenAndServe())
}
