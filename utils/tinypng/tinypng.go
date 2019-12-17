package tinypng

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/372572571/Exercise/utils/download"

	jsoniter "github.com/json-iterator/go"
)

var (
	CompressingUrl = "https://api.tinify.com/shrink"
)

const (
	FIX = "temp"
)

type input struct {
	Size int
	Type string
}
type output struct {
	Height int
	ratio  int
	Size   int
	Type   string
	Url    string
	Width  int
}

// TinyInfo ........
type TinyInfo struct {
	Input  input
	Output output
}

// TinyPng .......
type TinyPng struct {
	Email   string
	ApiKey  string
	Imgpath string
}

// Run 执行压缩图片
func (t *TinyPng) Run() (string, error) {

	// 创建Request
	PathExists(t.Imgpath, FIX)
	// return
	req, err := http.NewRequest(http.MethodPost, CompressingUrl, nil)
	if err != nil {
		return "", err
	}

	// 将鉴权信息写入Request
	req.SetBasicAuth(t.Email, t.ApiKey)

	// 将图片以二进制的形式写入Request
	data, err := ioutil.ReadFile(t.Imgpath)
	if err != nil {
		return "", err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(data))

	// 发起请求
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	// 解析请求
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var a = &TinyInfo{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	json.Unmarshal(data, a)
	var d = download.NewDownloadImg(a.Output.Url, a.Output.Size, t.Imgpath, FIX)
	p, err := d.Download()
	if err != nil {
		return "", err
	}
	return p, nil
}

// PathExists ...
// ......
func PathExists(dir string, fix string) (bool, error) {
	dir, _ = path.Split(dir)
	_, err := os.Stat(path.Join(dir, fix))
	if err == nil {
		return true, nil // 文件存在
	}
	if os.IsNotExist(err) {
		// fmt.Println(path.Join(dir, fix))
		os.Mkdir(path.Join(dir, fix), os.ModePerm)
		return false, nil
	}
	return false, err
}
