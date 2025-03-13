package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"unsafe"
)

// 与客户端通信
type ClientMsg struct {
	To      string  `json:"to"`
	Msg     string  `json:"msg"`
	Datalen uintptr `json:"datalen"` //消息长度，用于验证数据包长度
}

func help() {
	fmt.Println("1. set:your name")
	fmt.Println("2. all:your msg -- broadcast")
	fmt.Println("1. anyone:your msg -- private msg")
}

func handle_conn(conn net.Conn) {
	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("read error:", err)
		}
		fmt.Println(string(buf[:n]))
		fmt.Printf("firefly's chat >")
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println("connect error:", err)
	}
	defer conn.Close()
	go handle_conn(conn)
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Wlcome to firefly's chat room!\n")
	help()
	for {
		fmt.Printf("firefly's chat >")
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("read error:", err)
		}
		msg = strings.TrimSpace(msg)

		if msg == "quit" {
			fmt.Println("bye")
			break
		}
		if msg == "help" {
			help()
			continue
		}
		msgs := strings.Split(msg, ":")
		if len(msgs) == 2 {
			var climsg ClientMsg
			climsg.To = msgs[0]
			climsg.Msg = msgs[1]
			climsg.Datalen = unsafe.Sizeof(climsg)

			data, err := json.Marshal(climsg)
			if err != nil {
				log.Println("json marshal error:", err)
				continue
			}
			_, err = conn.Write(data)
			if err != nil {
				log.Println("write error:", err)
				break
			}
		}
	}

}
