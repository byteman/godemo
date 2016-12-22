// mydemo project main.go
package main

import (
	"fmt"
	"net"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		b := make([]byte, 100)
		n, err := c.Read(b)
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		if n == 0 {
			fmt.Println("recv 0")
			break
		}
		fmt.Println("read len= ", n)
		s := string(b)
		fmt.Println(s)
	}
}

func main() {
	fmt.Println("Ready Accept!")
	l, err := net.Listen("tcp", ":8889")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		fmt.Println("accept", c.RemoteAddr().String())
		go handleConn(c)

	}
}
