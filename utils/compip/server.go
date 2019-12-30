package compip

import (
	"bufio"
	"os"

	"github.com/372572571/Exercise/utils/message"
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
	Err     interface{} // 异常携带数据
}
type servicePoint interface {
	handle(*OnMsg, *Service)
}

// Service 连接管道服务
type Service struct {
	fun   map[string]interface{} // 功能方法存放容器
	queue *message.Queue
	// pool chan
}

// NewService .... 创建服务
func NewService() *Service {
	var s = &Service{}
	s.fun = make(map[string]interface{})
	s.queue = message.NewQueue(1)
	s.queue.Run()
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
		s.Output(-1, UNKNOWN, m, "invalid json data.")
		return
	}
	if m.Service == "" {
		s.Output(-1, UNKNOWN, m, "invalid service data.")
		return
	}
	f, ok := s.fun[m.Service]
	if !ok {
		s.Output(-1, UNKNOWN, m, "invalid service.")
		return
	}
	// 强制转换
	switch f.(type) {
	case ServiceHandle:
		f.(ServiceHandle)(m, s) // 执行
		break
	default:
		s.Output(-1, m.Service, m, "service error,invalid handle.")
		return
	}
}

func (s *Service) echo(code int8, name string, data interface{}, err interface{}) {
	var info = &sendMsg{}
	info.Code = code
	info.Data = data
	info.Service = name
	res, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(info)
	s.queue.Push(res)
}

// Output 向pip输出数据
func (s *Service) Output(code int8, name string, data interface{}, err interface{}) {
	s.echo(code, name, data, err)
}
