package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type ClientMessage struct {
	sender    Client
	message   []byte
	timestamp time.Time
}

type Client struct {
	connection net.Conn
	history    []ClientMessage
	lastRead   time.Time
}

type Account struct {
	username string
	password string
	client   *Client
}

var accounts = make(map[string]*Account)
var clients = make(map[string]Client)
var messages []ClientMessage

func broadcastBytes(msg []byte) {
	for _, v := range clients {
		_, _ = v.connection.Write(msg)
	}
}

func broadcastString(msg string) {
	broadcastBytes([]byte(msg))
}

func accountCreate(username string, password string) bool {
	account := accounts[username]
	if account == nil {
		account = new(Account)
		account.username = username
		account.password = password
		accounts[username] = account
		return true
	}
	return false
}

func clientSend(client Client, message []byte) bool {
	message = append(message, '\n')
	_, err := client.connection.Write(message)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func messageHandler(msg ClientMessage) bool {
	words := strings.Split(strings.ToLower(string(msg.message)), " ")
	if len(words) == 0 {
		return false
	}
	command := words[0]
	switch command {
	case "create":
		// account create
		ret := accountCreate(words[1], words[2])
		if ret == false {
			clientSend(msg.sender, []byte("false"))
		} else {
			clientSend(msg.sender, []byte("true"))
		}
	case "auth":
		// account auth
	case "pos":
		// position
	default:
		// invalid message
		fmt.Printf("%s sent an invalid message!\n", msg.sender.connection.RemoteAddr().String())
	}

	return true
}

func connHandler(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	fmt.Printf("[%s connected]\n", addr)
	broadcastString(fmt.Sprintf("[%s connected]\r\n", addr))
	client := Client{
		connection: conn,
		history:    make([]ClientMessage, 0),
		lastRead:   time.Now(),
	}
	clients[addr] = client

	defer conn.Close()
	readBuf := make([]byte, 1024)
	messageBuf := make([]byte, 1024)

	for {
		count, err := conn.Read(readBuf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("[read error: %s]\n", err)
			}
			break
		}

		if count == 0 {
			if time.Now().Sub(client.lastRead) > 60*time.Second {
				fmt.Printf("[%s timed out]\n", addr)
				break
			}
			continue
		}

		if len(messageBuf)+len(readBuf) > 1024 {
			fmt.Printf("[%s sent an invalid message]\n", addr)
			break
		}
		messageBuf = fmt.Appendf(messageBuf, "%s", readBuf)
		if strings.Contains(string(readBuf), "\n") {
			client.lastRead = time.Now()
			clientMessage := ClientMessage{
				sender:    client,
				message:   messageBuf,
				timestamp: time.Now(),
			}
			client.history = append(client.history, clientMessage)
			if len(client.history) > 64 {
				client.history = client.history[1:]
			}
			messages = append(messages, clientMessage)
			if len(messages) > 256 {
				messages = messages[1:]
			}
			tmp := string(messageBuf)
			fmt.Printf("<%s> %s", addr, tmp)
			// broadcastString(fmt.Sprintf("<%s> %s", addr, tmp))
			clear(messageBuf)
			messageHandler(clientMessage)
		}
		clear(readBuf)
	}
	fmt.Printf("[%s disconnected]\n", addr)
	broadcastString(fmt.Sprintf("[%s disconnected]\r\n", addr))
	delete(clients, addr)
}

func listen(port int, handler func(net.Conn)) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handler(conn)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: gogameserver <port>\n")
		return
	}
	port, _ := strconv.Atoi(os.Args[1])
	err := listen(port, connHandler)
	if err != nil {
		log.Fatal(err)
		return
	}
}
