package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("./client tag remoteIp remotePort Port")
		return
	}

	port, _ := strconv.Atoi(os.Args[4])
	tag := os.Args[1]
	remoteIp := os.Args[2]
	remotePort, _ := strconv.Atoi(os.Args[3])

	LocalAddr := &net.UDPAddr{Port: port}

	conn, err := net.DialUDP("udp", LocalAddr, &net.UDPAddr{IP: net.ParseIP(remoteIp), Port: remotePort})
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write([]byte("I'm: " + tag))

	buf := make([]byte, 256)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Close()
	toAddr := parseAddr(string(buf[:n]))
	fmt.Println("connect to ", toAddr)
	p2p(LocalAddr, &toAddr)
}

func parseAddr(addr string) net.UDPAddr {
	fmt.Println(addr)
	t := strings.Split(addr, ":")
	fmt.Println(t)
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{IP: net.ParseIP(t[0]), Port: port}
}

func p2p(srcAddr *net.UDPAddr, destAddr *net.UDPAddr) {
	conn, _ := net.DialUDP("udp", srcAddr, destAddr)
	conn.Write([]byte("hello\n"))

	go func() {
		buf := make([]byte, 256)
		for {
			n, _, _ := conn.ReadFromUDP(buf)
			if n > 0 {
				fmt.Print("receive message " + string(buf[:n]))
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("p2p >")
		data, _ := reader.ReadString('\n')
		conn.Write([]byte(data))
	}
}
