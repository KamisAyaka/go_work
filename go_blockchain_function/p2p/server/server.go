package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	listener, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 9527})
	defer listener.Close()

	fmt.Println("server start at", listener.LocalAddr().String())

	peers := make([]*net.UDPAddr, 2)
	buf := make([]byte, 256)

	n, addr, _ := listener.ReadFromUDP(buf)

	fmt.Printf("receive from %s,data:%s\n", addr.String(), buf[:n])
	peers[0] = addr

	n, addr, _ = listener.ReadFromUDP(buf)
	fmt.Printf("receive from %s,data:%s\n", addr.String(), buf[:n])
	peers[1] = addr

	fmt.Println("begin nat")

	listener.WriteToUDP([]byte(peers[0].String()), peers[1])
	listener.WriteToUDP([]byte(peers[1].String()), peers[0])
	time.Sleep(time.Second * 10)
}
