package types

import "time"

type Trade struct {
	Symbol   Symbol
	Time     time.Time
	TradeId  string
	Price    float64
	Quantity float64
}

type SubscriptionManager interface {
	Close()
}

type ExchangeImplementation interface {
	// Initialization
	Init() error

	// Basic data
	GetSymbols() ([]Symbol, error)
	GetExchange() Exchange
	GetUniqueQuoteAssets() ([]string, error)
	FormatSymbol(symbol Symbol) string

	GetAvailableTimeFrames() []TimeFrame
	FormatTimeFrame(timeFrame TimeFrame) string

	// Candle data
	GetCandles(symbol Symbol, timeFrame TimeFrame, limit int) ([]*Candle, error)
	GetHistoricalCandles(symbol Symbol, timeFrame TimeFrame, from time.Time, to time.Time) ([]*Candle, error)

	// Ticker data
	GetTicker(symbol Symbol) (Ticker, error)

	// Real-time
	WatchTrades(symbols []Symbol, handleTrade func(trade Trade), handleClose func(error)) (func(), error)
}
