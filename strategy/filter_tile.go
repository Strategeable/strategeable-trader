package strategy

import (
	"cex-bot/types"
)

type FilterTile struct {
	symbols   map[string]bool
	whitelist bool
}

func (w *FilterTile) HasSignal(candleCollection *types.CandleCollection, symbol types.Symbol, exchange types.Exchange) (bool, error) {
	return (w.whitelist && w.symbols[symbol.String()]) || (!w.whitelist && !w.symbols[symbol.String()]), nil
}
