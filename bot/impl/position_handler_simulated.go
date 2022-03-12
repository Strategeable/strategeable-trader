package impl

import (
	"errors"
	"time"

	"github.com/Stratomicl/Trader/types"
)

type simulatedPositionHandler struct {
	types.BasePositionHandler
}

func NewSimulatedPositionHandler(accountSize float64, positions []*types.Position) *simulatedPositionHandler {
	positionMapping := make(map[string]*types.Position)

	for _, position := range positions {
		positionMapping[position.Symbol().String()] = position
	}

	positionHandler := &simulatedPositionHandler{}
	positionHandler.Positions = positionMapping
	positionHandler.TotalBalance = accountSize

	return positionHandler
}

func (s *simulatedPositionHandler) ClosePosition(symbol types.Symbol, rate float64, time time.Time) error {
	s.PositionsLock.Lock()
	defer s.PositionsLock.Unlock()

	openPosition := s.Positions[symbol.String()]
	if openPosition == nil {
		return nil
	}

	openPosition.AddOrder(&types.Order{
		Time:   time,
		Side:   types.SELL,
		Active: false,
		Size:   openPosition.BaseSize(),
		Rate:   rate,
		Fills: []types.OrderFill{
			{
				Time:     time,
				Rate:     rate,
				Quantity: openPosition.BaseSize(),
				QuoteFee: rate * openPosition.BaseSize() * 0.00075,
			},
		},
	})

	s.TotalBalance += openPosition.ChangeAmount(rate)

	s.EmitEvent(types.PositionHandlerEvent{
		Time: time,
		Type: types.TOTAL_BALANCE_CHANGED,
		Data: s.TotalBalance,
	})

	openPosition.MarkClosed(time)

	s.EmitEvent(types.PositionHandlerEvent{
		Time: time,
		Type: types.POSITION_CLOSED,
		Data: openPosition,
	})

	return nil
}

func (s *simulatedPositionHandler) OpenPosition(symbol types.Symbol, rate float64, quoteSize float64, time time.Time) (*types.Position, error) {
	s.PositionsLock.Lock()
	defer s.PositionsLock.Unlock()

	existingPosition := s.Positions[symbol.String()]
	if existingPosition != nil && !existingPosition.IsClosed() {
		return nil, errors.New("duplicate position")
	}

	if s.GetAvailableBalance() < quoteSize {
		return nil, errors.New("insufficient balance available")
	}

	position := types.NewPosition(symbol, types.OPEN, time, nil, make([]*types.Order, 0))

	baseSize := quoteSize / rate
	fees := baseSize * 0.00075

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
				Quantity: baseSize - fees,
				QuoteFee: fees * rate,
			},
		},
	})

	s.Positions[symbol.String()] = position

	s.EmitEvent(types.PositionHandlerEvent{
		Time: time,
		Type: types.POSITION_CREATED,
		Data: position,
	})
	return position, nil
}
