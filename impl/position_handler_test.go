package impl

import (
	"cex-bot/types"
	"errors"
	"fmt"
	"time"
)

type testPositionHandler struct {
	types.BasePositionHandler
}

func NewTestPositionHandler(positions []*types.Position, closedPositions []*types.Position) *testPositionHandler {
	positionMapping := make(map[string]*types.Position)
	closedPositionMapping := make(map[string]*types.Position)

	for _, position := range positions {
		positionMapping[position.Symbol.String()] = position
	}
	for _, position := range closedPositions {
		closedPositionMapping[position.Symbol.String()] = position
	}

	positionHandler := &testPositionHandler{}
	positionHandler.ClosedPositions = closedPositionMapping
	positionHandler.Positions = positionMapping

	return positionHandler
}

func (t *testPositionHandler) ClosePosition(symbol types.Symbol) {
	t.PositionsLock.Lock()
	defer t.PositionsLock.Unlock()

	openPosition := t.Positions[symbol.String()]
	if openPosition == nil {
		return
	}

	openPosition.CloseTime = time.Now()

	t.ClosedPositionsLock.Lock()
	t.ClosedPositions[symbol.String()] = openPosition
	t.ClosedPositionsLock.Unlock()

	delete(t.Positions, symbol.String())

	fmt.Printf("Closed position for %s.\n", symbol.String())
}

func (t *testPositionHandler) OpenNewPosition(symbol types.Symbol) error {
	t.PositionsLock.Lock()
	defer t.PositionsLock.Unlock()

	if t.Positions[symbol.String()] != nil {
		return errors.New("duplicate position")
	}

	position := &types.Position{
		OpenTime: time.Now(),
		Symbol:   symbol,
	}

	t.Positions[symbol.String()] = position

	fmt.Printf("Created new position for %s.\n", symbol.String())
	return nil
}
