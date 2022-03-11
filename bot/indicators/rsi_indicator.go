package indicators

import (
	"github.com/Stratomicl/Trader/math"
	"github.com/Stratomicl/Trader/types"
)

type RsiIndicator struct {
	Source types.Indicator
	Period int
}

func (r *RsiIndicator) Calculate(input []*types.Candle, position *types.Position) []float64 {
	if len(input) <= r.Period {
		return make([]float64, 0)
	}

	values := r.Source.Calculate(input, position)

	return math.Rsi(values, r.Period)
}
