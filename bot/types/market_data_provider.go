package types

type MarketDataProvider interface {
	Init() error
	GetCandleCollection() *CandleCollection
	GetTradeCh() chan Trade
	RequiresAcks() bool
	GetAckCh() chan string
	GetCloseCh() chan string
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
