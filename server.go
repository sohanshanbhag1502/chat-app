package main

import (
	"encoding/json"
	"fmt"
	"net"
	str "strings"
	"time"
)

type Client struct {
	self  string
	conn  net.Conn
	other string
	queue chan string
}

type Message struct {
	Msg  string `json:"Msg"`
	Info string `json:"Info"`
}

var clients = make(map[string]Client)
var msgsend = make(chan string)

func main() {
	ln, err := net.Listen("tcp", "localhost:5555")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Server listening on localhost:5555")
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		buf := make([]byte, 2048)
		recv_len, err := conn.Read(buf)
		var client Client
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			hosts := str.Split(string(buf[:recv_len]), "-")
			hosts[0] = str.Trim(hosts[0], " ")
			hosts[1] = str.Trim(hosts[1], " ")
			client = Client{self: hosts[0], conn: conn, other: hosts[1]}
			client.queue = make(chan string, 10)
			_, exists := clients[hosts[0]]
			if exists {
				_, err := client.conn.Write(Serialize(Message{Msg: "",
					Info: "USERNAME_TAKEN"}))
				if err != nil {
					fmt.Println(err)
					return
				}
				client.conn.Close()
				continue
			} else {
				clients[hosts[0]] = client
				fmt.Println("Accepted Connection from", hosts[0])
				go sendQueuedMessages(client)
			}
		}

		go handleClient(client)
	}
}

func Serialize(msg Message) []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func DeSerialize(obj []byte) Message {
	var msg = Message{}
	err := json.Unmarshal(obj, &msg)
	if err != nil {
		fmt.Println(err)
	}
	return msg
}

func sendQueuedMessages(client Client) {
	value, exists:= clients[client.other]
	if (exists) {
		if (len(value.queue)!=0){
			for (len(value.queue) != 0) {
				_, err := client.conn.Write(
					Serialize(Message{Msg: <-value.queue, Info: ""}))
				time.Sleep(1 * time.Second)
				if (err != nil) {
					fmt.Println(err)
					break
				}
			}
		}
	}
}

func handleClient(client Client) {
	defer client.conn.Close()

	for {
		buf := make([]byte, 2048)
		recv_len, err := client.conn.Read(buf)
		msg := DeSerialize(buf[:recv_len])

		value, exists := clients[client.other]
		if err != nil {
			fmt.Println(err)
			return
		} else if msg.Info == "CLOSE" {
			fmt.Println("Closed Connection from", client.self)
			delete(clients, client.self)
			client.conn.Close()
			return
		} else if !exists {
			if len(client.queue) != 10 {
				client.queue <- msg.Msg
			} else {
				_, err := client.conn.Write(Serialize(Message{Msg: "",
					Info: "CLIENT_NOT_CONN"}))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		} else {
			otherconn := value.conn
			for len(client.queue) != 0 {
				_, err := otherconn.Write(Serialize(Message{Msg: <-client.queue,
					Info: ""}))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			_, err := otherconn.Write(Serialize(msg))
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = client.conn.Write(Serialize(Message{Msg: "",
				Info: "SUCCESS"}))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
