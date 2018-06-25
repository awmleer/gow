package main

import (
	"net"
	"bufio"
	"fmt"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":4785")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	flag := true
	for {
		message, _ := reader.ReadString('\n')
		if flag {
			firstLn := strings.Split(message, " ")
			method := firstLn[0]
			uri := firstLn[1]
			flag = false
		}
		// output message received
		fmt.Print(string(message))
		if string(message) == "\r\n" {
			fmt.Print("aaa")
		}
	}
	//// sample process for string received
	//newMessage := strings.ToUpper(message)
	//// send new string back to client
	//conn.Write([]byte(newMessage + "\n"))
}