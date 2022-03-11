package indicators

import (
	"github.com/Stratomicl/Trader/helpers"
	"github.com/Stratomicl/Trader/types"
)

type CandlePositionValueIndicator struct {
	CandlePosition helpers.CandlePosition
}

func (c *CandlePositionValueIndicator) Calculate(input []*types.Candle, position *types.Position) []float64 {
	return helpers.CandlesToValues(input, c.CandlePosition)
}
