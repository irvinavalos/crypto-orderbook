package coinbase_api

import (
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

const WebSocketURL = "wss://advanced-trade-ws.coinbase.com"
const BTC_USD = "BTC-USD"

type SubscribeMessage struct {
	Type       string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channel    string   `json:"channel"`
}

func NewSubscribeMessage(productIDs []string) SubscribeMessage {
	return SubscribeMessage{
		Type:       "subscribe",
		ProductIDs: productIDs,
		Channel:    "level2",
	}
}

func (sm *SubscribeMessage) stringProductIDs() string {
	var s strings.Builder
	for i := range sm.ProductIDs {
		s.WriteString(sm.ProductIDs[i])
	}
	return s.String()
}

func StartCoinbaseWS() {
	conn, _, _ := websocket.DefaultDialer.Dial(WebSocketURL, nil)
	defer conn.Close()

	log.Println("Connecting to Coinbase Advanced Trade WebSocket")

	// productIDs := makeProductIDs(BTC_USD)
	productIDs := []string{"BTC-USD"}
	subscriptionMessage := NewSubscribeMessage(productIDs)

	if err := conn.WriteJSON(subscriptionMessage); err != nil {
		log.Fatalf("Error with subscription message: %v\n", err)
	}

	log.Printf("Subscribed to [%s] level2\n", subscriptionMessage.stringProductIDs())

	var coinbaseMessage *CoinbaseMessage

	for {
		if err := conn.ReadJSON(&coinbaseMessage); err != nil {
			log.Printf("Error reading JSON: %v\n", err)
			return
		}

		log.Println(coinbaseMessage)
	}
}
