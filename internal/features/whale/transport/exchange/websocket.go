package exchange

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ukique/crypto-whale-tracker-api/internal/features/whale/models"
	"github.com/ukique/crypto-whale-tracker-api/internal/features/whale/service"
	"github.com/ukique/crypto-whale-tracker-api/internal/pb"
	"google.golang.org/protobuf/proto"
)

const (
	WebsocketURL = "wss://wbs-api.mexc.com/ws"
)

func ConnectWebSocket(whaleChan chan models.Whale) {
	//graceful shutdown when we use Ctrl + C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial(WebsocketURL, nil)
	if err != nil {
		log.Fatal("fail dial connection:", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("fail close connection", err)
		}
	}()

	done := make(chan struct{})

	subscriptionMessage := map[string]interface{}{
		"method": "SUBSCRIPTION",
		"params": []string{"spot@public.aggre.deals.v3.api.pb@100ms@BTCUSDT"},
	}
	subMsg, _ := json.Marshal(subscriptionMessage)

	//if there is no valid subscription on the websocket,
	//the MEXC server will actively disconnect after 30 second
	//so we need send PING
	pingMessage := map[string]string{
		"method": "PING",
	}
	pingMsg, _ := json.Marshal(pingMessage)

	if err = conn.WriteMessage(websocket.TextMessage, subMsg); err != nil {
		log.Println("fail to write subscription message:", err)
		return
	}
	log.Printf("Send subscription message: %s", subMsg)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	go func() {
		defer close(done)
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.TextMessage, pingMsg); err != nil {
					log.Println("fail write PING message:", err)
					return
				}
				log.Println("send PING message:", pingMsg)

			case <-interrupt:
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("fail write close message:", err)
					return
				}
				return
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("fail to read message:", err)
			return
		}

		// Skip system MEXC notifications (subscription confirmations, PONG)
		if len(message) > 0 && message[0] == '{' {
			continue
		}
		decodedMessage := &pb.PushDataV3ApiWrapper{}
		if err := proto.Unmarshal(message, decodedMessage); err != nil {
			log.Println("protobuf unmarshal error:", err)
			continue
		}

		for _, deal := range decodedMessage.GetPublicAggreDeals().Deals {
			if service.IsWhale(deal.Price, deal.Quantity) {
				whaleChan <- models.Whale{
					Price:    deal.Price,
					Quantity: deal.Quantity,
					Symbol:   decodedMessage.GetSymbol(),
					Time:     deal.Time,
				}
			}
		}
	}
}
