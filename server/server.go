package main

import (
	"fmt"
	"io"
	"net"
)

type Messager struct {
	conn    net.Conn
	content string
}

func writeActiveChan(conn net.Conn, activeChan chan Messager) {
	// fmt.Println(reflect.TypeOf(conn).Kind())

	str := ""
	fmt.Println("接收到服务端的消息，开始读取")
	// 循环读取打印conn中的内容
	for {
		fmt.Println("循环读取中")
		// 先创建一个buffer
		buf := make([]byte, 512)
		// 以buf为尺度去读取数据，类似于用勺子去桶里打水
		n, err := conn.Read(buf)
		fmt.Println("读到了数据")
		if err == io.EOF || n < 512 {
			fmt.Println("读到了最后一行")
			// 如果到末尾了
			// 把结果数据拼接
			str += string(buf[:n])
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
		msg := Messager{conn: conn, content: str}
		fmt.Printf("数据加到管道中了: %v", msg)
		activeChan <- msg
	}
	writeActiveChan(conn, activeChan)
}

func broadcast(connChan chan net.Conn, activeChan chan Messager) {
	for {
		if len(activeChan) == 0 {
			continue
		}
		fmt.Println("开始从管道取值，准备广播")
		// 循环的从chan中取出值
		messager := <-activeChan
		fmt.Println("获取到待同步数据: ", messager.content)
		cacheChan := make(chan net.Conn, 500)

		for len(connChan) > 0 {
			item := <-connChan
			cacheChan <- item
			if messager.conn.RemoteAddr() == item.RemoteAddr() {
				continue
			}
			fmt.Println("正在同步消息: ", messager.content)
			item.Write([]byte(messager.content))
		}
		connChan = cacheChan
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
	connChan := make(chan net.Conn, 500)

	// 创建一个保存活跃消息的管道
	activeChan := make(chan Messager, 500)
	go broadcast(connChan, activeChan)
	for {
		conn, err := listener.Accept()
		connChan <- conn

		// activeChan <- &conn
		if err != nil {
			fmt.Printf("客户端连接出错 %v \n", err)
		}

		fmt.Printf("客户端[%v]成功登录\n", conn.RemoteAddr())
		go writeActiveChan(conn, activeChan)

	}
}
