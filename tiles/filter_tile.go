package tiles

import (
	"cex-bot/handlers"
	"cex-bot/types"
)

type FilterTile struct {
	symbols   map[string]bool
	whitelist bool
}

func (w *FilterTile) HasSignal(candleCollection *handlers.CandleCollection, symbol types.Symbol) (bool, error) {
	return (w.whitelist && w.symbols[symbol.String()]) || (!w.whitelist && !w.symbols[symbol.String()]), nil
}
