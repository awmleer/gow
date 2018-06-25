package main

import (
	"net"
	"bufio"
	"fmt"
	"strings"
	"io/ioutil"
	"os"
)

var servePath string

func main() {
	ln, err := net.Listen("tcp", ":4785")
	if err != nil {
		// handle error
	}
	if len(os.Args) < 2 {
		fmt.Println("Please input a directory to serve.")
		return
	}
	servePath = os.Args[1]
	if servePath[len(servePath)-1:] == "/" {
		servePath = servePath[:len(servePath)-1]
	}
	fmt.Println("Serving directory: " + servePath)
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
		data, err := ioutil.ReadFile(servePath + uri)
		if err != nil {
			rw.WriteString("HTTP/1.1 404 NOT FOUND\r\n")
			rw.WriteString("Content-Length: 35\r\n")
			rw.WriteString("Content-Type: text/html\r\n")
			rw.WriteString("\r\n")
			rw.WriteString("<html><body>Not Found</body></html>")
		} else {
			rw.WriteString("HTTP/1.1 200 OK\r\n")
			contentLength := len(data)
			rw.WriteString("Content-Length: " + string(contentLength) + "\r\n")
			uriParts := strings.Split(uri, ".")
			fileExtention := uriParts[len(uriParts)-1]
			var contentType string
			switch fileExtention {
			case "html":
				contentType = "text/html"
			case "jpg":
				contentType = "image/jpeg"
			case "txt":
				contentType = "text/plain"
			default:
				contentType = "text/plain"
			}
			rw.WriteString("Content-Type: " + contentType + "\r\n")
			rw.WriteString("\r\n")
			rw.Write(data)
		}
		rw.Flush()
	}
	conn.Close()
}
