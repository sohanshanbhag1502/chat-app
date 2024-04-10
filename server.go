package main

import (
	"encoding/json"
	"fmt"
	"net"
	str "strings"
)

type Client struct {
	self  string
	conn  net.Conn
	other string
}

type Message struct {
	Msg  string `json:"Msg"`
	Info string `json:"Info"`
}

var clients = make(map[string]Client)

func main() {
	ln, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Server listening on localhost:5000")
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
			clients[hosts[0]] = client
			fmt.Println("Accepted Connection from", hosts[0])
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

func handleClient(client Client) {
	defer client.conn.Close()

	msgchan := make(chan string, 10)

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
			if len(msgchan) != 10 {
				msgchan <- msg.Msg
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
			for len(msgchan) != 0 {
				_, err := otherconn.Write(Serialize(Message{Msg: <-msgchan, 
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
