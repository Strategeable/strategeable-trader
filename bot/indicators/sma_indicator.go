package indicators

import (
	"github.com/Stratomicl/Trader/math"
	"github.com/Stratomicl/Trader/types"
)

type SmaIndicator struct {
	Source types.Indicator
	Period int
}

func (s *SmaIndicator) Calculate(input []*types.Candle, position *types.Position) []float64 {
	values := s.Source.Calculate(input, position)

	return math.Sma(values, s.Period)
}
