package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
)


func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// Uncomment this block to pass the first stage
	//
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()
	// defer func() {
	// 	listener.Close()
	// }()

	fmt.Println("Listening on " + listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
		}
		fmt.Println("Accepted connection from " + conn.RemoteAddr().String())
		// go handleConnection(conn)
		go echoExecute(conn)
	}
}



func handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	message := []byte("+PONG\r\n")
	fmt.Println("Handling connection")
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Failed to read data")
			return
		}
		conn.Write(message)
	}
}

func echoExecute(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Read input from client
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed")
			return
		}

		// Parse the command
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])
		args := make([]Value, len(parts[1:]))
		for i, arg := range parts[1:] {
			args[i] = Value{typ: "string", bulk: arg}
		}

		// Execute the command
		if handler, exists := Commands[command]; exists {
			result := handler(args)
			if result.typ == "error" {
				fmt.Fprintf(conn, "-%s\r\n", result.bulk)
			} else if result.typ == "string" {
				fmt.Fprintf(conn, "+%s\r\n", result.bulk)
			}
		} else {
			fmt.Fprintf(conn, "-ERR unknown command '%s'\r\n", command)
		}
	}
}