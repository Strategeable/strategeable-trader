package handlers

import (
	"github.com/Strategeable/Trader/strategy"
	"github.com/Strategeable/Trader/types"
)

// (trading) Engine serves as the main trading logic abstraction.
// The logic within Engine can be be used for any situation,
// like simulated mode, live trading mode and backtesting.
type Engine struct {
	marketDataProvider types.MarketDataProvider
	signalHandler      *SignalHandler

	// Used to shut down the engine once it has been started
	stopCh chan struct{}
}

// Create a new (trading) Engine based on all of its different components.
func NewEngine(strategy *strategy.Strategy, marketDataProvider types.MarketDataProvider, positionHandler types.PositionHandler) *Engine {
	return &Engine{
		marketDataProvider: marketDataProvider,
		signalHandler:      newSignalHandler(marketDataProvider, positionHandler, strategy),
		stopCh:             make(chan struct{}),
	}
}

// Initializes the market data and starts responding to
// new incoming trades.
func (e *Engine) Start() error {
	// Prepare required market data
	err := e.marketDataProvider.Init()
	if err != nil {
		return err
	}

	defer e.marketDataProvider.Close()

	for {
		select {
		case <-e.stopCh:
			// Break out of the engine loop
			return nil
		case trade := <-e.marketDataProvider.GetTradeCh():
			// Notify signal handler about the new trade
			e.signalHandler.handleTrigger(trade)

			// When for example running backtests, we want to check all
			// individual trades one by one and not simultaneously.
			// Acking here tells the MarketDataProvider that we handled this trade.
			if e.marketDataProvider.RequiresAcks() {
				e.marketDataProvider.GetAckCh() <- trade.TradeId
			}
		case <-e.marketDataProvider.GetCloseCh():
			return nil
		}
	}
}

// Shut down the engine, stop watching for trades.
func (e *Engine) Stop() {
	e.stopCh <- struct{}{}
}
