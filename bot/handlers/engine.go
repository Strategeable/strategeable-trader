package handlers

import (
	"fmt"

	"github.com/Strategeable/Trader/strategy"
	"github.com/Strategeable/Trader/types"
)

type Engine struct {
	MarketDataProvider types.MarketDataProvider
	PositionHandler    types.PositionHandler
	SignalHandler      *SignalHandler
	stopCh             chan struct{}
}

func NewEngine(strategy *strategy.Strategy, marketDataProvider types.MarketDataProvider, positionHandler types.PositionHandler) *Engine {
	return &Engine{
		MarketDataProvider: marketDataProvider,
		SignalHandler:      newSignalHandler(marketDataProvider, positionHandler, strategy),
		stopCh:             make(chan struct{}),
	}
}

func (e *Engine) Start() error {
	fmt.Println("Initializing market data provider.")
	err := e.MarketDataProvider.Init()
	if err != nil {
		return err
	}
	fmt.Println("Initialized market data provider.")

	defer e.MarketDataProvider.Close()

	for {
		select {
		case <-e.stopCh:
			return nil
		case trade := <-e.MarketDataProvider.GetTradeCh():
			e.SignalHandler.handleTrigger(trade)

			if e.MarketDataProvider.RequiresAcks() {
				e.MarketDataProvider.GetAckCh() <- trade.TradeId
			}
		case <-e.MarketDataProvider.GetCloseCh():
			return nil
		}
	}
}

func (e *Engine) Stop() {
	e.stopCh <- struct{}{}
}
