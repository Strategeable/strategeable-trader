package strategy

import (
	"cex-bot/types"
)

type Path struct {
	Tiles []Tile
}

func (p *Path) HasSignal(candleCollection *types.CandleCollection, symbol types.Symbol, exchange types.Exchange, position *types.Position) (bool, error) {
	for _, tile := range p.Tiles {
		signal, err := tile.HasSignal(candleCollection, symbol, exchange, position)
		if err != nil {
			return false, err
		}
		if !signal {
			return false, nil
		}
	}

	return true, nil
}
