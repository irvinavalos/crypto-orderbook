package orderbook

import "sync"

type Price int64
type Quantity int64
type Side string

const (
	Bid Side = "bid"
	Ask Side = "ask"
)

type PriceLevels map[Price]Quantity

type OrderbookUpdate struct {
	Side     Side
	Price    Price
	Quantity Quantity
}

type OrderbookUpdates []OrderbookUpdate

type Orderbook struct {
	Bids    PriceLevels
	Asks    PriceLevels
	bestBid Price
	bestAsk Price
	mu      sync.RWMutex
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Bids:    PriceLevels{},
		Asks:    PriceLevels{},
		bestBid: 0,
		bestAsk: 0,
	}
}

func (ob *Orderbook) BestBid() (Price, Quantity, bool) {
	if len(ob.Bids) == 0 {
		return 0, 0, false
	}
	return ob.bestBid, ob.Bids[ob.bestBid], true
}

func (ob *Orderbook) BestAsk() (Price, Quantity, bool) {
	if len(ob.Asks) == 0 {
		return 0, 0, false
	}
	return ob.bestAsk, ob.Asks[ob.bestAsk], true
}

func (ob *Orderbook) ApplyUpdate(ou OrderbookUpdate) {
	if ou.Price <= 0 {
		return
	}
	if ou.Quantity < 0 {
		return
	}
	ob.mu.Lock()
	defer ob.mu.Unlock()

	var best Price

	switch ou.Side {
	case Bid:
		best = ob.bestBid
	case Ask:
		best = ob.bestAsk
	}

	switch ou.Quantity {
	case 0:
		ob.removeOrder(ou.Price, best, ou.Side)
	default:
		ob.addOrder(ou.Price, ou.Quantity, ou.Side)
	}
}

func (ob *Orderbook) addOrder(price Price, quantity Quantity, side Side) {
	switch side {
	case Bid:
		ob.Bids[price] = quantity
		if price > ob.bestBid {
			ob.bestBid = price
		}
	case Ask:
		ob.Asks[price] = quantity
		if ob.bestAsk == 0 || price < ob.bestAsk {
			ob.bestAsk = price
		}
	}
}

func (ob *Orderbook) removeOrder(price Price, best Price, side Side) {
	switch side {
	case Bid:
		delete(ob.Bids, price)
	case Ask:
		delete(ob.Asks, price)
	}

	if price == best {
		ob.updateBestBidOrOffer(side)
	}
}

func (ob *Orderbook) updateBestBidOrOffer(side Side) {
	switch side {
	case Bid:
		ob.bestBid = ob.findBestBid()
	case Ask:
		ob.bestAsk = ob.findBestOffer()
	}
}

func (ob *Orderbook) findBestBid() Price {
	var best Price
	for price := range ob.Bids {
		if price > best {
			best = price
		}
	}
	return best
}

func (ob *Orderbook) findBestOffer() Price {
	var best Price
	for price := range ob.Asks {
		if best == 0 || price < best {
			best = price
		}
	}
	return best
}

func (ob *Orderbook) Spread() Price {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.bestAsk == 0 || ob.bestBid == 0 {
		return 0
	}
	return ob.bestAsk - ob.bestBid
}

func (ob *Orderbook) MidPoint() Price {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.bestAsk == 0 || ob.bestBid == 0 {
		return 0
	}
	return (ob.bestAsk + ob.bestBid) / 2
}
