package handlers

import (
	"github.com/Stratomicl/Trader/strategy"
	"github.com/Stratomicl/Trader/types"
)

// SignalHandler handles incoming trades and decides
// what actions need to be taken on existing positions,
// or what new positions need to be opened. All based on
// the supplied strategy.
type SignalHandler struct {
	marketDataProvider types.MarketDataProvider
	positionHandler    types.PositionHandler
	strategy           *strategy.Strategy
}

// Create a new SignalHandler based on a MarketDataProvider,
// PositionHandler and a Strategy.
func newSignalHandler(marketDataProvider types.MarketDataProvider, positionHandler types.PositionHandler, strategy *strategy.Strategy) *SignalHandler {
	return &SignalHandler{
		positionHandler:    positionHandler,
		strategy:           strategy,
		marketDataProvider: marketDataProvider,
	}
}

// Trades are able to trigger the SignalHandler,
// allowing it to then open or modify (existing) positions.
func (s *SignalHandler) handleTrigger(trade types.Trade) {
	candleCollection := s.marketDataProvider.GetCandleCollection()

	// Fetch current open position on this symbol
	position := s.positionHandler.GetPosition(trade.Symbol)

	if position != nil {
		// Position should not be in a transitioning state like CLOSING or OPENING
		if position.State() != types.OPEN {
			return
		}

		// Run sell checks
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

	// Fetch latest closed position on this symbol
	closedPosition := s.positionHandler.GetClosedPosition(trade.Symbol)

	if closedPosition != nil && trade.Time.Sub(*closedPosition.CloseTime()) < s.strategy.BuyCooldown {
		// Symbol is currently on a buy cooldown
		return
	}

	// Run buy checks
	buySignal, err := s.strategy.HasBuySignal(candleCollection, trade.Symbol)
	if err != nil {
		return
	}

	if !buySignal {
		return
	}

	availableBalance := s.positionHandler.GetAvailableBalance()

	// Determine the quote size we want to open this position with
	quoteSize := s.positionHandler.GetTotalBalance() / 100 * s.strategy.BuySize

	// Check if enough balance is available
	if quoteSize > availableBalance {
		// Check if the available balance is enough to cover the minimum
		// trade size on this symbol
		if availableBalance < trade.Symbol.MinQuoteSize() {
			return
		}

		quoteSize = availableBalance
	}

	// Open new position
	s.positionHandler.OpenPosition(trade.Symbol, trade.Price, quoteSize, trade.Time)
}
