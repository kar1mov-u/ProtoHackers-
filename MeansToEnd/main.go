package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	CONN_TYPE = "tcp"
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "8000"
)

func main() {
	//create tcp listener
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Failed to start server: ", err.Error())
	}
	defer listener.Close()

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed initatie connection")
			continue
		}

		go handleFunc(conn)
	}
}

func handleFunc(conn net.Conn) {
	data := [][]int32{}
	defer conn.Close()
	buf := make([]byte, 9)

	for {
		_, err := io.ReadFull(conn, buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
			} else {
				fmt.Println("Error: ", err)
			}
			break
		}
		// fmt.Println(string(buf))

		//proces mesage
		form, val1, val2 := procesMessage(buf)
		if form == 'I' {
			data = append(data, []int32{val1, val2})

		} else if form == 'Q' {
			var sum int64
			var count int64
			var mean int32
			for i := 0; i < len(data); i++ {
				if val1 <= data[i][0] && data[i][0] <= val2 {
					sum += int64(data[i][1])
					count += 1
				}
			}
			if count == 0 {
				mean = 0
			} else {
				mean = int32(float64(sum)/float64(count) + 0.5)
			}
			var response [4]byte
			binary.BigEndian.PutUint32(response[:], uint32(mean))
			conn.Write(response[:])
			fmt.Println(mean)
			break
		}
	}
}

func procesMessage(data []byte) (rune, int32, int32) {
	if len(data) != 9 {
		fmt.Println("not")
	}
	messageType := data[0]
	// Convert bytes to int32 using big-endian order
	val1 := int32(binary.BigEndian.Uint32(data[1:5])) // First int32 (timestamp or min time)
	val2 := int32(binary.BigEndian.Uint32(data[5:9])) // Second int32 (price or max time)
	fmt.Printf("Received message: Type=%c, Value1=%d, Value2=%d\n", messageType, val1, val2)

	return rune(messageType), val1, val2
}
