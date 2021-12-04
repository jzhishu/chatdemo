package main

import (
	"fmt"
	"io"
	"net"
	"reflect"
)

func writeActiveChan(conn net.Conn, activeChan chan net.Conn) {
	fmt.Println(reflect.TypeOf(conn).Kind())

	str := ""
	fmt.Println("接收到服务端的消息，开始读取")
	// 循环读取打印conn中的内容
	for {
		fmt.Println("循环读取中")
		// 先创建一个buffer
		buf := make([]byte, 512)
		// 以buf为尺度去读取数据，类似于用勺子去桶里打水
		n, err := conn.Read(buf)

		if err == io.EOF {
			// fmt.Println("数据读取完了")
			break
		}

		if err != nil {
			fmt.Printf("从Conn读取数据出错: %v\n", err)
			break
		}

		// 把结果数据拼接
		str += string(buf[:n])
	}
	fmt.Println("数据读取拼接完毕，将conn放到active管道中，准备想其他所有conn同步")
	if str != "" {
		activeChan <- conn
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8888")

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("服务成功启动")
	defer listener.Close()

	// 创建一个保存所有连接的切片
	connList := []net.Conn{}

	// 创建一个保存活跃消息的管道
	activeChan := make(chan net.Conn, 500)

	for {
		conn, err := listener.Accept()
		// connList = append(connList, &conn)

		// activeChan <- &conn
		if err != nil {
			fmt.Printf("客户端连接出错 %v \n", err)
		}

		fmt.Printf("客户端[%v]成功登录\n", conn.RemoteAddr())
		go writeActiveChan(conn, activeChan)
	}
}
