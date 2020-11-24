package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// DecodeJSONSumResult decodes JSON sum results
type DecodeJSONSumResult struct {
	SumResult []float64 `json:"sumResult"`
	MyCount   int       `json:"myCount"`
}

type EncodeJSONSumInput struct {
	InputArray []float64 `json:"inputArray"`
	IonCount   int       `json:"ionCount"`
	MyCount    int       `json:"myCount"`
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:3333")
	checkError(err)

	ln, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := ln.Accept()
		checkError(err)

		go handleConnection(conn.(*net.TCPConn))
	}
}

func handleConnection(conn *net.TCPConn) {
	var input [512]byte
	var jsonSumInput = []float64{1, 1, 3}

	n, err := conn.Read(input[0:])
	checkError(err)

	fmt.Println(string(input[0:n]))

	bytes := encodeJSONSumInput(EncodeJSONSumInput{jsonSumInput, 3, 0})

	_, err = conn.Write(bytes)
	checkError(err)

	n, err = conn.Read(input[0:])
	checkError(err)

	a := decodeJSONSumResult(string(input[0:n]))

	fmt.Println(a.SumResult)
	fmt.Println(a.MyCount)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func decodeJSONSumResult(input string) DecodeJSONSumResult {
	rawIn := json.RawMessage(input)

	bytes, err := rawIn.MarshalJSON()
	checkError(err)

	var d DecodeJSONSumResult

	err = json.Unmarshal(bytes, &d)
	checkError(err)

	return d
}

func encodeJSONSumInput(s EncodeJSONSumInput) []byte {

	bytes, err := json.Marshal(s)
	checkError(err)

	return bytes
}

func encodeJSONResult() {}

func decodeJSONResult() {}
