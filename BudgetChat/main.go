package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"unicode"
)

const (
	CONN_TYPE = "tcp"
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "8000"
)

func main() {
	// Start server
	fmt.Println("Listening on : 8000")
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Failed to satrt server:", err.Error())
		os.Exit(1)
	}

	connections := map[string]net.Conn{}
	var mu sync.RWMutex

	for {
		//Accept connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed  to connect the client", err.Error())
		}

		go handleReq(conn, &mu, &connections)

	}

}

func handleReq(conn net.Conn, mu *sync.RWMutex, connections *map[string]net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))
	var name string
	var err error
	// var nameError bool

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Bytes()
		// name := strings.TrimSpace(string(line))
		//sets name if valid
		if len(name) == 0 {
			//get name
			name, err = setName(line)
			if err != nil {
				conn.Write([]byte("Invalid value to name"))
				// nameError = true
				break
			}
			fmt.Println("User " + name + " is connected")

			//Lock and read users
			userList := listUsers(connections, mu)
			conn.Write([]byte(userList + "\n"))

			//Lock and add user to users
			mu.Lock()
			(*connections)[name] = conn
			mu.Unlock()

			//send others that this user have joined
			go broadcast(name, connections, "* "+name+" has entered the room\n")

			continue
		}
		//if user is logged in broadcast its message
		// conn.Write([]byte(line))
		go broadcast(name, connections, "["+name+"] "+string(line)+"\n")

	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Erorr from reading connection:", err)
	} else if name != "" {
		//user terminated connection
		fmt.Println(name + "left the group")
		delete(*connections, name)
		go broadcast(name, connections, "* "+name+" has left the room\n")

	}

}

// --------------------------------------------- Helpers

func broadcast(author string, connections *map[string]net.Conn, message string) {

	for user := range *connections {
		if user == author {
			continue
		}
		(*connections)[user].Write([]byte(message))
	}
}

func setName(name []byte) (string, error) {
	if len(name) == 0 {
		return "", errors.New("Neme should be at least one char")
	}
	for i := 0; i < len(name); i++ {
		if unicode.IsDigit(rune(name[i])) || unicode.IsLetter(rune(name[i])) {
			continue
		} else {
			return "", errors.New("Name should only include alphanumeric values")
		}
	}
	return string(name), nil
}

func listUsers(connections *map[string]net.Conn, mu *sync.RWMutex) string {
	mu.RLock()
	defer mu.RUnlock()
	var res strings.Builder
	res.WriteString("*Current users: ")
	for user := range *connections {
		res.WriteString(user + " , ")
	}
	return res.String()
}
