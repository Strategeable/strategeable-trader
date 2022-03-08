package indicators

import (
	"cex-bot/helpers"
	"cex-bot/math"
	"cex-bot/types"
)

type RsiIndicatorConfig struct {
	CandlePosition helpers.CandlePosition
	Period         int
}

type RsiIndicator struct {
	Config RsiIndicatorConfig
}

func (r *RsiIndicator) Calculate(input []*types.Candle, _ *types.Position) []float64 {
	if len(input) <= r.Config.Period {
		return make([]float64, 0)
	}

	values := helpers.CandlesToValues(input, r.Config.CandlePosition)

	return math.Rsi(values, r.Config.Period)
}
