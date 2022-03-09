package indicators

import (
	"github.com/Stratomicl/Trader/helpers"
	"github.com/Stratomicl/Trader/math"
	"github.com/Stratomicl/Trader/types"
)

type SmaIndicatorConfig struct {
	CandlePosition helpers.CandlePosition
	Period         int
}

type SmaIndicator struct {
	Config SmaIndicatorConfig
}

func (s *SmaIndicator) Calculate(input []*types.Candle, _ *types.Position) []float64 {
	values := helpers.CandlesToValues(input, s.Config.CandlePosition)

	return math.Sma(values, s.Config.Period)
}
