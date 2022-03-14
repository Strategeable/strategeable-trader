package types

type MarketDataProvider interface {
	Init() error
	Close()
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
	b.candleCollection = NewCandleCollection(1000)
}

func (b *BaseMarketDataProvider) GetCandleCollection() *CandleCollection {
	return b.candleCollection
}
