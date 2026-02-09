package coinbase_api

import (
	"strconv"
	"strings"

	"github.com/irvinavalos/crypto-orderbook/orderbook"
)

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

type Update struct {
	Side        string `json:"side"`
	EventTime   string `json:"event_time"`
	PriceLevel  string `json:"price_level"`
	NewQuantity string `json:"new_quantity"`
}

func (u *Update) price() int64 {
	return stringToInt64(u.PriceLevel)
}

func (u *Update) volume() int64 {
	return stringToInt64(u.NewQuantity)
}

type Event struct {
	Type      string   `json:"type"`
	ProductID string   `json:"product_id"`
	Updates   []Update `json:"updates"`
}

type CoinbaseMessage struct {
	Channel        string  `json:"channel"`
	Timestamp      string  `json:"timestamp"`
	SequenceNumber int     `json:"sequence_num"`
	Events         []Event `json:"events"`
}

func (cm *CoinbaseMessage) BookUpdates() []orderbook.BookUpdate {
	updates := make([]orderbook.BookUpdate, 0, len(cm.Events)*10)

	for _, event := range cm.Events {
		for _, update := range event.Updates {
			updates = append(updates, toBookUpdate(update))
		}
	}

	return updates
}

func toBookUpdate(u Update) orderbook.BookUpdate {
	var side orderbook.Side

	if u.Side == "bid" {
		side = orderbook.Bid
	} else {
		side = orderbook.Offer
	}

	price := orderbook.Price(u.price())
	volume := orderbook.Volume(u.volume())

	return orderbook.BookUpdate{
		Side:   side,
		Price:  price,
		Volume: volume,
	}
}
