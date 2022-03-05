package types

import "sync"

type PositionHandler interface {
	OpenNewPosition(symbol Symbol) error
	ClosePosition(symbol Symbol) error
	GetPosition(symbol Symbol) *Position
	GetClosedPosition(symbol Symbol) *Position
}

type BasePositionHandler struct {
	Positions           map[string]*Position
	PositionsLock       sync.RWMutex
	ClosedPositions     map[string]*Position
	ClosedPositionsLock sync.RWMutex
}

func (b *BasePositionHandler) GetPosition(symbol Symbol) *Position {
	b.PositionsLock.RLock()
	defer b.PositionsLock.RUnlock()

	return b.Positions[symbol.String()]
}

func (b *BasePositionHandler) GetClosedPosition(symbol Symbol) *Position {
	b.ClosedPositionsLock.RLock()
	defer b.ClosedPositionsLock.RUnlock()

	return b.ClosedPositions[symbol.String()]
}
