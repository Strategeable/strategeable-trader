package indicators

import (
	"github.com/Strategeable/Trader/types"
)

type PositionChangeIndicator struct{}

func (p *PositionChangeIndicator) Calculate(input []*types.Candle, position *types.Position) []float64 {
	result := make([]float64, 0)
	if position == nil {
		return result
	}

	for i := 0; i < len(input); i++ {
		changePercentage := position.ChangePercentage(input[i].Close)
		if input[i].OpenTime.Before(position.OpenTime()) {
			changePercentage = 0
		}
		result = append(result, changePercentage)
	}

	return result
}
