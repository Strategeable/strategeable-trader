package types

type MarketDataProvider interface {
	Init() error
	GetCandleCollection() *CandleCollection
	GetTradeCh() chan Trade
}
