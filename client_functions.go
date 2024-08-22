package main

import "fmt"

const CLIENT_FN_CREATE = "create"
const CLIENT_FN_LOGIN = "login"
const CLIENT_FN_LOGOUT = "logout"

var clientFunctions map[string]func(*Client, []string) bool

func loadClientFunctions() {
	clientFunctions = make(map[string]func(*Client, []string) bool)
	clientFunctions[CLIENT_FN_CREATE] = fnAccountCreate
	clientFunctions[CLIENT_FN_LOGIN] = fnAccountLogin
	clientFunctions[CLIENT_FN_LOGOUT] = fnAccountLogout
}

// handles client account-create command
func fnAccountCreate(client *Client, args []string) bool {
	if len(args) != 3 || client.account != nil {
		clientSend(client, []byte(fmt.Sprintf("%s false", CLIENT_FN_CREATE)))
		fmt.Printf("[account create failed, wrong arg count (%d)]\n", len(args))
		return false
	}
	username := args[1]
	password := args[2]
	account := accounts[username]
	if account != nil {
		clientSend(client, []byte(fmt.Sprintf("%s false", CLIENT_FN_CREATE)))
		fmt.Printf("[account create failed, '%s' already exists]\n", username)
		return false
	}
	account = new(Account)
	account.username = username
	account.password = password
	accounts[username] = account
	clientSend(client, []byte(fmt.Sprintf("%s true", CLIENT_FN_CREATE)))
	account.player = makePlayer(account)
	fmt.Printf("[account create success (%s)]\n", account.username)
	return true
}

// handles client account-login command
func fnAccountLogin(client *Client, args []string) bool {
	if len(args) != 3 {
		clientSend(client, []byte(fmt.Sprintf("%s false", CLIENT_FN_LOGIN)))
		fmt.Printf("[account login failed, wrong arg count (%d)]\n", len(args))
		return false
	}
	username := args[1]
	password := args[2]
	if client.account != nil || accounts[username] == nil || accounts[username].password != password {
		clientSend(client, []byte(fmt.Sprintf("%s false", CLIENT_FN_LOGIN)))
		fmt.Printf("[account login failed, already logged in or bad credentials]\n")
		return false
	}
	accounts[username].client = client
	client.account = accounts[username]
	clientSend(client, []byte(fmt.Sprintf("%s true", CLIENT_FN_LOGIN)))
	fmt.Printf("[account login success (%s)]\n", username)
	return true
}

// handles client account-logout command
func fnAccountLogout(client *Client, args []string) bool {
	if len(args) != 1 || client.account == nil {
		clientSend(client, []byte(fmt.Sprintf("%s false", CLIENT_FN_LOGOUT)))
		fmt.Printf("[account logout failed]\n")
		return false
	}
	accounts[client.account.username].client = nil
	fmt.Printf("[account logout success (%s)]\n", client.account.username)
	client.account = nil
	clientSend(client, []byte(fmt.Sprintf("%s true", CLIENT_FN_LOGOUT)))
	return true
}
