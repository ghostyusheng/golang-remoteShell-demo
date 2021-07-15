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

func stdinput() string {
	fmt.Println("<==== input command ====>")
	var cmd, _ = reader.ReadString('\n')
	cmd = strings.Replace(cmd, "\n", "", -1)
	var cmdArr = strings.Split(cmd, " ")
	if len(cmdArr) < 2 {
		fmt.Println("you need to input port and command like: $port $command")
		return ""
	}
	return cmd
}

func inArray(need string, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

func sendCommand(conn *net.TCPConn, cmd string) {
	n, err := conn.Write([]byte(cmd))
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	fmt.Printf("send %d bytes to %s\n", n, conn.RemoteAddr())
}

func connHandler(conn *net.TCPConn) {
	tick := time.Tick(3 * time.Second)
	for now := range tick {
		if globalCommand == "" {
			continue
		}

		var arr = strings.Split(globalCommand, " ")
		var _port = arr[0]
		var cmd = strings.Join(arr[1:], " ")
		var host_port_ = conn.RemoteAddr().String()
		var port_ = strings.Split(host_port_, ":")[1]
		if _port == port_ && cmd != M[_port] {
			fmt.Println(dt(now))
			sendCommand(conn, cmd)
			M[_port] = cmd
		}
	}
}

func dt(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
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

var globalCommand = ""
var M = make(map[string]string)
var S = []string{}

func globalInputScopeControl() {
	tick := time.Tick(2 * time.Second) // 五秒的心跳间隔
	for now := range tick {
		fmt.Println(dt(now), S)
		globalCommand = stdinput()
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
		S = append(S, conn.RemoteAddr().String())
		go globalInputScopeControl()
		go connHandler(conn)
		go reply(conn)
	}
}
