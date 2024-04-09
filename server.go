package main

import (
	"fmt"
	"net"
	str "strings"
	"bytes"
	"encoding/gob"
)

type Client struct{
	conn net.Conn
	other string
}

var clients map[string]Client

func main(){
	ln, err := net.Listen("tcp", "localhost:8080")
	if (err!=nil){
		fmt.Println(err)
		return
	}else{
		fmt.Println("Server listening on localhost:8080")
	}

	buf:=make([]byte, 2048)

	for{
		conn, err := ln.Accept()
		if (err!=nil){
			fmt.Println(err)
			continue
		}

		_, err=conn.Read(buf)
		var client Client
		if (err!=nil){
			fmt.Println(err)
			continue
		}else{
			hosts:=str.Split(string(buf), "-")
			client=Client{conn:conn, other:hosts[1]}
			clients[hosts[0]]=client
			fmt.Println("Accepted Connection from ", hosts[0])
		}

		go handleClient(client)
	}
}

func handleClient(client Client){
	defer client.conn.Close()

	msgchan:=make(chan string, 10)

	msgoutbuf:=new(bytes.Buffer)
	encoder:=gob.NewEncoder(msgoutbuf)

	msginbuf:=new(bytes.Buffer)
	decoder:=gob.NewDecoder(msginbuf)

	for{
		buf:=make([]byte, 2048)
		_, err:=client.conn.Read(buf)
		msginbuf.Write(buf)

		var msg []string
		decoder.Decode(&msg)

		value, exists:=clients[client.other]

		if (err!=nil){
			fmt.Println(err)
			return
		}else if (!exists){
			if (len(msgchan)==10){
				msgchan<-msg[0]
			}else{
				encoder.Encode([]string{"","CLIENT_NOT_CONN"})
				_, err:=client.conn.Write(msgoutbuf.Bytes())
				if (err!=nil){
					fmt.Println(err)
					return
				}
			}
		}else{
			otherconn:=value.conn
			for (len(msgchan)!=0){
				encoder.Encode([]string{<-msgchan, ""})
				_, err:=otherconn.Write(msgoutbuf.Bytes())
				if (err!=nil){
					fmt.Println(err)
					return
				}
			}
			encoder.Encode([]string{msg[0], ""})
			_, err:=otherconn.Write(msgoutbuf.Bytes())
			if (err!=nil){
				fmt.Println(err)
				return
			}
		}
	}
}