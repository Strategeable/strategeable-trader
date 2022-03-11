package indicators

import (
	"github.com/Stratomicl/Trader/math"
	"github.com/Stratomicl/Trader/types"
)

type ULMLine int

const (
	UPPER ULMLine = iota
	MIDDLE
	LOWER
)

type BollingerBandIndicator struct {
	Source        types.Indicator
	Period        int
	DeviationUp   float64
	DeviationDown float64
	MaType        math.MaType

	Line ULMLine
}

func (b *BollingerBandIndicator) Calculate(input []*types.Candle, position *types.Position) []float64 {
	values := b.Source.Calculate(input, position)

	lower, middle, upper := math.BBands(values, b.Period, b.DeviationUp, b.DeviationDown, b.MaType)

	switch b.Line {
	case UPPER:
		return upper
	case MIDDLE:
		return middle
	case LOWER:
		return lower
	}

	return []float64{}
}
