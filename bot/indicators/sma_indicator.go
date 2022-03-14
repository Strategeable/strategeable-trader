package indicators

import (
	"github.com/Strategeable/Trader/math"
	"github.com/Strategeable/Trader/types"
)

type SmaIndicator struct {
	Source types.Indicator
	Period int
}

func (s *SmaIndicator) Calculate(input []*types.Candle, position *types.Position) []float64 {
	values := s.Source.Calculate(input, position)

	return math.Sma(values, s.Period)
}
