package strategy

import (
	"github.com/Stratomicl/Trader/types"
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

func (p *Path) GetTimeFrames() []types.TimeFrame {
	timeFrameMap := make(map[types.TimeFrame]bool)

	for _, tile := range p.Tiles {
		for _, timeFrame := range tile.GetTimeFrames() {
			timeFrameMap[timeFrame] = true
		}
	}

	timeFrames := make([]types.TimeFrame, 0)
	for timeFrame := range timeFrameMap {
		timeFrames = append(timeFrames, timeFrame)
	}
	return timeFrames
}
