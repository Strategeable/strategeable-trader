package types

import (
	"time"
)

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
	Size    float64
	Rate    float64
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
	symbol    Symbol
	state     PositionState
	openTime  time.Time
	closeTime *time.Time
	orders    []*Order
}

func NewPosition(symbol Symbol, state PositionState, openTime time.Time, closeTime *time.Time, orders []*Order) *Position {
	return &Position{
		symbol:    symbol,
		state:     state,
		openTime:  openTime,
		closeTime: closeTime,
		orders:    orders,
	}
}

func (p *Position) ChangePercentage(rate float64) float64 {
	if len(p.orders) == 0 {
		return 0
	}

	entryRate := p.AverageEntryRate()
	exitRate := p.AverageExitRate(rate)

	return (exitRate - entryRate) / entryRate * 100
}

func (p *Position) ChangeAmount(rate float64) float64 {
	if len(p.orders) == 0 {
		return 0
	}

	entryQuoteSize := p.QuoteCost()
	exitQuoteSize := p.QuoteValue(rate)

	return exitQuoteSize - entryQuoteSize
}

func (p *Position) AddOrder(order *Order) {
	p.orders = append(p.orders, order)
}

func (p *Position) State() PositionState {
	return p.state
}

func (p *Position) Symbol() *Symbol {
	return &p.symbol
}

func (p *Position) CloseTime() *time.Time {
	return p.closeTime
}

func (p *Position) IsClosed() bool {
	return p.closeTime != nil
}

func (p *Position) MarkClosed(time time.Time) {
	p.closeTime = &time
	p.SetState(CLOSED)
}

func (p *Position) SetState(state PositionState) {
	p.state = state
}

func (p *Position) BaseSize() float64 {
	total := float64(0)

	for _, order := range p.orders {
		if order.Side == BUY {
			total += order.FilledSize()
		}
	}

	return total
}

func (p *Position) QuoteCost() float64 {
	total := float64(0)

	for _, order := range p.orders {
		if order.Side == BUY {
			total += order.QuoteValue()
		}
	}

	return total
}

func (p *Position) QuoteValue(rate float64) float64 {
	total := p.BaseSize()

	product := float64(0)

	for _, order := range p.orders {
		if order.Side == SELL {
			product += order.QuoteValue()
			total -= order.FilledSize()
		}
	}

	product += total * rate

	return product
}

func (p *Position) AverageEntryRate() float64 {
	return p.averageRate(0, BUY)
}

func (p *Position) AverageExitRate(rate float64) float64 {
	return p.averageRate(rate, SELL)
}

func (p *Position) averageRate(rate float64, side OrderSide) float64 {
	total := p.BaseSize()

	product := float64(0)

	for _, order := range p.orders {
		if order.Side == side {
			product += order.QuoteValue()
			total -= order.FilledSize()
		}
	}

	product += total * rate

	return product / p.BaseSize()
}
