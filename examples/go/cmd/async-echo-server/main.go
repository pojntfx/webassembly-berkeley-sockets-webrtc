package main

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/pojntfx/webassembly-berkeley-sockets-via-webrtc/examples/go/pkg/sockets"
)

var (
	LOCAL_HOST = []byte{10, 0, 0, 240}
)

const (
	LOCAL_PORT = 1234
	BACKLOG    = 1

	BUFFER_LENGTH = 1024
)

func main() {
	// Create address
	serverAddress := sockets.SockaddrIn{
		SinFamily: sockets.PF_INET,
		SinPort:   sockets.Htons(LOCAL_PORT),
		SinAddr: struct{ SAddr uint32 }{
			SAddr: uint32(binary.LittleEndian.Uint32(LOCAL_HOST)),
		},
	}
	serverAddressReadable := fmt.Sprintf("%v:%v", LOCAL_HOST, LOCAL_PORT)

	// Create socket
	serverSocket := sockets.Socket(sockets.PF_INET, sockets.SOCK_STREAM, 0)
	if serverSocket == -1 {
		log.Fatalf("[ERROR] Could not create socket %v: %v\n", serverAddressReadable, serverSocket)
	}

	// Bind
	if err := sockets.Bind(serverSocket, &serverAddress); err == -1 {
		log.Fatalf("[ERROR] Could not bind socket %v: %v\n", serverAddressReadable, err)
	}

	// Listen
	if err := sockets.Listen(serverSocket, BACKLOG); err == -1 {
		log.Fatalf("[ERROR] Could not listen on socket %v: %v\n", serverAddressReadable, err)
	}

	log.Println("[INFO] Listening on", serverAddressReadable)

	// Accept loop
	for {
		log.Println("[DEBUG] Accepting on", serverAddressReadable)

		clientAddress := sockets.SockaddrIn{}

		// Accept
		clientSocket := sockets.Accept(serverSocket, &clientAddress)
		if clientSocket == -1 {
			log.Println("[ERROR] Could not accept, continuing:", clientSocket)

			continue
		}

		go func(innerClientSocket int32, innerClientAddress sockets.SockaddrIn) {
			clientHost := make([]byte, 4) // xxx.xxx.xxx.xxx
			binary.LittleEndian.PutUint32(clientHost, uint32(innerClientAddress.SinAddr.SAddr))

			clientAddressReadable := fmt.Sprintf("%v:%v", clientHost, innerClientAddress.SinPort)

			log.Println("[INFO] Accepted client", clientAddressReadable)

			// Receive loop
			for {
				log.Printf("[DEBUG] Waiting for client %v to send\n", clientAddressReadable)

				// Receive
				receivedMessage := make([]byte, BUFFER_LENGTH)

				receivedMessageLength := sockets.Recv(innerClientSocket, &receivedMessage, BUFFER_LENGTH, 0)
				if receivedMessageLength == -1 {
					log.Printf("[ERROR] Could not receive from client %v, dropping message: %v\n", clientAddressReadable, receivedMessageLength)

					continue
				}

				if receivedMessageLength == 0 {
					break
				}

				log.Printf("[DEBUG] Received %v bytes from %v\n", receivedMessageLength, clientAddressReadable)

				// Send
				sentMessage := []byte(fmt.Sprintf("You've sent: %v", string(receivedMessage))) // TODO: Access the received message here

				sentMessageLength := sockets.Send(innerClientSocket, sentMessage, 0)
				if sentMessageLength == -1 {
					log.Printf("[ERROR] Could not send to client %v, dropping message: %v\n", clientAddressReadable, sentMessageLength)

					return
				}

				log.Printf("[DEBUG] Sent %v bytes to %v\n", sentMessageLength, clientAddressReadable)
			}
		}(clientSocket, clientAddress)
	}
}