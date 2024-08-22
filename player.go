package main

import (
	"math/rand/v2"
	"time"
)

const PLAYER_UPDATES_PER_SECOND = 1

func makePlayer(account *Account) *Player {
	pl := new(Player)
	pl.account = account
	pl.lastUpdate = time.Now()
	pl.position.X = rand.Int() % 1024
	pl.position.Y = rand.Int() % 1024
	return pl
}
