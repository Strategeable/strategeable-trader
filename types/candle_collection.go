package types

import (
	"sync"
)

type CandleCollection struct {
	windowSize int
	caches     map[string]map[TimeFrame]*CandleCache
	lock       sync.RWMutex
}

func NewCandleCollection() *CandleCollection {
	return &CandleCollection{
		windowSize: 1000,
		caches:     make(map[string]map[TimeFrame]*CandleCache),
	}
}

func (c *CandleCollection) AddTrade(symbol Symbol, trade Trade) {
	caches := c.GetCaches(symbol)

	for _, cache := range caches {
		cache.AddTrade(trade.Price, trade.Quantity, trade.Time)
	}
}

func (c *CandleCollection) InitializeTimeFrame(symbol Symbol, timeFrame TimeFrame, candles []*Candle) {
	caches := c.GetCaches(symbol)

	caches[timeFrame] = NewCandleCache(candles, timeFrame, c.windowSize)
}

func (c *CandleCollection) GetCaches(symbol Symbol) map[TimeFrame]*CandleCache {
	c.lock.Lock()
	defer c.lock.Unlock()

	caches := c.caches[symbol.String()]
	if caches == nil {
		caches = make(map[TimeFrame]*CandleCache)
		c.caches[symbol.String()] = caches
	}

	return caches
}

func (c *CandleCollection) GetCache(symbol Symbol, timeFrame TimeFrame) *CandleCache {
	c.lock.RLock()
	defer c.lock.RUnlock()

	caches := c.caches[symbol.String()]
	if caches == nil {
		return nil
	}

	return caches[timeFrame]
}

func (c *CandleCollection) RegisterSymbol(symbol Symbol) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.caches[symbol.String()] != nil {
		return true
	}

	c.caches[symbol.String()] = make(map[TimeFrame]*CandleCache)
	return false
}
