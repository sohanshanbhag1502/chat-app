package main

import (
	"encoding/json"
	"fmt"
	"net"
	str "strings"
	"time"
	dotenv "github.com/joho/godotenv"
)

type Message struct {
	Msg       string `json:"Msg"`
	Info      string `json:"Info"`
	Time_stmp string `json:"Time_stmp"`
}

type Client struct {
	self  string
	conn  net.Conn
	other string
	queue chan Message
}

var ip_port string;
var clients = make(map[string]Client)
var msgsend = make(chan string)

func main() {
	envFile, err := dotenv.Read(".env")
	if (err != nil) {
		ip_port = "localhost:8080"
	} else{
		ip_port = envFile["HOST"] + ":" + envFile["PORT"]
	}


	ln, err := net.Listen("tcp", ip_port)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Print("Server listening on " + ip_port + "\n")
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
			client.queue = make(chan Message, 10)
			_, exists := clients[hosts[0]]
			if exists {
				_, err := client.conn.Write(Serialize(Message{Msg: "",
					Info: "USERNAME_TAKEN", Time_stmp: ""}))
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
				go handleClient(client)
			}
		}
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
	value, exists := clients[client.other]
	if exists {
		if len(value.queue) != 0 {
			for len(value.queue) != 0 {
				_, err := client.conn.Write(Serialize(<-value.queue))
				time.Sleep(1 * time.Second)
				if err != nil {
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
				client.queue <- msg
			} else {
				_, err := client.conn.Write(Serialize(Message{Msg: "",
					Info: "CLIENT_NOT_CONN_BUFFER_FULL", Time_stmp: ""}))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		} else {
			otherconn := value.conn
			for len(client.queue) != 0 {
				_, err := otherconn.Write(Serialize(<-client.queue))
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
				Info: "SUCCESS", Time_stmp: ""}))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
