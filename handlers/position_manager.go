package handlers

import (
	"cex-bot/types"
	"errors"
	"fmt"
	"sync"
	"time"
)

type PositionManager struct {
	Positions           map[string]*types.Position
	positionsLock       sync.RWMutex
	ClosedPositions     map[string]*types.Position
	closedPositionsLock sync.RWMutex
}

func newPositionManager(positions []*types.Position, closedPositions []*types.Position) *PositionManager {
	positionMapping := make(map[string]*types.Position)
	closedPositionMapping := make(map[string]*types.Position)

	for _, position := range positions {
		positionMapping[position.Symbol.String()] = position
	}
	for _, position := range closedPositions {
		closedPositionMapping[position.Symbol.String()] = position
	}

	return &PositionManager{
		Positions:       positionMapping,
		ClosedPositions: closedPositionMapping,
	}
}

func (p *PositionManager) closePosition(symbol types.Symbol) {
	p.positionsLock.Lock()
	defer p.positionsLock.Unlock()

	openPosition := p.Positions[symbol.String()]
	if openPosition == nil {
		return
	}

	openPosition.CloseTime = time.Now()

	p.closedPositionsLock.Lock()
	p.ClosedPositions[symbol.String()] = openPosition
	p.closedPositionsLock.Unlock()

	delete(p.Positions, symbol.String())

	fmt.Printf("Closed position for %s.\n", symbol.String())
}

func (p *PositionManager) openNewPosition(symbol types.Symbol) error {
	p.positionsLock.Lock()
	defer p.positionsLock.Unlock()

	if p.Positions[symbol.String()] != nil {
		return errors.New("duplicate position")
	}

	position := &types.Position{
		OpenTime: time.Now(),
		Symbol:   symbol,
	}

	p.Positions[symbol.String()] = position

	fmt.Printf("Created new position for %s.\n", symbol.String())
	return nil
}

func (p *PositionManager) getPosition(symbol types.Symbol) *types.Position {
	p.positionsLock.RLock()
	defer p.positionsLock.RUnlock()

	return p.Positions[symbol.String()]
}

func (p *PositionManager) getClosedPosition(symbol types.Symbol) *types.Position {
	p.closedPositionsLock.RLock()
	defer p.closedPositionsLock.RUnlock()

	return p.ClosedPositions[symbol.String()]
}
