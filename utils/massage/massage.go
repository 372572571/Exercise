package massage

import (
	jsoniter "github.com/json-iterator/go"
)

// Massage ...
type Massage struct {
	Code int
	Key  string
	Type string
	Msg  string
	Data interface{}
}

// CreateMassage 创建信息结构
func CreateMassage() *Massage {
	var res = new(Massage)
	return res
}

// DefError 默认错误信息
func (msg *Massage) DefError() *Massage {
	msg.Code = 0
	msg.Type = "ERROR"
	msg.Msg = "defult error info"
	return msg
}

// SetError 传递错误信息
func (msg *Massage) SetError(info string) *Massage {
	msg.Code = 0
	msg.Type = "ERROR"
	msg.Msg = info
	return msg
}

// SetSuccess 传递成功信息
func (msg *Massage) SetSuccess(info string, data interface{}) *Massage {
	msg.Code = 1
	msg.Type = "SUCCESS"
	msg.Msg = info
	msg.Data = data
	return msg
}

// ToString 结构转字符串
func (msg *Massage) ToString() string {
	res, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(msg)
	if err != nil {
		return "{}"
	}
	return string(res)
}

func (msg *Massage) SetCode(code int) *Massage {
	msg.Code = code
	return msg
}

func (msg *Massage) SetKey(Key string) *Massage {
	msg.Key = Key
	return msg
}

func (msg *Massage) SetData(data interface{}) *Massage {
	msg.Data = data
	return msg
}

func (msg *Massage) SetMsg(m string) *Massage {
	msg.Msg = m
	return msg
}
