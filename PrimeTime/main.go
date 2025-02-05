package main

import (
	"bufio"
	"encoding/json"
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
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
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

	if err != nil {
		fmt.Println("Failed to start a server:", err.Error())
		os.Exit(1)
	}
	defer listener.Close() //close connection

	//Use loop to accept every request to the server
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to connect the client:", err.Error())
			continue
		}

		//Concurently hande each request
		go handleReq(conn)
	}
}

func handleReq(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Bytes()
		resp, check := validateJson(line, len(line))

		_, err := conn.Write(append(resp, '\n'))
		fmt.Println(" response is " + string(resp) + " request was " + string(line[:]))
		if err != nil {
			fmt.Println("failed to write response: ", err)
			return
		}

		if !check {
			fmt.Println("request was malformed, terminated connection", conn.RemoteAddr().String())
			return
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erorr from reading connection:", err)
	}

}

func floatToInt(f float64) (int, error) {

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
	input := Input{Number: -10100}

	//serialize the json to struct
	err := json.Unmarshal(body[:n], &input)
	//return malresp if cannot unmarshal

	if err != nil {
		fmt.Println("Failed to serailize Json", err)
		r, _ := json.Marshal(malResp{Method: "error", Prime: false})
		return r, false
	}

	num, err := floatToInt(input.Number)
	if err != nil {
		fmt.Println("number is not integer:", err)
		r, _ := json.Marshal(malResp{Method: "error", Prime: false})
		return r, false
	}

	//check if its correct
	if input.Method == "isPrime" && input.Number != -10100 {
		isPrime := isPrimeNum(num)
		resp := Resp{Method: "isPrime", Prime: isPrime}
		res, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("Failed to Serialize into json", err)
		}
		return res, true

	}
	res, _ := json.Marshal(malResp{Method: "error", Prime: false})
	return res, false
}
