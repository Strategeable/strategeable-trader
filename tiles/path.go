package tiles

import (
	"cex-bot/handlers"
	"cex-bot/types"
)

type Path struct {
	Tiles []Tile
}

func (p *Path) HasSignal(candleCollection *handlers.CandleCollection, symbol types.Symbol) (bool, error) {
	for _, tile := range p.Tiles {
		signal, err := tile.HasSignal(candleCollection, symbol)
		if err != nil {
			return false, err
		}
		if !signal {
			return false, nil
		}
	}

	return true, nil
}
