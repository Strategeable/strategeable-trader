package indicators

import (
	"cex-bot/helpers"
	"cex-bot/math"
	"cex-bot/types"
)

type ULMLine int

const (
	UPPER ULMLine = iota
	MIDDLE
	LOWER
)

type BollingerBandIndicatorConfig struct {
	CandlePosition helpers.CandlePosition
	Period         int
	DeviationUp    float64
	DeviationDown  float64
	MaType         math.MaType

	Line ULMLine
}

type BollingerBandIndicator struct {
	Config BollingerBandIndicatorConfig
}

func (b *BollingerBandIndicator) Calculate(input []*types.Candle) []float64 {
	values := helpers.CandlesToValues(input, b.Config.CandlePosition)

	lower, middle, upper := math.BBands(values, b.Config.Period, b.Config.DeviationUp, b.Config.DeviationDown, b.Config.MaType)

	switch b.Config.Line {
	case UPPER:
		return upper
	case MIDDLE:
		return middle
	case LOWER:
		return lower
	}

	return []float64{}
}
