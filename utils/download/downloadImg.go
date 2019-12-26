package download

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

// DownloadImg ...结构
// ...结构
type DownloadImg struct {
	OutDirPath string
	Url        string
	Size       int
	Name       string
	ResPath    string
}

// NewDownloadImg ...
// 创建下载对象
func NewDownloadImg(url string, size int, _path string, fix string) *DownloadImg {
	var res = &DownloadImg{}
	res.Url = url
	res.Size = size
	res.OutDirPath, res.Name = path.Split(_path)
	res.OutDirPath = path.Join(res.OutDirPath, fix)
	res.ResPath = path.Join(res.OutDirPath, res.Name)
	return res
}

// Download 下载图片并保存到目录
func (down DownloadImg) Download() (string, error) {
	req, err := http.NewRequest(http.MethodGet, down.Url, nil)

	if err != nil {
		return "", errors.New("创建请求失败:" + err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("发起请求失败:" + err.Error())
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("读取失败:" + err.Error()) //
	}
	file, err := os.Create(down.ResPath)
	if err != nil {

		return "", errors.New("创建文件失败:" + err.Error())
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	io.Copy(writer, bytes.NewReader(data)) // 数据写入
	return down.ResPath, nil
}
