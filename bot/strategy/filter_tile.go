package strategy

import (
	"github.com/Strategeable/Trader/types"
)

type FilterTile struct {
	symbols   map[string]bool
	whitelist bool
}

func (w *FilterTile) HasSignal(candleCollection *types.CandleCollection, symbol types.Symbol, exchange types.Exchange, position *types.Position) (bool, error) {
	return (w.whitelist && w.symbols[symbol.String()]) || (!w.whitelist && !w.symbols[symbol.String()]), nil
}
