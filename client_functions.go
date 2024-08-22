package main

var clientFunctions map[string]func(*Client, []string) bool

func loadClientFunctions() {
	clientFunctions = make(map[string]func(*Client, []string) bool)
	clientFunctions["create"] = fnAccountCreate
	clientFunctions["login"] = fnAccountLogin
	clientFunctions["logout"] = fnAccountLogout
}

// handles client account-create command
func fnAccountCreate(client *Client, args []string) bool {
	if len(args) != 3 || client.account != nil {
		clientSend(client, []byte("create false"))
		return false
	}
	username := args[1]
	password := args[2]
	account := accounts[username]
	if account != nil {
		clientSend(client, []byte("create false"))
		return false
	}
	account = new(Account)
	account.username = username
	account.password = password
	accounts[username] = account
	clientSend(client, []byte("create true"))
	return true
}

// handles client account-login command
func fnAccountLogin(client *Client, args []string) bool {
	if len(args) != 3 {
		clientSend(client, []byte("login false"))
		return false
	}
	username := args[1]
	password := args[2]
	if client.account != nil || accounts[username] == nil || accounts[username].password != password {
		clientSend(client, []byte("login false"))
		return false
	}
	accounts[username].client = client
	client.account = accounts[username]
	clientSend(client, []byte("login true"))
	return true
}

// handles client account-logout command
func fnAccountLogout(client *Client, args []string) bool {
	if len(args) != 1 || client.account == nil {
		clientSend(client, []byte("logout false"))
		return false
	}
	accounts[client.account.username].client = nil
	client.account = nil
	clientSend(client, []byte("logout true"))
	return true
}
