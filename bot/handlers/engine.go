package handlers

import (
	"fmt"

	"github.com/Stratomicl/Trader/strategy"
	"github.com/Stratomicl/Trader/types"
)

type Engine struct {
	MarketDataProvider types.MarketDataProvider
	PositionHandler    types.PositionHandler
	SignalHandler      *SignalHandler
	stopCh             chan struct{}
}

func NewEngine(strategy strategy.Strategy, marketDataProvider types.MarketDataProvider, positionHandler types.PositionHandler) *Engine {
	return &Engine{
		MarketDataProvider: marketDataProvider,
		SignalHandler:      newSignalHandler(marketDataProvider, positionHandler, strategy),
		stopCh:             make(chan struct{}),
	}
}

func (e *Engine) Start() error {
	err := e.MarketDataProvider.Init()
	if err != nil {
		return err
	}

	for {
		select {
		case <-e.stopCh:
			return nil
		case trade := <-e.MarketDataProvider.GetTradeCh():
			e.SignalHandler.handleTrigger(trade)
			fmt.Println(trade.Time)
		}
	}
}

func (e *Engine) Stop() {
	e.stopCh <- struct{}{}
}
