package types

// Keeps a cache of known candles and
// offers a stream of new incoming trades
// that keep the cache up to date.
type MarketDataProvider interface {
	// Prepare initial candle caches
	Init() error
	IsInitialized() bool
	// Stop watching new trades and clear cache
	Close()
	GetCandleCollection() *CandleCollection
	// Channel offering all new incoming trades
	GetTradeCh() chan Trade
	// Does the provider need to know if a trade has been handled?
	// Used in backtests
	RequiresAcks() bool
	GetAckCh() chan string
	// This channel gets closed once the MarketDataProvider is done/closed
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
