package main

import (
	"fmt"
	"net"
	"os"
)

type Value struct {
    typ  string
    str  string
    bulk string
}

var Commands = map[string]func([]Value) Value{"ECHO":echo,}

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
		//go handleConnection(conn)
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



func echo(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'echo' command"}
	}
	return Value{typ: "string", str: args[0].bulk}
}