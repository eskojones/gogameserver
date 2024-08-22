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

var accounts = make(map[string]*Account)
var clients = make(map[string]*Client)
var messages []*ClientMessage

func broadcastBytes(msg []byte) {
	for _, v := range clients {
		_, _ = v.connection.Write(msg)
	}
}

func broadcastString(msg string) {
	broadcastBytes([]byte(msg))
}

func connectionHandler(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	fmt.Printf("[%s connected]\n", addr)
	broadcastString(fmt.Sprintf("[%s connected]\r\n", addr))
	client := new(Client)
	client.connection = conn
	client.history = make([]*ClientMessage, 0)
	client.lastRead = time.Now()
	clients[addr] = client

	defer conn.Close()
	readBuf := make([]byte, 1024)
	messageBuf := make([]byte, 1024)
	var bytesReadCount int
	for {
		_ = conn.SetReadDeadline(time.Now().Add((1000.0 / 20.0) * time.Millisecond))
		count, err := conn.Read(readBuf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else if !errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Printf("[read error: %s]\n", err)
				break
			}
		}

		clientSendUpdate(client)

		if count == 0 {
			if time.Now().Sub(client.lastRead) > 60*time.Second {
				fmt.Printf("[%s timed out]\n", addr)
				break
			}
			continue
		}

		messageBuf = fmt.Appendf(messageBuf[:bytesReadCount], "%s", readBuf[:count])
		bytesReadCount += count

		if bytesReadCount > 1024 {
			fmt.Printf("[%s sent an invalid message (too long)]\n", addr)
			break
		}

		if strings.Contains(string(readBuf), "\n") {
			client.lastRead = time.Now()
			clientMessageHandler(client, messageBuf[:], bytesReadCount)
			clear(messageBuf)
			bytesReadCount = 0
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

	loadClientFunctions()

	port, _ := strconv.Atoi(os.Args[1])
	err := listen(port, connectionHandler)
	if err != nil {
		log.Fatal(err)
		return
	}
}
