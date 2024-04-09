package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to localhost:8080")

	reader := bufio.NewReader(os.Stdin)
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

	buffer := new(bytes.Buffer)
	buffer_out := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	decoder := gob.NewDecoder(buffer_out)

	go func() {
		for {
			var msg []string
			buf := make([]byte, 2048)
			_, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error decoding message:", err)
			}
			buffer_out.Write(buf)
			err = decoder.Decode(&msg)
			if err != nil {
				fmt.Println("Error decoding message:", err)

			}
			if msg[1] == "CLIENT_NOT_CONN" {
				fmt.Println("Recipient client not connected.")

			}
			fmt.Printf("%s : %s \n", recipient, msg[0])
		}
	}()

	for {

		message, _ := reader.ReadString('\n')

		if message == "exit" {
			err := encoder.Encode([]string{"", "CLOSE"})
			if err != nil {
				fmt.Println("Error encoding message:", err)
			}
			break
		}

		err := encoder.Encode([]string{message, ""})

		if err != nil {
			fmt.Println("Error encoding message:", err)
			break
		}
		conn.Write(buffer.Bytes())
		fmt.Println("You : ", message)

	}
}
