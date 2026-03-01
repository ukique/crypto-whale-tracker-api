package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ukique/crypto-whale-tracker-api/internal/features/whale/models"
	"github.com/ukique/crypto-whale-tracker-api/internal/features/whale/transport/exchange"
	"github.com/ukique/crypto-whale-tracker-api/internal/features/whale/transport/handlers"
)

func main() {
	r := gin.Default()

	whaleChan := make(chan models.Whale)
	go exchange.ConnectWebSocket(whaleChan)

	r.GET("/whale", handlers.WhaleHandler(whaleChan))
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Fail to start on port 8081:", err)
	}
}
