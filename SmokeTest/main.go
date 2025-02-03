package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections
	fmt.Println("Listening on:", CONN_PORT)
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	for {
		// Accept incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		// Handle connection in a new goroutine
		go handleReq(conn)
	}
}

// Handle Incoming request
func handleReq(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Serving client: %s\n", conn.RemoteAddr().String())

	buf := make([]byte, 4096) // Buffer for reading data

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected:", conn.RemoteAddr().String())
			} else {
				fmt.Println("Read error:", err)
			}
			break
		}

		// Send the received message back (echo)
		_, writeErr := conn.Write(buf[:n])
		if writeErr != nil {
			fmt.Println("Write error:", writeErr)
			break
		}

		fmt.Printf("Echoed message from %s: %s\n", conn.RemoteAddr().String(), string(buf[:n]))
	}
}
