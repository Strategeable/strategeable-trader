package indicators

import (
	"github.com/Strategeable/Trader/types"
)

type PositionHoldTimeIndicator struct{}

func (p *PositionHoldTimeIndicator) Calculate(input []*types.Candle, position *types.Position) []float64 {
	result := make([]float64, 0)
	if position == nil {
		return result
	}

	for i := 0; i < len(input); i++ {
		holdTime := input[i].OpenTime.Sub(position.OpenTime())
		if input[i].OpenTime.Before(position.OpenTime()) {
			holdTime = 0
		}
		result = append(result, holdTime.Seconds())
	}

	return result
}
