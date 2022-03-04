package types

import "time"

type ExchangeImplementation interface {
	// Initialization
	Init() error

	// Basic data
	GetSymbols() ([]Symbol, error)
	GetUniqueQuoteAssets() ([]string, error)
	FormatSymbol(symbol Symbol) string

	GetAvailableTimeFrames() []TimeFrame
	FormatTimeFrame(timeFrame TimeFrame) string

	// Candle data
	GetCandles(symbol Symbol, timeFrame TimeFrame, limit int) ([]*Candle, error)
	GetHistoricalCandles(symbol Symbol, timeFrame TimeFrame, from time.Time, to time.Time, limit int) ([]*Candle, error)

	// Ticker data
	GetTicker(symbol Symbol) (Ticker, error)
}
