package handlers

import (
	"cex-bot/types"
	"sync"
)

type CandleCollection struct {
	windowSize int
	caches     map[string]map[types.TimeFrame]*types.CandleCache
	lock       sync.RWMutex
}

func NewCandleCollection() *CandleCollection {
	return &CandleCollection{
		windowSize: 1000,
		caches:     make(map[string]map[types.TimeFrame]*types.CandleCache),
	}
}

func (c *CandleCollection) AddTrade(symbol types.Symbol, trade types.Trade) {
	caches := c.GetCaches(symbol)

	for _, cache := range caches {
		cache.AddTrade(trade.Price, trade.Quantity, trade.Time)
	}
}

func (c *CandleCollection) InitializeTimeFrame(symbol types.Symbol, timeFrame types.TimeFrame, candles []*types.Candle) {
	caches := c.GetCaches(symbol)

	caches[timeFrame] = types.NewCandleCache(candles, timeFrame, c.windowSize)
}

func (c *CandleCollection) GetCaches(symbol types.Symbol) map[types.TimeFrame]*types.CandleCache {
	c.lock.Lock()
	defer c.lock.Unlock()

	caches := c.caches[symbol.String()]
	if caches == nil {
		caches = make(map[types.TimeFrame]*types.CandleCache)
		c.caches[symbol.String()] = caches
	}

	return caches
}

func (c *CandleCollection) GetCache(symbol types.Symbol, timeFrame types.TimeFrame) *types.CandleCache {
	c.lock.RLock()
	defer c.lock.RUnlock()

	caches := c.caches[symbol.String()]
	if caches == nil {
		return nil
	}

	return caches[timeFrame]
}

func (c *CandleCollection) RegisterSymbol(symbol types.Symbol) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.caches[symbol.String()] != nil {
		return true
	}

	c.caches[symbol.String()] = make(map[types.TimeFrame]*types.CandleCache)
	return false
}
