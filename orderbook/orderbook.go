package orderbook

type Price int64
type Volume int64

type Side string

const (
	Bid   Side = "bid"
	Offer Side = "offer"
)

type PriceList map[Price]Volume

type Orderbook struct {
	Bids      PriceList
	Offers    PriceList
	bestBid   Price
	bestOffer Price
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Bids:      make(PriceList),
		Offers:    make(PriceList),
		bestBid:   0,
		bestOffer: 0,
	}
}

func (ob *Orderbook) BestBid() Price {
	return ob.bestBid
}

func (ob *Orderbook) BestOffer() Price {
	return ob.bestOffer
}

func (ob *Orderbook) ApplyUpdate(update BookUpdate) {
	if update.Side == Bid {
		updateBids(ob, update)
	} else {
		updateOffers(ob, update)
	}
}

func updateBids(ob *Orderbook, u BookUpdate) bool {
	oldBestBid := ob.bestBid
	price := u.Price
	volume := u.Volume

	if volume == 0 {
		delete(ob.Bids, price)
		if price == oldBestBid {
			ob.bestBid = findBestBid(ob.Bids)
		}
	} else {
		ob.Bids[price] = volume
		if price > ob.bestBid {
			ob.bestBid = price
		}
	}

	return ob.bestBid != oldBestBid
}

func findBestBid(bids PriceList) Price {
	var best Price

	for price := range bids {
		if price > best {
			best = price
		}
	}

	return best
}

func updateOffers(ob *Orderbook, u BookUpdate) bool {
	oldBestOffer := ob.bestOffer
	price := u.Price
	volume := u.Volume

	if volume == 0 {
		delete(ob.Offers, price)
		if price == oldBestOffer {
			ob.bestOffer = findBestOffer(ob.Offers)
		}
	} else {
		ob.Offers[price] = volume
		if ob.bestOffer == 0 || price < ob.bestOffer {
			ob.bestOffer = price
		}
	}

	return ob.bestOffer != oldBestOffer
}

func findBestOffer(bids PriceList) Price {
	var best Price

	for price := range bids {
		if best == 0 || price < best {
			best = price
		}
	}

	return best
}

type BookUpdate struct {
	Side   Side
	Price  Price
	Volume Volume
}
