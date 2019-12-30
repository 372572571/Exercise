package httpservice

import (
	// "fmt"
	"net/http"
	"reflect"
	"strings"
	// "A.brikc/controller/common"
)

// webRoute 路由信息
var webRoute map[string]reflect.Type = make(map[string]reflect.Type)

// 添加路由
func router(pattern string, t reflect.Type) {
	webRoute[strings.ToLower(pattern)] = t
}

// Router 储存反射信息到路由
func Router(pattern string, app IApp) {
	refV := reflect.ValueOf(app)
	refT := reflect.Indirect(refV).Type()
	router(pattern, refT)
}

// AutoRouter 注册路由
func AutoRouter(app IApp) {
	refV := reflect.ValueOf(app)
	refT := reflect.Indirect(refV).Type()
	// 去掉controller 并小写
	refName := strings.TrimSuffix(strings.ToLower(refT.Name()), "controller")
	// fmt.Println(refName, refT)
	router(refName, refT) //储存
}

// Static 静态服务容器1
var Static map[string]string = make(map[string]string) // 静态服务

func serveStatic(w http.ResponseWriter, r *http.Request) bool {
	for prefix, static := range Static { // 循环注册静态服务
		if strings.HasPrefix(r.URL.Path, prefix) {
			if prefix != "/storage" {
				file := static + r.URL.Path[len(prefix):]
				// fmt.Println(file)
				http.ServeFile(w, r, file)
				return true
			}
			// 静态文件权限,登陆后才能调用
			/*if common.VerifyUser(r)&&prefix=="/storage"{
				// fmt.Println("注册")
				file := static + r.URL.Path[len(prefix):]
				http.ServeFile(w, r, file)
				return true
			}*/
		}
	}

	return false
}
