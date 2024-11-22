package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// temp code to ensure gofmt doesn't remove the imports while unreferenced
var _ = net.Listen
var _ = os.Exit

const port = 4221

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %v\n", port)
		os.Exit(1)
	}

	fmt.Printf("Listening on port %v\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to read data from connection")
		return
	}

	request := string(buf[:n])

	requestParser := RequestParserImpl{}
	requestLine, err := requestParser.parseRequestLine(request)

	if err != nil {
		fmt.Println("Invalid request")
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	fmt.Printf("verb: %v\n", requestLine.Verb)
	fmt.Printf("path: %v\n", requestLine.Path)
	fmt.Printf("version: %v\n", requestLine.Version)

	switch {
	case requestLine.Path == "/":
		handleRootRequest(requestLine.Path, conn)
	case strings.HasPrefix(requestLine.Path, "/echo"):
		handleEchoRequest(requestLine.Path, conn)
	case strings.HasPrefix(requestLine.Path, "/user-agent"):
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	default:
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func handleRootRequest(path string, conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}

func handleEchoRequest(path string, conn net.Conn) {
	fmt.Printf("handling echo request for path: %v\n", path)
	// extract the message from the path
	message := strings.Split(path, "/")[2]
	fmt.Printf("The message is: %v\n", message)

	// construct the response
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)

	conn.Write([]byte(response))
}