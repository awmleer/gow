package main

import (
	"net"
	"bufio"
	"fmt"
	"strings"
	"io/ioutil"
	"os"
	"strconv"
)

var servePath string

func main() {
	ln, err := net.Listen("tcp", ":4785")
	if err != nil {
		fmt.Println("Failed to listen at port 4785.")
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("Please input a directory to serve.")
		return
	}
	servePath = os.Args[1]
	if strings.HasSuffix(servePath, "/") {
		servePath = servePath[:len(servePath)-1]
	}
	fmt.Println("Serving directory: " + servePath)
	for {
		conn, _ := ln.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	flag := true
	var requestContentLength int
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
		}else {
			if strings.HasPrefix(message, "Content-Length: "){
				requestContentLength, _ = strconv.Atoi(message[16:len(message)-2])
			}
		}
		// output message received
		fmt.Print(string(message))
		if string(message) == "\r\n" {
			break
		}
	}
	switch method {
	case "GET":
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
			fileExtension := uriParts[len(uriParts)-1]
			var contentType string
			switch fileExtension {
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
	case "POST":
		if uri != "/dopost" {
			rw.WriteString("HTTP/1.1 404 NOT FOUND\r\n")
			rw.WriteString("Content-Length: 35\r\n")
			rw.WriteString("Content-Type: text/html\r\n")
			rw.WriteString("\r\n")
			rw.WriteString("<html><body>Not Found</body></html>")
		} else {
			postBodyBytes := make([]byte, requestContentLength)
			rw.Read(postBodyBytes)
			postBody := string(postBodyBytes)
			fmt.Println(postBody)
			items := strings.Split(postBody, "&")
			var login string
			var pass string
			for _, v := range items {
				a := strings.Split(v, "=")
				switch a[0] {
				case "login":
					login = a[1]
				case "pass":
					pass = a[1]
				}
			}
			rw.WriteString("HTTP/1.1 200 OK\r\n")
			rw.WriteString("Content-Length: 30\r\n")
			rw.WriteString("Content-Type: text/html\r\n")
			rw.WriteString("\r\n")
			if login == "3150104785" && pass == "4785" {
				rw.WriteString("<html><body>登录成功</body></html>")
			} else {
				rw.WriteString("<html><body>登录失败</body></html>")
			}
		}
	}
	rw.Flush()
	conn.Close()
}
