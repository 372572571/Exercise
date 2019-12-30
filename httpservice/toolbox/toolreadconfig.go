package toolbox

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)
// 读取配置文件信息操作

// ASCII 码表
const (
	COMMENTS=byte(35) // 对应#
	WARP=byte(10) // 对应换行键
)

// 错误处理
func errnil(err error){
	if err!=nil{log.Println(err)}
	if err!=nil{panic(err)}
}

// 获得配置文件的设置 返回键值对
func OpenConfig(filename string)map[string]string{
	configinfo,err:=os.Open(filename)//open file get *File
	errnil(err)
	fileinfo,err:=configinfo.Stat()
	errnil(err)
	filesize:=fileinfo.Size()
	filebuff:=make([]byte,filesize)
	configinfo.Read(filebuff)
	configinfo.Close()
	return ConfigRead(filebuff)

}

//判断有多少行
func lines (iobyte []byte) int{
	lens:=0
	ascii:=byte(10)
	for _,v:=range iobyte{
		if ascii==v {
			lens++
		}
	}
	return lens
}

// 返回
// 解析过的K和V MAP[string][string]
func ConfigRead(iobyte []byte)map[string]string{
	flag:=0 // 记录有效参数开始的位子
	iolen:=len(iobyte)
	m:=make(map[string]string)
	for i:=0;i<iolen;i++{
		if isComments(iobyte[i]) {
			for i<iolen&&noWarp(iobyte[i]){
				i++ // 如果是一个注释符，那么丢弃其中#到\n数据
			}
		}else{
			flag=i
			for i<iolen&&noWarp(iobyte[i]) {
				i++
			}
			m[parse(string(iobyte[flag:i]))[0]]=parse(string(iobyte[flag:i]))[1]
		}

	}
	return m
}

// 拆分字符串返回map[string]string
func parse(str string)[]string{
	str=strings.TrimSpace(str)
	if strings.Index(str,"=")!=-1 {
		return strings.Split(str,"=")
	}
	return []string{"0","0"}
}

// 比较一个byte是不是注释符
func isComments(b byte)bool{
	return b==COMMENTS
}

// 是否换行,否为真
func noWarp(b byte)bool{
	return b!=WARP
}

// 修改配置，filename 配置文件名称 需要修改的配置键值对
func ChangeSetting(filename string ,s map[string]string)[]string{
	m:=OpenConfig(filename)
	m=SettingMatch(m,s)
	return mapToString(m)

}

// 遍历s中的键，m中是否存在，存在则用s的值覆盖m的值
func SettingMatch(m,s map[string]string) map[string]string{
	for k,v:=range s {
		if _,ok:=m[k];ok {
			m[k]=v
		}
	}
	return  m
}

// m的键值对拼接成数组字符串
func mapToString(m map[string]string)[]string{
	flag:=0
	lens:=len(m)
	s:=make([]string,lens)
	for k,v:=range m{
		s[flag]=k+"="+v
		flag++
	}
	return s
}

// 判断源配置文件是否存在
func IsConfigSource(programpath string)bool{
	_,err:=os.Stat(programpath)
	if err==nil{
		return true
	}
	return false
}
// 字符串数组加\n后字符长度
func strArrGetByteSize(s []string)int{
	bytelen:=0
	for _,v:=range s{
		bytelen+=len([]byte(v))
	}
	return bytelen+len(s)
}
// 输出配置文件
// filename 文件输出的路径和文件名 不存在则自动创建，如果存在则覆盖
// s []string 配置数据
func OutConfig(filename string,s []string){
	data:=make([]byte,strArrGetByteSize(s))
	flag:=0
	for _,v:=range s{
		d:=[]byte(v)
		for _,v1:=range d{
			data[flag]=v1
			flag++
		}
		data[flag]='\n'
		flag++
	}
	err:=ioutil.WriteFile(filename,data,0644)
	if err!=nil {
		log.Println(err)
	}

}