package handlers

import (
	"cex-bot/strategy"
	"cex-bot/types"
	"fmt"
)

type ExchangeDetails struct {
	Name     string
	Exchange types.Exchange
	Live     bool
}

type Bot struct {
	Name             string
	ExchangeDetails  ExchangeDetails
	Strategy         strategy.Strategy
	Positions        []*types.Position
	stopCh           chan struct{}
	TradeCh          chan types.Trade
	CandleCollection *types.CandleCollection
}

func NewBot(name string, exchangeDetails ExchangeDetails, strategy strategy.Strategy, positions []*types.Position, candleCollection *types.CandleCollection) *Bot {
	return &Bot{
		Name:             name,
		ExchangeDetails:  exchangeDetails,
		Strategy:         strategy,
		Positions:        positions,
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
			fmt.Println(b.Name, trade)
		}
	}
}

func (b *Bot) Stop() {
	b.stopCh <- struct{}{}
}
