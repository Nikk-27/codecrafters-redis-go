package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
	"strconv"
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

// Parse RESP input
func parseRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)

	// Check for an array start (e.g., "*2")
	if !strings.HasPrefix(line, "*") {
		return nil, fmt.Errorf("invalid RESP format: %s", line)
	}

	// Parse the number of arguments
	numArgs, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid argument count: %s", line[1:])
	}

	// Read the arguments
	args := make([]string, numArgs)
	for i := 0; i < numArgs; i++ {
		// Read the bulk string indicator (e.g., "$4")
		sizeLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		sizeLine = strings.TrimSpace(sizeLine)
		if !strings.HasPrefix(sizeLine, "$") {
			return nil, fmt.Errorf("invalid bulk string format: %s", sizeLine)
		}

		// Parse the bulk string length
		argSize, err := strconv.Atoi(sizeLine[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid bulk string size: %s", sizeLine[1:])
		}

		// Read the bulk string itself
		arg := make([]byte, argSize)
		_, err = reader.Read(arg)
		if err != nil {
			return nil, err
		}
		args[i] = string(arg)

		// Read and discard the trailing \r\n
		_, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
	}
	return args, nil
}

// Handle ECHO command execution
func echoExecute(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Parse RESP input
		args, err := parseRESP(reader)
		fmt.Println("Hi")
		if err != nil {
			fmt.Fprintf(conn, "-ERR %s\r\n", err.Error())
			return
		}

		// Check if the command exists
		command := strings.ToUpper(args[0])
		fmt.Println(command)
		if handler, exists := Commands[command]; exists {
			// Execute the command
			commandArgs := make([]Value, len(args[1:]))
			for i, arg := range args[1:] {
				fmt.Println(arg)
				commandArgs[i] = Value{typ: "string", bulk: arg}
				fmt.Println("commandArgs", commandArgs)
			}

			result := handler(commandArgs)
			fmt.Println("result", result)
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