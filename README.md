# Chat App

A simple chat app built on Golang using socket programming in CLI.

## Dependencies:
- [github.com/TwiN/go-color](https://github.com/TwiN/go-color)
    ```
    go get github.com/TwiN/go-color
    ```
- [github.com/joho/godotenv](https://github.com/joho/godotenv)
    ```
    go get github.com/joho/godotenv
    ```

## Instructions to run the program:
- Create a .env file with following keys:
    * HOST: The ip address of server.
    * PORT: The port number on which the server should accept new connections.

- To run the server script run the following command:
    ```
    go run server.go
    ```

- To run the client script run the following command:
    ```
    go run client.go
    ```

- To build the executable files run the following commands:
    ```
    go build server.go
    go build client.go
    ```
