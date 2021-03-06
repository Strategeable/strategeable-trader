package strategy

import (
	"github.com/Strategeable/Trader/types"
)

type Tile interface {
	HasSignal(candleCollection *types.CandleCollection, symbol types.Symbol, exchange types.Exchange, position *types.Position) (bool, error)
	GetTimeFrames() []types.TimeFrame
}
