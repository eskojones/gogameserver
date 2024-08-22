package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func makeClient(conn net.Conn) *Client {
	client := new(Client)
	client.connection = conn
	client.history = make([]*ClientMessage, 0)
	client.lastRead = time.Now()
	clients[conn.RemoteAddr().String()] = client
	return client
}

func deleteClient(client *Client) {
	delete(clients, client.connection.RemoteAddr().String())
}

func clientSend(client *Client, message []byte) bool {
	if client.connection == nil {
		return false
	}
	message = append(message, '\n')
	_, err := client.connection.Write(message)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func clientUpdate(client *Client) {
	if client.account == nil || client.connection == nil {
		return
	}

	playerUpdateView(client, false)
}

func makeClientMessage(client *Client, msgBuffer []byte, msgLength int) *ClientMessage {
	msg := new(ClientMessage)
	msg.sender = client
	msg.message = msgBuffer
	msg.length = msgLength
	msg.timestamp = time.Now()

	client.history = append(client.history, msg)
	if len(client.history) > CLIENT_MSG_HISTORY_LEN {
		client.history = client.history[1:]
	}

	messages = append(messages, msg)
	if len(messages) > CLIENT_MSG_QUEUE_LEN {
		messages = messages[1:]
	}
	fmt.Printf("<%s> %s", client.connection.RemoteAddr().String(), string(msg.message))
	return msg
}

func clientMessageHandler(client *Client, msgBuffer []byte, msgLength int) {
	msg := makeClientMessage(client, msgBuffer, msgLength)
	words := strings.Split(strings.ToLower(string(msg.message)), " ")
	if len(words) == 0 {
		return
	}
	for i := range words {
		for strings.Contains(words[i], "\n") {
			words[i] = strings.ReplaceAll(words[i], "\n", "")
		}
	}
	fn := clientFunctions[words[0]]
	if fn == nil {
		onClientInvalidMessage(client, words)
		return
	}
	fn(msg.sender, words)
}

func onClientInvalidMessage(client *Client, args []string) {
	fmt.Printf("%s sent an invalid message (%s)!\n", client.connection.RemoteAddr().String(), args[0])
	clientSend(client, []byte("invalid message"))
}
