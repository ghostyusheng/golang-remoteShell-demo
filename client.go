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
	n, err := conn.Write([]byte("HEAD / HTTP/1.1\r\n\r\n"))
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
	for {
		n2, err2 := conn.Read(tmp)
		if err2 != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		fmt.Println("got", n, "bytes.")
		fmt.Println("got", string(tmp))
		buf = append(buf, tmp[:n2]...)

	}

	log.Fatal(n)
}
