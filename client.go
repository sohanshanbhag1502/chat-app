package main

import (
	"fmt"
	"net"
	"encoding/json"
)

type Message struct{
	Msg string `json:"Msg"`
	Info string `json:"Info"`
}

func Serialize(msg Message) []byte {
	b, err:=json.Marshal(msg)
	if (err!=nil){
		fmt.Println(err)
	}
	fmt.Println(b)
	return b
}

func DeSerialize(obj []byte) Message{
	var msg=Message{}
	err:=json.Unmarshal(obj, &msg)
	if (err!=nil){
		fmt.Println(err)
	}
	return msg
} 

func main() {
	conn, err := net.Dial("tcp", "localhost:5000")

	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to localhost:5000")

	name := ""
	fmt.Print("Enter your name: ")
	fmt.Scanln(&name) 
	recipient := ""
	fmt.Print("Enter the Recipient's Name: ")
	fmt.Scanln(&recipient)

	_, err = conn.Write([]byte(name + "-" + recipient))
	if err != nil {
		fmt.Println("Error sending name:", err)
		return
	}

	fmt.Println("Connected to", recipient)
	fmt.Println("Type 'exit' to quit")

	go func() {
		for {
			var msg Message=Message{}
			buf := make([]byte, 2048)
			recv_len, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error decoding message:", err)
			}
			msg=DeSerialize(buf[:recv_len])
			fmt.Println(msg)
			break
			if msg.Info == "CLIENT_NOT_CONN" {
				fmt.Println("Recipient client not connected.")
			}else if msg.Info=="SUCCESS"{
				continue
			}else{
				fmt.Printf("%s : %s \n", recipient, msg.Msg)
			}
		}
	}()

	for {
		var message string=""
		fmt.Scanln(&message)

		if message == "exit" {
			conn.Write(Serialize(Message{Msg:"", Info: "CLOSE"}))
			break
		}
		conn.Write(Serialize(Message{Msg:message, Info:""}))
		fmt.Println("You : ", message)
	}
}
