package main

import (
	"fmt"
	"net"
	"os"
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
	reader := bufio.NewReader(os.Stdin) // Read input from Codecrafters
	for {
		// Read a line of input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("ERR reading input:", err)
			continue
		}

		// Trim and split the input into command and arguments
		input = strings.TrimSpace(input)
		parts := strings.Split(input, " ")

		if len(parts) < 1 {
			fmt.Println("ERR empty command")
			continue
		}

		command := strings.ToUpper(parts[0])
		args := make([]app.Value, len(parts[1:]))
		for i, arg := range parts[1:] {
			args[i] = app.Value{typ: "string", bulk: arg}
		}

		// Find and execute the command
		handler, exists := app.Commands[command]
		if !exists {
			fmt.Println("ERR unknown command:", command)
			continue
		}

		// Get the result and output it
		result := handler(args)
		if result.typ == "error" {
			fmt.Println(result.str)
		} else if result.typ == "string" {
			fmt.Println(result.str)
		}
	}
}
