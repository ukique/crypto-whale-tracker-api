package main

import (
	"log"
	"net/http"

	"github.com/ukique/crypto-whale-tracker-api/internal/features/whale/transport/exchange"
)

func main() {
	exchange.ConnectWebSocket()
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Fail to start on port 8081:", err)
	}
}
