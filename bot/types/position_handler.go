package types

import (
	"sync"
	"time"
)

type PositionHandlerEventType string

type PositionHandlerEvent struct {
	Time time.Time
	Type PositionHandlerEventType
	Data interface{}
}

const (
	TOTAL_BALANCE_CHANGED  PositionHandlerEventType = "TOTAL_BALANCE_CHANGED"
	POSITION_CREATED       PositionHandlerEventType = "POSITION_CREATED"
	POSITION_STATE_CHANGED PositionHandlerEventType = "POSITION_STATE_CHANGED"
	POSITION_CLOSED        PositionHandlerEventType = "POSITION_CLOSED"
)

// Handles opening and closing positions.
// Also maintains the current total balance and
// available balance.
type PositionHandler interface {
	// Get total balance without change on current open positions
	GetTotalBalance() float64
	// Get balance without all initial quote costs of all open positions
	GetAvailableBalance() float64
	// Open a new position on a symbol with a specific quote size
	OpenPosition(symbol Symbol, rate float64, quoteSize float64, time time.Time) (*Position, error)
	// Close a position on a symbol
	ClosePosition(symbol Symbol, rate float64, time time.Time) error
	// Get an open position
	GetPosition(symbol Symbol) *Position
	// Get a closed position
	GetClosedPosition(symbol Symbol) *Position
	// Subscribe to all position handler events
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
		ch <- event
	}
}

func (b *BasePositionHandler) SubscribeEvents(channel chan PositionHandlerEvent) {
	b.eventSubscriptions = append(b.eventSubscriptions, channel)
}

func (b *BasePositionHandler) GetPosition(symbol Symbol) *Position {
	b.PositionsLock.RLock()
	defer b.PositionsLock.RUnlock()

	position := b.Positions[symbol.String()]

	if position != nil && position.IsClosed() {
		return nil
	}

	return position
}

func (b *BasePositionHandler) GetClosedPosition(symbol Symbol) *Position {
	b.PositionsLock.RLock()
	defer b.PositionsLock.RUnlock()

	position := b.Positions[symbol.String()]

	if position != nil && !position.IsClosed() {
		return nil
	}

	return position
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
