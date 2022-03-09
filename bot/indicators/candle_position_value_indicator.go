package indicators

import (
	"github.com/Stratomicl/Trader/helpers"
	"github.com/Stratomicl/Trader/types"
)

type CandlePositionValueIndicatorConfig struct {
	CandlePosition helpers.CandlePosition
}

type CandlePositionValueIndicator struct {
	Config CandlePositionValueIndicatorConfig
}

func (c *CandlePositionValueIndicator) Calculate(input []*types.Candle, _ *types.Position) []float64 {
	return helpers.CandlesToValues(input, c.Config.CandlePosition)
}
