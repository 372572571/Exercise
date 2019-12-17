package compip

import (
	"bufio"
	"fmt"
	"os"
	jsoniter "github.com/json-iterator/go"
)

const (
	// UNKNOWN 未知的命令
	UNKNOWN = "unknown"
)
// ServiceHandle ...... 处理函数
type ServiceHandle func(*OnMsg, *Service)

// OnMsg 统一接收消息的格式
type OnMsg struct {
	Service string                 // 请求的服务
	Data    map[string]interface{} // 携带的数据
}

// sendMsg
type sendMsg struct {
	Service string // 请求的服务
	Code    int8
	Data    interface{} // 携带的数据
}
type servicePoint interface {
	handle(*OnMsg, *Service)
}

// Service 连接管道服务
type Service struct {
	fun map[string]interface{} // 功能方法存放容器
	// pool chan
}

// NewService .... 创建服务
func NewService() *Service {
	var s = &Service{}
	s.fun = make(map[string]interface{})
	return s
}

// Open 开启服务
func (s *Service) Open() {
	// 接收管道消息
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		go s.Unmarshal(input.Text())
	}
}

// Registered ..... 注册服务
func (s *Service) Registered(key string, handle ServiceHandle) {
	s.fun[key] = handle
}

// Unmarshal ..... 解析传入数据
func (s *Service) Unmarshal(data string) {
	var m = &OnMsg{}
	err := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(data, m)
	if err != nil {
		s.echo(-1, UNKNOWN, "invalid json data.")
		return
	}
	if m.Service == "" {
		s.echo(-1, UNKNOWN, "invalid service data.")
		return
	}
	f, ok := s.fun[m.Service]
	if !ok {
		s.echo(-1, UNKNOWN, "invalid service.")
		return
	}
	// 强制转换
	switch f.(type) {
	case ServiceHandle:
		f.(ServiceHandle)(m, s) // 执行
		break
	default:
		s.echo(-1, m.Service, "service error,invalid handle.")
		return
	}
}

func (s *Service) echo(code int8, name string, data interface{}) {
	var info = &sendMsg{}
	info.Code = code
	info.Data = data
	info.Service = name
	res, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(info)
	fmt.Println(string(res))
}

// Output 向pip输出数据
func (s *Service) Output(code int8, name string, data interface{}) {
	s.echo(code, name, data)
}
