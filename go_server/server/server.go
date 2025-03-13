package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"unsafe"
)

type ChatMessage struct {
	From, To, Msg string
}

// 与客户端通信
type ClientMsg struct {
	To      string  `json:"to"`
	Msg     string  `json:"msg"`
	Datalen uintptr `json:"datalen"` //消息长度，用于验证数据包长度
}

// channel 消息中心
var chan_msgcenter chan ChatMessage

// 私信的时候必须使用昵称
// 定义名字与地址的映射
var mapName2CliAddr map[string]string
var mapCliaddr2Clients map[string]net.Conn

func handle_conn(conn net.Conn) {
	// 处理新的连接
	from := conn.RemoteAddr().String()
	mapCliaddr2Clients[from] = conn
	msg := ChatMessage{from, "all", from + "-> login"}
	chan_msgcenter <- msg
	// 处理断开连接的事件
	defer logout(conn, from)

	// 分析消息
	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("read error:", err)
			break
		}
		if n > 0 {
			var climsg ClientMsg
			err := json.Unmarshal(buf[:n], &climsg) // json解析
			if err != nil {
				fmt.Println("json error:", err, string(buf[:n]))
				continue
			}
			if climsg.Datalen != unsafe.Sizeof(climsg) {
				fmt.Println("data error:", climsg.Datalen, unsafe.Sizeof(climsg))
			}

			// 处理消息发送到消息中心
			chatmsg := ChatMessage{from, "all", climsg.Msg}
			switch climsg.To {
			case "all":
			case "set": // 设置用户名
				mapName2CliAddr[climsg.Msg] = from
				chatmsg.Msg = from + "set name = " + climsg.Msg + " success"
				chatmsg.From = "server"
			default:
				chatmsg.To = climsg.To
			}
			chan_msgcenter <- chatmsg
		}
	}
}

func logout(conn net.Conn, from string) {
	defer conn.Close()
	delete(mapCliaddr2Clients, from) // 删除map中的元素
	// 通知消息中心
	msg := ChatMessage{from, "all", from + "-> logout"}
	chan_msgcenter <- msg
}

func msg_center() {
	for {
		msg := <-chan_msgcenter
		go send_msg(msg)
	}
}

func send_msg(msg ChatMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("json.Marshal err:", err)
		return
	}
	if msg.To == "all" {
		for _, v := range mapCliaddr2Clients {
			if msg.From != v.RemoteAddr().String() {
				v.Write(data)
			}
		}
	} else {
		from, ok := mapName2CliAddr[msg.To]
		if !ok {
			fmt.Println("no such user", msg.To)
			return
		}
		conn, ok := mapCliaddr2Clients[from]
		if !ok {
			fmt.Println("no such client", from, msg.To)
			return
		}
		conn.Write(data)
	}
}

func main() {
	mapCliaddr2Clients = make(map[string]net.Conn)
	mapName2CliAddr = make(map[string]string)
	chan_msgcenter = make(chan ChatMessage)

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer listener.Close()

	go msg_center()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handle_conn(conn)
	}
}
