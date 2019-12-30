package compip

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)

// 客户端的实现
type Client struct {
	sync.Mutex
	Path      string // 可执行程序目录
	clientCmd *exec.Cmd
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Run(name string, arg ...string) {
	c.clientCmd = exec.Command(name, "logcat")
	// c.clientCmd.Args = arg
	// c.clientCmd.Path = c.Path

	c.read()
	err := c.clientCmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.clientCmd.Wait()
	if err != nil {
		fmt.Println(err)
		return
	}
	// for {
	// }
}

func (c *Client) read() {
	cmdReader, err := c.clientCmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()
}
