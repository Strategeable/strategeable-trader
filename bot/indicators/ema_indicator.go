package indicators

import (
	"github.com/Stratomicl/Trader/helpers"
	"github.com/Stratomicl/Trader/math"
	"github.com/Stratomicl/Trader/types"
)

type EmaIndicatorConfig struct {
	CandlePosition helpers.CandlePosition
	Period         int
}

type EmaIndicator struct {
	Config EmaIndicatorConfig
}

func (e *EmaIndicator) Calculate(input []*types.Candle, _ *types.Position) []float64 {
	values := helpers.CandlesToValues(input, e.Config.CandlePosition)

	return math.Ema(values, e.Config.Period)
}
