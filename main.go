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
	client := makeClient(conn)
	defer conn.Close()
	readBuf := make([]byte, NET_MSG_MAX_LEN)
	messageBuf := make([]byte, NET_MSG_MAX_LEN)
	var bytesReadCount int
	for {
		_ = conn.SetReadDeadline(time.Now().Add(NET_READ_DEADLINE * time.Millisecond))
		count, err := conn.Read(readBuf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else if !errors.Is(err, os.ErrDeadlineExceeded) {
				fmt.Printf("[read error: %s]\n", err)
				break
			}
		}

		clientUpdate(client)

		if count == 0 {
			if time.Now().Sub(client.lastRead) > NET_TIMEOUT*time.Second {
				fmt.Printf("[%s timed out]\n", addr)
				break
			}
			continue
		}

		messageBuf = fmt.Appendf(messageBuf[:bytesReadCount], "%s", readBuf[:count])
		bytesReadCount += count

		if bytesReadCount > NET_MSG_MAX_LEN {
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
	deleteClient(client)
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
