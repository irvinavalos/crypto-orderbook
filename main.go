package main

import (
	"github.com/irvinavalos/crypto-orderbook/coinbase_api"
)

// func containsKey[M ~map[K]V, K comparable, V any](m M, k K) bool {
// 	_, ok := m[k]
// 	return ok
// }

func main() {
	coinbase_api.StartCoinbaseWS()
}
