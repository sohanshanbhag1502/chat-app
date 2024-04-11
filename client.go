package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	str "strings"
	"time"
	"github.com/TwiN/go-color"
)

type Message struct {
	Msg  string `json:"Msg"`
	Info string `json:"Info"`
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

func main() {

	var conn, err = net.Dial("tcp", "localhost:5555")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to localhost:5555\n")

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

	fmt.Println("\nPress Ctrl+C to quit\n")

	go func() {
		for {
			var msg Message = Message{}
			buf := make([]byte, 2048)
			recv_len, err := conn.Read(buf)
			if err != nil {
				return
			}
			msg = DeSerialize(buf[:recv_len])
			if msg.Info == "USERNAME_TAKEN" {
				fmt.Printf("\033[1A\033[K")
				fmt.Printf("\033[1A\033[K")

				fmt.Printf("\b\b\b\b\b")
				fmt.Println(color.Colorize(color.Red, "Username already taken."))
				os.Exit(2)
			}
			if msg.Info == "CLIENT_NOT_CONN" {
				fmt.Printf("\b\b\b\b\b")
				fmt.Println(color.Colorize(color.Red, "You can only send up to 10 messages when client in disconnected."))
			} else if msg.Info == "SUCCESS" {
				continue
			} else {
				fmt.Printf("\b\b\b\b\b")
				fmt.Printf(color.Colorize(color.Cyan, time.Now().Format("15:04:05"))+" - "+color.Colorize(color.Yellow, recipient)+": %s\n", msg.Msg)
				fmt.Printf(color.Colorize(color.Blue, "You: "))
			}
		}
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf(color.Colorize(color.Blue, "You: "))
		var message string = ""
		message, err = reader.ReadString('\n')
		if (err!=nil){
			fmt.Println("")
			fmt.Printf("\033[1A\033[K")
			fmt.Println("\nClosed Connection")
			conn.Write(Serialize(Message{Msg: "", Info: "CLOSE"}))
			return
		}
		fmt.Printf("\033[1A\033[K")
		fmt.Printf(color.Colorize(color.Green, time.Now().Format("15:04:05")+color.Colorize(color.Blue, " - You: ")) + message)
		message = str.Trim(message, "\n")
		conn.Write(Serialize(Message{Msg: message, Info: ""}))
	}
}
