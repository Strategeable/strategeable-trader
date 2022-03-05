package indicators

import (
	"cex-bot/helpers"
	"cex-bot/math"
	"cex-bot/types"
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
