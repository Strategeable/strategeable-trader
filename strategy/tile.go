package strategy

import (
	"cex-bot/types"
)

type Tile interface {
	HasSignal(candleCollection *types.CandleCollection, symbol types.Symbol) (bool, error)
}
