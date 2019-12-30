package returnjson

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//ImassageJson... 返回信息接口定义
type ImassageJson interface {
	Massage() []byte
	NewImassage(servertype int, massage string) ImassageJson
}

// json massage
type massageJson struct {
	ServerType    int // error int
	ServerMassage string
}

// init massage struct return *MassageJson
func NewMassage(servertype int, massage string) *massageJson {
	return &massageJson{ServerType: servertype, ServerMassage: massage}
}

// return ImassageJson
func (m *massageJson) NewImassage(servertype int, massage string) ImassageJson {
	return NewMassage(servertype, massage)
}

// return json massage
func (m *massageJson) Massage() []byte {
	s, err := json.Marshal(m)
	if err != nil {
		return NewMassage(500, err.Error()).Massage()
	}
	return s
}

// put json client
func PutJson(w http.ResponseWriter, imassageJson ImassageJson) {
	fmt.Fprint(w, string(imassageJson.Massage()))
}

// is error!= nil put Json error
func PutErrorJson(err error, w http.ResponseWriter) {
	if err != nil {
		imassageJson := NewMassage(500, err.Error())
		PutJson(w, imassageJson)
		return
	}
	return
}
