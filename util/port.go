package util

import (
	"fmt"
	"net"
)

// todo 生成随机端口号
func GenRandomPort() int {
	//	解析TCP地址（当指定的端口为0时，系统会自动配置一个随机可用端口）
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	//	根据解析的地址获取监听器
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	//	延迟关闭监听器
	defer listener.Close()
	//	获取监听器绑定的端口
	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Println(port)
	return port
}
