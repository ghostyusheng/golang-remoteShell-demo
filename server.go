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
	var lastCmdPort string = ""
	var lastCmd string = ""
	for now := range tick {
		var rAddr string = conn.RemoteAddr().String()
		var cmd string
		if lastCmdPort == "" {
			fmt.Println("<==== input command ====>")
			cmd, _ = reader.ReadString('\n')
			cmd = strings.Replace(cmd, "\n", "", -1)
			fmt.Println(now, "cmd: ", cmd)
		}
		var cmdArr = strings.Split(lastCmd, " ")
		var chickenPort string = cmdArr[0]
		fmt.Println("remote address: ", rAddr)
		fmt.Println("remote port: ", chickenPort)
		if lastCmd != "" && len(cmdArr) < 2 {
			fmt.Println("you need to input port and command like: $port $command")
			continue
		}
		var realCmd = strings.Join(cmdArr[1:], " ")
		fmt.Println("remote command: ", realCmd)
		var controlPort string = strings.Split(rAddr, ":")[1]
		if lastCmdPort == controlPort {
			fmt.Println("last command: ", lastCmd)
			n, err := conn.Write([]byte(lastCmd))
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
			fmt.Printf("send %d bytes to %s\n", n, rAddr)
			fmt.Printf("send %s\n", string(lastCmd))
			lastCmdPort = ""
		} else {
			//fmt.Println("alternate ${port} ${cmd}: ", rAddr)
			lastCmdPort = chickenPort
			lastCmd = cmd
			continue
		}
	}
}

func reply(conn *net.TCPConn) {
	for {
		tmp := make([]byte, 256)
		m, _ := conn.Read(tmp)
		if m > 0 {
			resp_str := string(tmp[:m])
			fmt.Println(resp_str)
			tmp = make([]byte, 256)
		}
	}
}

var reader = bufio.NewReader(os.Stdin)

var ip = "0.0.0.0"

func main() {
	address := net.TCPAddr{
		IP:   net.ParseIP(ip),
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
		go echo(conn)
		go reply(conn)
	}
}
