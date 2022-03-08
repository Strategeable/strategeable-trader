package types

import (
	"sync"
	"time"
)

type PositionHandlerEventType int

type PositionHandlerEvent struct {
	Type PositionHandlerEventType
	Data interface{}
}

const (
	TOTAL_BALANCE_CHANGED PositionHandlerEventType = iota
	POSITION_CREATED
	POSITION_STATE_CHANGED
	POSITION_CLOSED
)

type PositionHandler interface {
	GetTotalBalance() float64
	GetAvailableBalance() float64
	OpenPosition(symbol Symbol, rate float64, quoteSize float64, time time.Time) (*Position, error)
	ClosePosition(symbol Symbol, rate float64, time time.Time) error
	GetPosition(symbol Symbol) *Position
	SubscribeEvents(chan PositionHandlerEvent)
}

type BasePositionHandler struct {
	TotalBalance       float64
	Positions          map[string]*Position
	PositionsLock      sync.RWMutex
	eventSubscriptions []chan PositionHandlerEvent
}

func (b *BasePositionHandler) EmitEvent(event PositionHandlerEvent) {
	for _, ch := range b.eventSubscriptions {
		select {
		case ch <- event:
		default:
		}
	}
}

func (b *BasePositionHandler) SubscribeEvents(channel chan PositionHandlerEvent) {
	b.eventSubscriptions = append(b.eventSubscriptions, channel)
}

func (b *BasePositionHandler) GetPosition(symbol Symbol) *Position {
	b.PositionsLock.RLock()
	defer b.PositionsLock.RUnlock()

	return b.Positions[symbol.String()]
}

func (b *BasePositionHandler) GetTotalBalance() float64 {
	return b.TotalBalance
}

func (b *BasePositionHandler) GetAvailableBalance() float64 {
	balance := b.TotalBalance

	for _, position := range b.Positions {
		if position.IsClosed() {
			continue
		}
		balance -= position.QuoteCost()
	}

	return balance
}
