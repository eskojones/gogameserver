package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func clientSend(client *Client, message []byte) bool {
	message = append(message, '\n')
	_, err := client.connection.Write(message)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func clientSendUpdate(client *Client) {
	if client.account == nil {
		return
	}

	// send position update to client
	clientSend(client, []byte("update"))
}

func clientMessageHandler(client *Client, msgBuffer []byte, msgLength int) {
	msg := new(ClientMessage)
	msg.sender = client
	msg.message = msgBuffer
	msg.length = msgLength
	msg.timestamp = time.Now()

	client.history = append(client.history, msg)
	if len(client.history) > 64 {
		client.history = client.history[1:]
	}
	messages = append(messages, msg)
	if len(messages) > 256 {
		messages = messages[1:]
	}
	fmt.Printf("<%s> %s", client.connection.RemoteAddr().String(), string(msg.message))

	words := strings.Split(strings.ToLower(string(msg.message)), " ")
	if len(words) == 0 {
		return
	}
	for i := range words {
		for strings.Contains(words[i], "\n") {
			words[i] = strings.ReplaceAll(words[i], "\n", "")
		}
		// fmt.Printf("\"%s\"(%d) ", words[i], len(words[i]))
	}
	// fmt.Printf("\n")

	fn := clientFunctions[words[0]]
	if fn == nil {
		fmt.Printf("%s sent an invalid message (%s)!\n", msg.sender.connection.RemoteAddr().String(), words[0])
		clientSend(msg.sender, []byte("invalid message"))
		return
	}
	fn(msg.sender, words)
}
