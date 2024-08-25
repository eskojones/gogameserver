package main

const NET_TIMEOUT = 60       // time in seconds a client may go without sending a message
const NET_READ_DEADLINE = 5  // time in milliseconds to wait for a socket read
const NET_MSG_MAX_LEN = 1024 // maximum length of the client's read buffer (per message)

const CLIENT_FN_CREATE = "create" // command to create an account
const CLIENT_FN_LOGIN = "login"   // command to login to an account
const CLIENT_FN_LOGOUT = "logout" // command to logout from an account
const CLIENT_FN_UPDATE = "update" // command to update a player's position
const CLIENT_FN_QUERY = "query"

const CLIENT_MSG_HISTORY_LEN = 50 // keep this many of the client's messages
const CLIENT_MSG_QUEUE_LEN = 500  // keep this many client messages in the global history

const PLAYER_UPDATE_PER_SECOND = 20 // player view updates per second
const PLAYER_UPDATE_DISTANCE = 1024 // distance a player is updated of another player
