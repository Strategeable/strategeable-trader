package helpers

import (
	"math"

	"github.com/Stratomicl/Trader/types"
)

type CandlePosition int

const (
	OPEN CandlePosition = iota
	HIGH
	LOW
	CLOSE
	VOLUME
)

func CandlesToValues(candles []*types.Candle, candlePosition CandlePosition) []float64 {
	result := make([]float64, 0)

	for _, candle := range candles {
		switch candlePosition {
		case OPEN:
			result = append(result, candle.Open)
		case HIGH:
			result = append(result, candle.High)
		case LOW:
			result = append(result, candle.Low)
		case CLOSE:
			result = append(result, candle.Close)
		case VOLUME:
			result = append(result, candle.Volume)
		}
	}

	return result
}

func CandlesCopyToHeikinAshi(candles []*types.Candle) []*types.Candle {
	result := make([]*types.Candle, 0)

	for i, candle := range candles {
		prevOpen := candle.Open
		prevClose := candle.Close
		if i > 0 {
			prevOpen = result[i-1].Open
			prevClose = result[i-1].Close
		}

		result = append(result, types.NewCandle(
			candle.OpenTime,
			candle.CloseTime,
			0.5*(prevOpen+prevClose),
			math.Max(math.Max(candle.High, candle.Open), candle.Close),
			math.Min(math.Min(candle.Low, candle.Open), candle.Close),
			0.25*(candle.Open+candle.High+candle.Low+candle.Close),
			candle.Volume,
		))
	}

	return result
}