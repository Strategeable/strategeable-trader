package handlers

import (
	"cex-bot/strategy"
	"cex-bot/types"
	"fmt"
	"time"
)

type SignalHandler struct {
	marketDataProvider types.MarketDataProvider
	positionHandler    types.PositionHandler
	strategy           strategy.Strategy
}

func newSignalHandler(marketDataProvider types.MarketDataProvider, positionHandler types.PositionHandler, strategy strategy.Strategy) *SignalHandler {
	return &SignalHandler{
		positionHandler:    positionHandler,
		strategy:           strategy,
		marketDataProvider: marketDataProvider,
	}
}

func (s *SignalHandler) handleTrigger(trade types.Trade) {
	candleCollection := s.marketDataProvider.GetCandleCollection()

	openPosition := s.positionHandler.GetPosition(trade.Symbol)

	if openPosition != nil {
		if openPosition.State != types.OPEN {
			return
		}

		sellSignal, err := s.strategy.HasSellSignal(candleCollection, trade.Symbol, openPosition)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !sellSignal {
			return
		}

		// Sell position
		s.positionHandler.ClosePosition(trade.Symbol)
		return
	}

	closedPosition := s.positionHandler.GetClosedPosition(trade.Symbol)
	if closedPosition != nil && time.Since(closedPosition.CloseTime) < s.strategy.BuyCooldown {
		// Symbol is currently on a buy cooldown
		return
	}

	buySignal, err := s.strategy.HasBuySignal(candleCollection, trade.Symbol)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !buySignal {
		return
	}

	// Open new position
	s.positionHandler.OpenNewPosition(trade.Symbol)
}
