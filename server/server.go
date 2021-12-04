package main

import (
	"fmt"
	"io"
	"net"
)

func process(conn net.Conn) {
	// 函数执行完毕则关闭conn
	defer conn.Close()

	// 循环读取打印conn中的内容
	for {
		// 先创建一个buffer
		buf := make([]byte, 512)
		// 以buf为尺度去读取数据，类似于用勺子去桶里打水
		n, err := conn.Read(buf)

		if err == io.EOF {
			fmt.Println("数据读取完了")
			return
		}

		if err != nil {
			fmt.Printf("从Conn读取数据出错: %v\n", err)
			return
		}

		// 打印数据
		fmt.Print(string(buf[:n]))
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8888")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("客户端连接出错 %v \n", err)
		}

		fmt.Printf("与客户端[%v]建立连接\n", conn.RemoteAddr())
		go process(conn)
	}
}
