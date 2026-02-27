package main

import (
	"log"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Fail to start on port 8081:", err)
	}
}
