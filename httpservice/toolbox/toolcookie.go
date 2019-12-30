package toolbox

import "net/http"

// 设置cookie判断cookie存在与否.

// 设置一个临时的Cookit.
// 浏览器关闭后失效.
func SetTempCookie(k , v string,w http.ResponseWriter){
	// Path 如果不设置会导致只能在设置这个Cookie的控制器下才能使用
	// 导致不能全站使用
	cookie:=http.Cookie{Name: k, Value: v,Path:"/"}
	http.SetCookie(w,&cookie)
}

// 判断是否有这个cookie.Nema
// 存在返回true.
func IsCookieStr(r *http.Request,str string)bool{
	cookies:=r.Cookies()
	for _,cookie:=range cookies{
		if cookie.Name==str {
			return true
		}
	}
	return false
}

// 获得一个cookie.Name的cookie.Value
// 如果不存在返回空字符串.
func GetCookieStr(r *http.Request,str string)string{
	cookies:=r.Cookies()
	for _,cookie:=range cookies{
		if cookie.Name==str {
			return string(cookie.Value)
		}
	}
	return ""
}