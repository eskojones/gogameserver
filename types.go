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
	position   image.Point
	lastUpdate time.Time
}

const NET_TIMEOUT = 60       // time in seconds a client may go without sending a message
const NET_READ_DEADLINE = 50 // time in milliseconds to wait for a socket read
const NET_MSG_MAX_LEN = 1024 // maximum length of the client's read buffer (per message)

const CLIENT_FN_CREATE = "create" // command to create an account
const CLIENT_FN_LOGIN = "login"   // command to login to an account
const CLIENT_FN_LOGOUT = "logout" // command to logout from an account
const CLIENT_FN_UPDATE = "update" // command to update a player's position

const CLIENT_MSG_HISTORY_LEN = 50 // keep this many of the client's messages
const CLIENT_MSG_QUEUE_LEN = 500  // keep this many client messages in the global history

const PLAYER_UPDATE_PER_SECOND = 1  // player view updates per second
const PLAYER_UPDATE_DISTANCE = 1024 // distance a player is updated of another player

var accounts = make(map[string]*Account)                    // map of username -> Account
var clients = make(map[string]*Client)                      // map of remote address -> Client
var messages []*ClientMessage                               // global message history
var clientFunctions map[string]func(*Client, []string) bool // map of command -> client function
