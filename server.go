package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func echo(conn *net.TCPConn) {
	tick := time.Tick(2 * time.Second) // 五秒的心跳间隔
	for now := range tick {
		cmd, _ := reader.ReadString('\n')
		cmd = strings.Replace(cmd, "\n", "", -1)
		fmt.Println(now, "cmd: ", cmd)

		n, err := conn.Write([]byte(cmd))
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		fmt.Printf("send %d bytes to %s\n", n, conn.RemoteAddr())
		fmt.Printf("send %s\n", string(cmd))
	}
}

var reader = bufio.NewReader(os.Stdin)

func main() {
	address := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"), // 把字符串IP地址转换为net.IP类型
		Port: 8000,
	}
	listener, err := net.ListenTCP("tcp4", &address) // 创建TCP4服务器端监听器
	if err != nil {
		log.Fatal(err) // Println + os.Exit(1)
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err) // 错误直接退出
		}
		fmt.Println("remote address:", conn.RemoteAddr())
		echo(conn)
	}
}
