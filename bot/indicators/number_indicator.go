package indicators

import (
	"github.com/Strategeable/Trader/types"
)

type NumberIndicator struct {
	Number float64
}

func (r *NumberIndicator) Calculate(input []*types.Candle, _ *types.Position) []float64 {
	result := make([]float64, 0)

	for i := 0; i < len(input); i++ {
		result = append(result, r.Number)
	}

	return result
}
