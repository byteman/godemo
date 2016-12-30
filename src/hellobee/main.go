package main

import (
	"fmt"
	_ "hellobee/routers"
	"net"

	"hellobee/czserver"

	"github.com/astaxie/beego"
)

func handleConn(c *czserver.NetClient) {
	defer c.Con.Close()
	defer czserver.RemoveClient(c.Con)
	for {
		b := make([]byte, 512)
		n, err := c.Con.Read(b)
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		if n == 0 {
			fmt.Println("recv 0")

			break
		}
		c.Handle(b, n)
	}
}

func startTcpServer() {
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

		client := czserver.CreateClient(c)
		go handleConn(client)

	}
}
func main() {

	go startTcpServer()
	beego.SetStaticPath("/static/html", "device")
	beego.Run()
}
