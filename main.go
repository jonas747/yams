package main

import (
	"github.com/jonas747/yams/connection"
	"github.com/jonas747/yams/game"
	"log"
)

func main() {
	game := game.NewGame()

	connManager := &connection.ConnectionManager{
		Upgrader: game.Upgrader,
	}

	err := connManager.Listen(":25565")
	if err != nil {
		log.Fatal("Failed listening for connection: ", err)
	}
}
