package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

const PLAYER_UPDATE_PER_SECOND = 1
const PLAYER_UPDATE_DISTANCE = 10

func makePlayer(account *Account) *Player {
	pl := new(Player)
	pl.account = account
	pl.lastUpdate = time.Now()
	pl.position.X = rand.Int() % 1024
	pl.position.Y = rand.Int() % 1024
	return pl
}

func playerUpdateView(client *Client, force bool) {
	if !force {
		if time.Now().Sub(client.account.player.lastUpdate).Milliseconds() < 1000/PLAYER_UPDATE_PER_SECOND*time.Millisecond.Milliseconds() {
			return
		}
	}

	pos := client.account.player.position

	for _, cl := range clients {
		if cl.account == nil || cl.connection == nil {
			continue
		}
		other := cl.account.player.position
		dX := pos.X - other.X
		dX *= dX
		if dX > PLAYER_UPDATE_DISTANCE {
			continue
		}
		dY := pos.Y - other.Y
		if dY > PLAYER_UPDATE_DISTANCE {
			continue
		}
		dY *= dY
		distance := math.Sqrt(float64(dX + dY))
		if distance > PLAYER_UPDATE_DISTANCE {
			continue
		}
		clientSend(client, []byte(fmt.Sprintf("update %s %d %d", cl.account.username, cl.account.player.position.X, cl.account.player.position.Y)))
	}
	client.account.player.lastUpdate = time.Now()
}
