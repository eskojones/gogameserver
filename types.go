package main

import (
	"image"
	"net"
	"sync"
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
	account    *Account
	lastRead   time.Time
}

type Account struct {
	username string
	password string
	client   *Client
	player   *Player
}

type Player struct {
	account    *Account
	position   Point
	sprite     []image.Point
	lastUpdate time.Time
}

type Point struct {
	X float64
	Y float64
}

var accounts = make(map[string]*Account)                    // map of username -> Account
var clients = make(map[string]*Client)                      // map of remote address -> Client
var messages []*ClientMessage                               // global message history
var clientFunctions map[string]func(*Client, []string) bool // map of command -> client function
var entityMutex sync.Mutex
