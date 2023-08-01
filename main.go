package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/seew0/player-score-management-system/server"
)

func main() {

	engine := gin.Default()
	port := ":4000"

	Server, err := server.NewServer(port, engine)
	if err != nil {
		log.Printf("err occured while creating server: %v", err)
	}

	Server.Start()
}
