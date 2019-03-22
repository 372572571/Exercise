package chanrpc_test

import (
	"fmt"
	"os"
	"time"
	"github.com/372572571/Exercise/chanrpc"
)

// Example单元测试
func Example() {
	server := chanrpc.NewServer(10)
	server.Registered("T", func(args []interface{}) ([]interface{}, error) {
		arr := []interface{}{1, 2, 3}
		fmt.Println(args)
		return arr, nil
	})

	server.Run()

	client := chanrpc.NewClient(10) // 创建客户
	client.Link(server)             // 连接服务

	client.AsyncRequest("T", []interface{}{9, 9, 9, 9, 9}, func(res *chanrpc.Result) {
		fmt.Println(res, "异步回调")
	})
	client.Request("T", []interface{}{1, 2, 3, 5, 6}, func(res *chanrpc.Result) {
		fmt.Println(res, "回调")
	})
	time.Sleep(5 * time.Second)
	client.AsyncRun()
	var a string

	for {
		i, err := fmt.Scan(&a)
		fmt.Println(i)
		fmt.Println(err)
		fmt.Println(a)
		if a == "exit" {
			os.Exit(0)
		}
	}

}

// GetCurrentPath 获取当前程序运行路径
// func GetCurrentPath() (string, error) {
// 	file, err := exec.LookPath(os.Args[0])
// 	if err != nil {
// 		return "", err
// 	}
// 	path, err := filepath.Abs(file)
// 	if err != nil {
// 		return "", err
// 	}
// 	i := strings.LastIndex(path, "/")
// 	if i < 0 {
// 		i = strings.LastIndex(path, "\\")
// 	}
// 	if i < 0 {
// 		return "", errors.New(`error: Can't find "/" or "\".`)
// 	}
// 	return string(path[0 : i+1]), nil
// }
