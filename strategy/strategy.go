package strategy

import (
	"cex-bot/types"
	"time"
)

type Strategy struct {
	BuyPaths         []*Path
	SellPaths        []*Path
	Symbols          []types.Symbol
	DefaultTimeFrame types.TimeFrame
	BuyCooldown      time.Duration
}

func (s *Strategy) GetQuoteAsset() string {
	return s.Symbols[0].QuoteAsset
}

func (s *Strategy) HasBuySignal(candleCollection *types.CandleCollection, symbol types.Symbol) (bool, error) {
	return hasSignal(s.BuyPaths, candleCollection, symbol)
}

func (s *Strategy) HasSellSignal(candleCollection *types.CandleCollection, symbol types.Symbol) (bool, error) {
	return hasSignal(s.SellPaths, candleCollection, symbol)
}

func hasSignal(paths []*Path, candleCollection *types.CandleCollection, symbol types.Symbol) (bool, error) {
	var pathError error

	for _, path := range paths {
		signal, err := path.HasSignal(candleCollection, symbol)
		if err != nil {
			pathError = err
			continue
		}

		if signal {
			return true, nil
		}
	}

	return false, pathError
}
