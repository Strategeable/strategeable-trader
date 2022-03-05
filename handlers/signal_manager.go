package handlers

import (
	"cex-bot/strategy"
	"cex-bot/types"
	"fmt"
	"time"
)

type SignalManager struct {
	candleCollection *types.CandleCollection
	positionManager  *PositionManager
	strategy         strategy.Strategy
}

func newSignalManager(candleCollection *types.CandleCollection, positionManager *PositionManager, strategy strategy.Strategy) *SignalManager {
	return &SignalManager{
		positionManager:  positionManager,
		strategy:         strategy,
		candleCollection: candleCollection,
	}
}

func (s *SignalManager) handleTrigger(trade types.Trade) {
	openPosition := s.positionManager.getPosition(trade.Symbol)

	if openPosition != nil {
		sellSignal, err := s.strategy.HasSellSignal(s.candleCollection, trade.Symbol, openPosition)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !sellSignal {
			return
		}

		// Sell position
		s.positionManager.closePosition(trade.Symbol)
		return
	}

	closedPosition := s.positionManager.getClosedPosition(trade.Symbol)
	if closedPosition != nil && time.Since(closedPosition.CloseTime) < s.strategy.BuyCooldown {
		// Symbol is currently on a buy cooldown
		return
	}

	buySignal, err := s.strategy.HasBuySignal(s.candleCollection, trade.Symbol)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !buySignal {
		return
	}

	// Open new position
	s.positionManager.openNewPosition(trade.Symbol)
}
