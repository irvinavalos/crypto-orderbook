package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// https://stackoverflow.com/questions/62652236/can-i-have-a-function-to-check-if-a-key-is-in-a-map
func containsKey[M ~map[K]V, K comparable, V any](m M, k K) bool {
	_, ok := m[k]
	return ok
}

func stringToInt64(s string) int64 {
	var mult int64 = 100_000_000

	if !strings.Contains(s, ".") {
		res, _ := strconv.ParseInt(s, 10, 64)
		return res * mult
	}

	strSlice := strings.Split(s, `.`)
	strLeft, strRight := strSlice[0], strSlice[1]
	intLeft, _ := strconv.ParseInt(strLeft, 10, 64)
	res := intLeft * mult

	if len(strRight) > 0 {
		if len(strRight) > 8 {
			strRight = strRight[:8]
		}

		padding := strings.Repeat("0", 8-len(strRight))
		intRight, _ := strconv.ParseInt(strRight+padding, 10, 64)
		res += intRight
	}

	return res
}

func int64ToString(i int64) string {
	var mult int64 = 100_000_000

	strLeft := strconv.FormatInt(i/mult, 10)
	strRight := strconv.FormatInt(i%mult, 10)

	padding := 8 - len(strRight)

	if padding > 0 {
		strRight = strings.Repeat("0", padding) + strRight
	}

	return strLeft + "." + strRight
}

const wsURL = "wss://advanced-trade-ws.coinbase.com"

type SubscribeMessage struct {
	Type       string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channel    string   `json:"channel"`
}

func makeSubscribeMessage(productIDs []string) SubscribeMessage {
	return SubscribeMessage{Type: "subscribe", ProductIDs: productIDs, Channel: "level2"}
}

func stringifyProductIDs(productIDs []string) string {
	var s strings.Builder
	for i := range productIDs {
		s.WriteString(productIDs[i])
	}
	return s.String()
}

type Update struct {
	Side        string `json:"side"`
	EventTime   string `json:"event_time"`
	PriceLevel  string `json:"price_level"`
	NewQuantity string `json:"new_quantity"`
}

type Event struct {
	Type      string   `json:"type"`
	ProductID string   `json:"product_id"`
	Updates   []Update `json:"updates"`
}

type CoinbaseResponse struct {
	Channel        string  `json:"channel"`
	Timestamp      string  `json:"timestamp"`
	SequenceNumber int     `json:"sequence_num"`
	Events         []Event `json:"events"`
}

type Orderbook struct {
	Bids         map[int64]int64
	Offers       map[int64]int64
	SortedBids   []int64
	SortedOffers []int64
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Bids:         make(map[int64]int64),
		Offers:       make(map[int64]int64),
		SortedBids:   make([]int64, 0, 1000),
		SortedOffers: make([]int64, 0, 1000),
	}
}

func (u *Update) getPrice() int64 {
	return stringToInt64(u.PriceLevel)
}

func (u *Update) getSize() int64 {
	return stringToInt64(u.NewQuantity)
}

func main() {
	conn, _, _ := websocket.DefaultDialer.Dial(wsURL, nil)
	defer conn.Close()

	log.Println("Connecting to Coinbase Advanced Trade WebSocket")

	// productIDs := []string{"ETH-USD"}
	productIDs := []string{"BTC-USD"}
	subMsg := makeSubscribeMessage(productIDs)

	if err := conn.WriteJSON(subMsg); err != nil {
		log.Fatal("Subscript Error:", err)
	}

	log.Printf("Subscribed to [%s] level2\n", stringifyProductIDs(productIDs))

	var coinbaseResponse *CoinbaseResponse

	for {
		// _, msg, err := conn.ReadMessage()
		// if err != nil {
		// 	log.Println("Read Error:", err)
		// }
		// log.Println(string(msg))
		if err := conn.ReadJSON(&coinbaseResponse); err != nil {
			log.Println("Read Error:", err)
			return
		}
		log.Println(coinbaseResponse)
	}
}
