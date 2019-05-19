package main

import (
	"github.com/jonas747/yams/connection"
	"log"
)

func main() {
	manager := &connection.ConnectionManager{}
	err := manager.Listen(":25565")
	if err != nil {
		log.Fatal("Failed listening for connection: ", err)
	}
}
