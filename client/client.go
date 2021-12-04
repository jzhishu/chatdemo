package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	// 创建一个命令行输入的reader
	reader := bufio.NewReader(os.Stdin)

	for {
		// 从reader中读取数据, ReadString可以读到指定符号就终止

		str, err := reader.ReadString('\n')

		line := strings.Trim(str, " \r\n")
		if line == "exit" {
			os.Exit(1)
			break
		}

		if err != nil {
			fmt.Println("读取命令行输入失败", err)
			return
		}

		// 将从控制台读取到的数据写入到conn中
		conn.Write([]byte(str))
	}
}
