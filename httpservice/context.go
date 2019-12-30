package httpservice

import (
	"net/http"
)

type IContext interface {
	Config(w http.ResponseWriter,r *http.Request)
}

// 请求信息的封装
type Context struct {
	w http.ResponseWriter	// 响应
	r *http.Request		// 请求
}

// 数据储存到结构
func (c *Context) Config(w http.ResponseWriter,r *http.Request){
	c.w=w
	c.r=r
}