package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	str "strings"
	"syscall"
	"time"

	"github.com/TwiN/go-color"
)

var ip_port = "10.20.200.141:64000"

type Message struct {
	Msg       string `json:"Msg"`
	Info      string `json:"Info"`
	Time_stmp string `json:"Time_stmp"`
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
	os_windows := false
	os_windows = runtime.GOOS=="windows"

	var conn, err = net.Dial("tcp", ip_port)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Print("Connected to " + ip_port + "\n")

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
				// fmt.Println("Error recieving message:", err)
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
			if msg.Info == "CLIENT_NOT_CONN_BUFFER_FULL" {
				fmt.Printf("\b\b\b\b\b")
				fmt.Println(color.Colorize(color.Red, "You can only send up to 10 messages when client in disconnected."))
				fmt.Printf(color.Colorize(color.Blue, "You: "))
			} else if msg.Info == "SUCCESS" {
				continue
			} else {
				fmt.Printf("\b\b\b\b\b")
				fmt.Printf(color.Colorize(color.Cyan, msg.Time_stmp)+" - "+color.Colorize(color.Yellow, recipient)+": %s\n", msg.Msg)
				fmt.Printf(color.Colorize(color.Blue, "You: "))
			}
		}
	}()

	if (!os_windows) {
		fmt.Println("Running...")
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			fmt.Println("\nClosed Connection...")
			conn.Write(Serialize(Message{Msg: "", Info: "CLOSE", Time_stmp: ""}))
			os.Exit(1)
		}()
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf(color.Colorize(color.Blue, "You: "))
		var message string = ""
		message, err = reader.ReadString('\n')

		if (os_windows && err != nil) {
			fmt.Println("")
			fmt.Printf("\033[1A\033[K")
			fmt.Println("\nClosed Connection")
			conn.Write(Serialize(Message{Msg: "", Info: "CLOSE", Time_stmp: ""}))
			return
		}
		fmt.Printf("\033[1A\033[K")
		fmt.Printf(color.Colorize(color.Green, time.Now().Format("15:04")+color.Colorize(color.Blue, " - You: ")) + message)
		message = str.Trim(message, "\n")
		message = str.Trim(message, "\r")
		_, err := conn.Write(Serialize(Message{Msg: message, Info: "", Time_stmp: time.Now().Format("15:04")}))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}
