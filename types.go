package main

import (
	"image"
	"net"
	"time"
)

type ClientMessage struct {
	sender    *Client
	message   []byte
	length    int
	timestamp time.Time
}

type Client struct {
	connection net.Conn
	history    []*ClientMessage
	lastRead   time.Time
	account    *Account
}

type Account struct {
	username string
	password string
	client   *Client
	player   *Player
}

type Player struct {
	account  *Account
	position image.Point
}
