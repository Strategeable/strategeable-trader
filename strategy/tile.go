package strategy

import (
	"cex-bot/handlers"
	"cex-bot/types"
)

type Tile interface {
	HasSignal(candleCollection *handlers.CandleCollection, symbol types.Symbol) (bool, error)
}
