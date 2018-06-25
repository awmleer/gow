package main

import (
	"net"
	"bufio"
	"fmt"
	"strings"
	"io/ioutil"
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
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	flag := true
	var method string
	var uri string
	for {
		message, _ := rw.ReadString('\n')
		if flag {
			firstLn := strings.Split(message, " ")
			method = firstLn[0]
			uri = firstLn[1]
			fmt.Println(uri)
			flag = false
		}
		// output message received
		fmt.Print(string(message))
		if string(message) == "\r\n" {
			fmt.Print("aaa")
			break
		}
	}
	if method == "GET" {
		rw.WriteString("HTTP/1.1 200 OK\r\n")
		rw.WriteString("Content-Type: text/html\r\n")
		rw.WriteString("Content-Length: 151\r\n")
		rw.WriteString("\r\n")
		bytes, _ := ioutil.ReadFile("/Users/awmleer/Project/go/src/github.com/awmleer/gow/test/index.html")
		rw.Write(bytes)
		rw.Flush()
	}
	conn.Close()

	//// sample process for string received
	//newMessage := strings.ToUpper(message)
	//// send new string back to client
	//conn.Write([]byte(newMessage + "\n"))
}