package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

func makePlayer(account *Account) *Player {
	pl := new(Player)
	pl.account = account
	pl.lastUpdate = time.Now()
	pl.position.X = float64(rand.Int() % 1024)
	pl.position.Y = float64(rand.Int() % 1024)
	return pl
}

func getPointDistance(a Point, b Point) float64 {
	dX := a.X - b.X
	dY := a.Y - b.Y
	dX *= dX
	dY *= dY
	return math.Sqrt(dX + dY)
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
		distance := getPointDistance(pos, other)
		if distance > PLAYER_UPDATE_DISTANCE {
			continue
		}
		clientSend(client, []byte(fmt.Sprintf("%s %s %.4f %.4f", CLIENT_FN_UPDATE, cl.account.username, cl.account.player.position.X, cl.account.player.position.Y)))
	}
	client.account.player.lastUpdate = time.Now()
}
