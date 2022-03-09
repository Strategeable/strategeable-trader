package types

type Ticker struct {
	BidPrice float64
	AskPrice float64
	BidQty   float64
	AskQty   float64
}

func NewTicker(bidPrice float64, askPrice float64, bidQty float64, askQty float64) Ticker {
	return Ticker{
		BidPrice: bidPrice,
		AskPrice: askPrice,
		BidQty:   bidQty,
		AskQty:   askQty,
	}
}
