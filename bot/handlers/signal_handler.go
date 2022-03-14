package handlers

import (
	"github.com/Strategeable/Trader/strategy"
	"github.com/Strategeable/Trader/types"
)

type SignalHandler struct {
	marketDataProvider types.MarketDataProvider
	positionHandler    types.PositionHandler
	strategy           *strategy.Strategy
}

func newSignalHandler(marketDataProvider types.MarketDataProvider, positionHandler types.PositionHandler, strategy *strategy.Strategy) *SignalHandler {
	return &SignalHandler{
		positionHandler:    positionHandler,
		strategy:           strategy,
		marketDataProvider: marketDataProvider,
	}
}

func (s *SignalHandler) handleTrigger(trade types.Trade) {
	candleCollection := s.marketDataProvider.GetCandleCollection()

	position := s.positionHandler.GetPosition(trade.Symbol)

	if position != nil && !position.IsClosed() {
		if position.State() != types.OPEN {
			return
		}

		sellSignal, err := s.strategy.HasSellSignal(candleCollection, trade.Symbol, position)
		if err != nil {
			return
		}

		if !sellSignal {
			return
		}

		// Sell position
		s.positionHandler.ClosePosition(trade.Symbol, trade.Price, trade.Time)
		return
	}

	if position != nil && position.IsClosed() && trade.Time.Sub(*position.CloseTime()) < s.strategy.BuyCooldown {
		// Symbol is currently on a buy cooldown
		return
	}

	buySignal, err := s.strategy.HasBuySignal(candleCollection, trade.Symbol)
	if err != nil {
		return
	}

	if !buySignal {
		return
	}

	availableBalance := s.positionHandler.GetAvailableBalance()

	quoteSize := s.positionHandler.GetTotalBalance() / 100 * s.strategy.BuySize
	if quoteSize > availableBalance {
		if availableBalance < trade.Symbol.MinQuoteSize() {
			return
		}

		quoteSize = availableBalance
	}

	// Open new position
	s.positionHandler.OpenPosition(trade.Symbol, trade.Price, quoteSize, trade.Time)
}
