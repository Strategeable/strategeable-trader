package strategy

import "github.com/Stratomicl/Trader/types"

type AnySignalTile struct {
	SignalTiles []SignalTile
	Amount      int
}

func (a *AnySignalTile) HasSignal(candleCollection *types.CandleCollection, symbol types.Symbol, exchange types.Exchange, position *types.Position) (bool, error) {
	var err error

	totalSignalsMet := 0

	for _, tile := range a.SignalTiles {
		signal, sigErr := tile.HasSignal(candleCollection, symbol, exchange, position)
		if err != nil {
			err = sigErr
			continue
		}

		if signal {
			totalSignalsMet++

			if totalSignalsMet >= a.Amount {
				return true, err
			}
		}
	}

	return false, err
}
