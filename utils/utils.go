/*
 * @Description: In User Settings Edh
 * @Author: your name
 * @Date: 2019-09-19 16:43:45
 * @LastEditTime: 2019-09-19 21:31:50
 * @LastEditors: Please set LastEditors
 */

package utils

import (
	"io"
	"net"
	"net/http"
	"os"
)

// GetLoopbackAddress 获取本机本地地址(回环地址)
// ip 如果没有获取到本机的ipv4 默认“0.0.0.0”
// err 获取接口地址数组现象错误
func GetLoopbackAddress() (ip string, err error) {
	addres, err := net.InterfaceAddrs()
	ip = "127.0.0.1"
	if err != nil {
		return
	}
	for _, val := range addres {
		// val.(*net.IPNet) 强制类型转换
		if ipNet, ok := val.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			// ipNet, _ := val.(*net.IPNet)
			if ipNet.IP.String() != ip {
				// fmt.Println(key, "<=>", ipNet.IP.String())
				ip = ipNet.IP.String()
				return
			}
		}
	}
	return
}

// GetIPAddress ... 根据url 获取IP地址.
// url 根据这个url获取地址
// utils.GetIPAddress("http://myexternalip.com/raw")
func GetIPAddress(url string) (ip string, err error) {
	ip = "0.0.0.0"
	err = nil

	resp, err := http.Get(url)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	os.Exit(0)
	return
}
