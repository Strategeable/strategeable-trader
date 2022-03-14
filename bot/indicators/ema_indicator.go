package indicators

import (
	"github.com/Strategeable/Trader/math"
	"github.com/Strategeable/Trader/types"
)

type EmaIndicator struct {
	Source types.Indicator
	Period int
}

func (e *EmaIndicator) Calculate(input []*types.Candle, position *types.Position) []float64 {
	values := e.Source.Calculate(input, position)

	return math.Ema(values, e.Period)
}
