package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// 连接服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("已连接到服务器")

	// 获取用户输入的数据
	var input string
	fmt.Print("请输入要发送给服务器的消息: ")
	fmt.Scanln(&input)

	// 发送数据到服务器
	_, err = conn.Write([]byte(input))
	if err != nil {
		fmt.Println("向服务器发送数据失败:", err)
		os.Exit(1)
	}

	// 接收服务器响应的数据
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("接收服务器响应失败:", err)
		os.Exit(1)
	}
	receivedData := string(buffer[:n])
	fmt.Println("从服务器接收到响应:", receivedData)
}
