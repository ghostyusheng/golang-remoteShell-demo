package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s host:port", os.Args[0])
	}
	service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		buf = append(buf, tmp[:n]...)

		if n > 0 {
			runCmd(string(tmp))
		}
	}

	log.Fatal("finish")
}

func runCmd(_cmd string) {
	fmt.Println(_cmd)
	//var cmd *exec.Cmd
	//cmd = exec.Command("date")
	//args := strings.Fields(_cmd + " ")
	//if len(args) < 1 {
	//	return
	//} else if len(args) == 2 {
	//	fmt.Println("receive: ", args)
	//	cmd = exec.Command(args[0])
	//} else {
	//	fmt.Println("receive: ", args)
	//}
	//var stdout, stderr bytes.Buffer
	//cmd.Stdout = &stdout
	//cmd.Stderr = &stderr
	//err := cmd.Run()
	//if err != nil {
	//	log.Println("cmd.Run() failed with %s\n", err)
	//	return
	//}
	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}
