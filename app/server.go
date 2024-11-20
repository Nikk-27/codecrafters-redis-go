package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
)

/*
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
		go handleConnection(conn)
	}
}
*/


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

func main() {
    scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// Read input from Codecrafters
		input := scanner.Text()
		parts := strings.SplitN(input, " ", 2)

		// Parse command and arguments
		command := parts[0]
		var args []Value
		if len(parts) > 1 {
			args = []Value{{bulk: parts[1]}}
		}

		// Execute the command
		if cmdFunc, exists := Commands[command]; exists {
			result := cmdFunc(args)
			if result.typ == "error" {
				fmt.Println(result.str)
			} else {
				fmt.Println(result.bulk)
			}
		} else {
			fmt.Println("ERR unknown command")
		}
	}
}