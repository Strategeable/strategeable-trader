package types

import "time"

type PositionState int

const (
	OPENING PositionState = iota
	OPEN
	CLOSING
	CLOSED
)

type OrderSide string

const (
	BUY  OrderSide = "BUY"
	SELL OrderSide = "SELL"
)

type OrderFill struct {
	Time     time.Time
	Rate     float64
	Quantity float64
}

type Order struct {
	OrderId string
	Time    time.Time
	Side    OrderSide
	Active  bool
	Fills   []OrderFill
}

func (o *Order) FilledSize() float64 {
	total := float64(0)

	for _, fill := range o.Fills {
		total += fill.Quantity
	}

	return total
}

func (o *Order) QuoteValue() float64 {
	total := float64(0)

	for _, fill := range o.Fills {
		total += fill.Quantity * fill.Rate
	}

	return total
}

type Position struct {
	Symbol    Symbol
	State     PositionState
	OpenTime  time.Time
	CloseTime time.Time
	Orders    []*Order
}

func (p *Position) BaseSize() float64 {
	total := float64(0)

	for _, order := range p.Orders {
		if order.Side == BUY {
			total += order.FilledSize()
		}
	}

	return total
}

func (p *Position) QuoteCost() float64 {
	total := float64(0)

	for _, order := range p.Orders {
		if order.Side == BUY {
			total += order.QuoteValue()
		}
	}

	return total
}

func (p *Position) QuoteValue(rate float64) float64 {
	total := p.BaseSize()

	product := float64(0)

	for _, order := range p.Orders {
		if order.Side == SELL {
			product += order.QuoteValue()
			total -= order.FilledSize()
		}
	}

	product += total * rate

	return product
}

func (p *Position) AverageQuoteRate(rate float64) float64 {
	total := p.BaseSize()

	product := float64(0)

	for _, order := range p.Orders {
		if order.Side == SELL {
			product += order.QuoteValue()
			total -= order.FilledSize()
		}
	}

	product += total * rate

	return product / p.BaseSize()
}
