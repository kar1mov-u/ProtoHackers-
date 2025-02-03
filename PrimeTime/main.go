package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
)

type Input struct {
	Method string  `json:"method"`
	Number float64 `json:"number"`
}

type Resp struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

type malResp struct {
	Method string `json:"mithod"`
	Prime  bool   `json:"prim"`
}

const (
	CONN_TYPE = "tcp"
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "8000"
)

func main() {

	fmt.Println("Listening on port:", CONN_PORT)
	//start a server
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)

	// defer listener.Close() //close connection

	if err != nil {
		fmt.Println("Failed to start a server:", err.Error())
		os.Exit(1)
	}
	//Use loop to accept every request to the server
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to connect the client:", conn.RemoteAddr().String(), err.Error())
		}

		//Concurently hande each request
		go handleReq(conn)

	}

}

func handleReq(conn net.Conn) {
	buff := make([]byte, 4096) //create a buffer to hold values

	for {
		n, readErr := conn.Read(buff)
		if readErr != nil {
			fmt.Println("Failed to read", readErr.Error())
			conn.Close()
		}
		//Parse and validate Json
		resp, check := validateJson(buff, n)
		_, writeErr := conn.Write(resp)
		if writeErr != nil {
			fmt.Println("failed to write response", writeErr)
		}
		if !check {
			break
		}

	}
	defer conn.Close()
}

func floatToInt(f float64) (int, error) {
	if f != float64(int(f)) {
		return 1, errors.New("float is not a whole number")
	}
	return int(f), nil
}

func isPrimeNum(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func validateJson(body []byte, n int) ([]byte, bool) {
	var input Input
	//serialize the json to struct
	err := json.Unmarshal(body[:n], &input)
	if err != nil {
		fmt.Println("Failed to serailize Json", err)
	}
	num, _ := floatToInt(input.Number)

	if input.Method == "isPrime" && (num >= 1) {
		isPrime := isPrimeNum(num)
		resp := Resp{Method: "isPrime", Prime: isPrime}
		res, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("Failed to Serialize into json", err)
		}
		return res, true

	} else {
		resp := malResp{Method: "Nigga", Prime: false}
		res, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("Failed to Serialize into json", err)
		}
		return res, false
	}

}
