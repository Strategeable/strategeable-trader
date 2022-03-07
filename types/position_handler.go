package types

import (
	"sync"
	"time"
)

type PositionHandler interface {
	GetTotalBalance() float64
	GetAvailableBalance() float64
	OpenNewPosition(symbol Symbol, rate float64, quoteSize float64, time time.Time) (*Position, error)
	ClosePosition(symbol Symbol, rate float64, time time.Time) error
	GetPosition(symbol Symbol) *Position
}

type BasePositionHandler struct {
	TotalBalance  float64
	Positions     map[string]*Position
	PositionsLock sync.RWMutex
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
		balance -= position.QuoteCost()
	}

	return balance
}
