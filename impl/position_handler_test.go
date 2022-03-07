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

func NewTestPositionHandler(accountSize float64, positions []*types.Position) *testPositionHandler {
	positionMapping := make(map[string]*types.Position)

	for _, position := range positions {
		positionMapping[position.Symbol().String()] = position
	}

	positionHandler := &testPositionHandler{}
	positionHandler.Positions = positionMapping
	positionHandler.TotalBalance = accountSize

	return positionHandler
}

func (t *testPositionHandler) ClosePosition(symbol types.Symbol, rate float64, time time.Time) {
	t.PositionsLock.Lock()
	defer t.PositionsLock.Unlock()

	openPosition := t.Positions[symbol.String()]
	if openPosition == nil {
		return
	}

	openPosition.AddOrder(&types.Order{
		Time:   time,
		Side:   types.BUY,
		Active: false,
		Size:   openPosition.BaseSize(),
		Rate:   rate,
		Fills: []types.OrderFill{
			{
				Time:     time,
				Rate:     rate,
				Quantity: openPosition.BaseSize(),
			},
		},
	})

	openPosition.MarkClosed(time)

	fmt.Printf("Closed position for %s.\n", symbol.String())
}

func (t *testPositionHandler) OpenNewPosition(symbol types.Symbol, rate float64, quoteSize float64, time time.Time) (*types.Position, error) {
	t.PositionsLock.Lock()
	defer t.PositionsLock.Unlock()

	if t.Positions[symbol.String()] != nil {
		return nil, errors.New("duplicate position")
	}

	if t.GetAvailableBalance() < quoteSize {
		return nil, errors.New("insufficient balance available")
	}

	position := types.NewPosition(symbol, types.OPEN, time, nil, make([]*types.Order, 0))

	baseSize := quoteSize / rate

	position.AddOrder(&types.Order{
		Time:   time,
		Side:   types.BUY,
		Active: false,
		Size:   baseSize,
		Rate:   rate,
		Fills: []types.OrderFill{
			{
				Time:     time,
				Rate:     rate,
				Quantity: baseSize,
			},
		},
	})

	t.Positions[symbol.String()] = position
	return position, nil
}
