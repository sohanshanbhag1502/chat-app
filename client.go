package main

import (
	"fmt"
	"net"
	str "strings"
	"bytes"
	"encoding/gob"
)


func main(){

	conn, err := net.Dial("tcp", "localhost:8080")
	if (err!=nil){
		fmt.Println(err)
		return
	}else{
		fmt.Println("Connected to localhost:8080")
	}

	
	name :=""
	fmt.Println("Enter your name: ")
	fmt.Scanf("%s", &name)

	name_with :=""
	fmt.Println("Enter the Recipent Name: ")
	fmt.Scanf("%s", &name_with)


	_, err=conn.Write([]byte(name+"-"+name_with))

	if (err!=nil){
		fmt.Println(err)
		return
	}
	
	fmt.Println("Connected to ", name_with)

	buffer := new(bytes.Buffer);

	encoder := gob.NewEncoder(buffer)

	message := ""

	for{fmt.Println("Enter your message: ")
	fmt.Scanf("%s",&message)

	if (message == "exit"){
		break
	}

	encoder.Encode(message)
	conn.Write(buffer.Bytes())

	}

	conn.Close()

}