package httpservice

import (
	"html/template"
	"log"
	"net/http"
	"reflect"
)

type IApp interface {
	Init(ctx *Context) //初始化context
	W() http.ResponseWriter
	R() *http.Request
	Echo(htmlpath ...string)
}

type App struct {
	ctx  *Context               //存放 请求信息
	Data map[string]interface{} //信息数据
}

// 初始化
func (a *App) Init(ctx *Context) {
	a.ctx = ctx
	a.Data = make(map[string]interface{})
}

// 获得 ResponseWriter
func (a *App) W() http.ResponseWriter {
	return a.ctx.w
}

// 获取Request
func (a *App) R() *http.Request {
	return a.ctx.r
}

// 解析模板并输出
func (a *App) Echo(htmlpath ...string) {
	if htmlpath[0] != "" {
		t, err := template.ParseFiles(htmlpath[0]) //解析模板
		a.appError(err)
		t.Execute(a.W(), a.Data)
		return
	}
	http.NotFound(a.W(), a.R())
}

// 构建一个空的*App结构
func newApp() *App {
	return &App{}
}

func (a *App) appError(err error) {
	if err != nil {
		log.Println(err)
	}
	// http.NotFound(a.W(), a.R())
}

// 返回App定义的函数
func appLimit() (limitfunc map[string]int) {
	limitfunc = make(map[string]int)
	app := newApp()
	appt := reflect.TypeOf(app)
	for i := 0; i < appt.NumMethod(); i++ {
		limitfunc[appt.Method(i).Name] = 0
	}
	return
}

/*
//Display
//tpls ...string 多参数 数组形式访问
func (a *App) Display(tpls ...string) {
	if len(tpls) == 0 {
		return
	}
	//tpls[0] path Base返回 分割path返回最后一个路径元素
	name := filepath.Base(tpls[0])
	//Must函数用于包装返回(*Template, error)的函数/方法调用，它会在err非nil时panic，一般用于变量初始化
	//ParseFiles 创建，解析模板。
	t := template.Must(template.ParseFiles(tpls...))
	t.ExecuteTemplate(a.W(), name, a.Data)
}

func (a *App) DisplayWithFuncs(funcs template.FuncMap, tpls ...string) {
	if len(tpls) == 0 {
		return
	}

	name := filepath.Base(tpls[0])
	t := template.Must(template.New(name).Funcs(funcs).ParseFiles(tpls...))
	t.ExecuteTemplate(a.W(), name, a.Data)
}*/
