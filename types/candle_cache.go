package types

type CandleCache struct {
	candles   []*Candle
	timeFrame TimeFrame
}

func NewCandleCache(candles []*Candle, timeFrame TimeFrame) *CandleCache {
	return &CandleCache{
		candles:   candles,
		timeFrame: timeFrame,
	}
}
