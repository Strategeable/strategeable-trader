package handlers

import (
	"cex-bot/strategy"
	"cex-bot/types"
)

type ExchangeDetails struct {
	Name     string
	Exchange types.Exchange
	Live     bool
}

type Bot struct {
	Name             string
	ExchangeDetails  ExchangeDetails
	PositionManager  *PositionManager
	SignalManager    *SignalManager
	Strategy         strategy.Strategy
	stopCh           chan struct{}
	TradeCh          chan types.Trade
	CandleCollection *types.CandleCollection
}

func NewBot(name string, exchangeDetails ExchangeDetails, strategy strategy.Strategy, positions []*types.Position, closedPositions []*types.Position, candleCollection *types.CandleCollection) *Bot {
	positionManager := newPositionManager(positions, closedPositions)

	return &Bot{
		Name:             name,
		ExchangeDetails:  exchangeDetails,
		Strategy:         strategy,
		PositionManager:  positionManager,
		SignalManager:    newSignalManager(candleCollection, positionManager, strategy),
		CandleCollection: candleCollection,
		stopCh:           make(chan struct{}),
		TradeCh:          make(chan types.Trade, 5),
	}
}

func (b *Bot) Init() error {
	return nil
}

func (b *Bot) RunLoop() {
	for {
		select {
		case <-b.stopCh:
			return
		case trade := <-b.TradeCh:
			b.SignalManager.handleTrigger(trade)
		}
	}
}

func (b *Bot) Stop() {
	b.stopCh <- struct{}{}
}
