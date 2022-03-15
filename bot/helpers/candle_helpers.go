package helpers

import (
	"math"

	"github.com/Strategeable/Trader/types"
)

type CandlePosition string

const (
	OPEN   CandlePosition = "OPEN"
	HIGH   CandlePosition = "HIGH"
	LOW    CandlePosition = "LOW"
	CLOSE  CandlePosition = "CLOSE"
	VOLUME CandlePosition = "VOLUME"
)

// Converts candles to a float slice based on the
// desired position, like candle open, candle close...
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

// Returns a slice of new candle objects in Heikin Ashi values
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
			candle.Exchange,
			candle.Symbol,
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
