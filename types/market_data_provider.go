package types

type MarketDataProvider interface {
	Init() error
	GetCandleCollection() *CandleCollection
	GetTradeCh() chan Trade
}

type BaseMarketDataProvider struct {
	candleCollection *CandleCollection
}

func (b *BaseMarketDataProvider) InitCandleCollection() {
	b.candleCollection = NewCandleCollection()
}

func (b *BaseMarketDataProvider) GetCandleCollection() *CandleCollection {
	return b.candleCollection
}
