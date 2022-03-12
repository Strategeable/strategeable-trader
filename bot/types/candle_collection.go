package types

import (
	"sync"
	"time"
)

type DateRange struct {
	From time.Time
	To   time.Time
}

type CandleCollection struct {
	windowSize int
	caches     map[Exchange]map[string]map[TimeFrame]*CandleCache
	lock       sync.RWMutex
}

func NewCandleCollection(windowSize int) *CandleCollection {
	return &CandleCollection{
		windowSize: windowSize,
		caches:     make(map[Exchange]map[string]map[TimeFrame]*CandleCache),
	}
}

func (c *CandleCollection) AddTrade(exchange Exchange, symbol Symbol, trade Trade) {
	caches := c.GetCaches(exchange, symbol)

	for _, cache := range caches {
		cache.AddTrade(trade.Price, trade.Quantity, trade.Time)
	}
}

func (c *CandleCollection) InitializeTimeFrame(exchange Exchange, symbol Symbol, timeFrame TimeFrame, candles []*Candle) {
	caches := c.GetCaches(exchange, symbol)

	caches[timeFrame] = NewCandleCache(candles, timeFrame, c.windowSize)
}

func (c *CandleCollection) GetCaches(exchange Exchange, symbol Symbol) map[TimeFrame]*CandleCache {
	c.lock.Lock()
	defer c.lock.Unlock()

	exchangeCaches := c.caches[exchange]
	if exchangeCaches == nil {
		exchangeCaches = make(map[string]map[TimeFrame]*CandleCache)
		c.caches[exchange] = exchangeCaches
	}

	caches := exchangeCaches[symbol.String()]
	if caches == nil {
		caches = make(map[TimeFrame]*CandleCache)
		exchangeCaches[symbol.String()] = caches
	}

	return caches
}

func (c *CandleCollection) GetCache(exchange Exchange, symbol Symbol, timeFrame TimeFrame) *CandleCache {
	c.lock.RLock()
	defer c.lock.RUnlock()

	exchangeCaches := c.caches[exchange]
	if exchangeCaches == nil {
		exchangeCaches = make(map[string]map[TimeFrame]*CandleCache)
		c.caches[exchange] = exchangeCaches
	}

	caches := exchangeCaches[symbol.String()]
	if caches == nil {
		return nil
	}

	return caches[timeFrame]
}

func (c *CandleCollection) RegisterSymbol(exchange Exchange, symbol Symbol) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	exchangeCaches := c.caches[exchange]
	if exchangeCaches == nil {
		exchangeCaches = make(map[string]map[TimeFrame]*CandleCache)
		c.caches[exchange] = exchangeCaches
	}

	if exchangeCaches[symbol.String()] != nil {
		return true
	}

	exchangeCaches[symbol.String()] = make(map[TimeFrame]*CandleCache)
	return false
}
