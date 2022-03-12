package types

import (
	"sync"
	"time"
)

type CandleCache struct {
	candles           []*Candle
	candleTimeMapping map[time.Time]*Candle
	mappingLock       sync.RWMutex
	timeFrame         TimeFrame
	windowSize        int
}

func NewCandleCache(candles []*Candle, timeFrame TimeFrame, windowSize int) *CandleCache {
	mapping := make(map[time.Time]*Candle)
	initialCandles := candles
	if len(initialCandles) > windowSize {
		initialCandles = initialCandles[len(initialCandles)-windowSize:]
	}

	for _, candle := range initialCandles {
		mapping[candle.OpenTime] = candle
	}

	return &CandleCache{
		candles:           initialCandles,
		candleTimeMapping: mapping,
		timeFrame:         timeFrame,
		windowSize:        windowSize,
	}
}

func (c *CandleCache) GetAvailableRanges(exchange Exchange, symbol Symbol, timeFrame TimeFrame) []DateRange {
	ranges := make([]DateRange, 0)

	return ranges
}

func (c *CandleCache) GetSize() int {
	return len(c.candles)
}

func (c *CandleCache) GetCandles() []*Candle {
	return c.candles
}

func (c *CandleCache) GetCurrentCandle() *Candle {
	return c.candles[len(c.candles)-1]
}

func (c *CandleCache) GetCurrentRate() float64 {
	return c.candles[len(c.candles)-1].Close
}

func (c *CandleCache) GetCandleAt(time time.Time) *Candle {
	c.mappingLock.RLock()
	defer c.mappingLock.RUnlock()

	return c.candleTimeMapping[time]
}

func (c *CandleCache) AddTrade(rate float64, volume float64, currentTime time.Time) bool {
	timeNano := currentTime.UnixNano()
	candleStartTime := timeNano - timeNano%CandleDurations[c.timeFrame].Nanoseconds()
	openTime := time.Unix(0, candleStartTime)

	var isNew bool
	var candle *Candle

	if len(c.candles) > 0 {
		currentCandle := c.candles[len(c.candles)-1]

		if currentCandle.OpenTime.UnixNano() == openTime.UnixNano() {
			candle = currentCandle
		}
	}

	if candle == nil {
		closeTime := time.Unix(0, candleStartTime+CandleDurations[c.timeFrame].Nanoseconds()-1)

		candle = NewCandle(
			openTime,
			closeTime,
			rate,
			rate,
			rate,
			rate,
			0,
		)

		c.candles = append(c.candles, candle)
		isNew = true

		c.mappingLock.Lock()
		c.candleTimeMapping[candle.OpenTime] = candle
		c.mappingLock.Unlock()

		if len(c.candles) > c.windowSize {
			c.candles = c.candles[len(c.candles)-c.windowSize:]
		}
	}

	candle.addTrade(rate, volume, currentTime)
	return isNew
}
