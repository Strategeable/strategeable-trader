package strategy

import (
	"cex-bot/handlers"
	"cex-bot/types"
	"time"
)

type Strategy struct {
	BuyPaths         []*Path
	SellPaths        []*Path
	DefaultTimeFrame types.TimeFrame
	BuyCooldown      time.Duration
}

func (s *Strategy) HasBuySignal(candleCollection *handlers.CandleCollection, symbol types.Symbol) (bool, error) {
	return hasSignal(s.BuyPaths, candleCollection, symbol)
}

func (s *Strategy) HasSellSignal(candleCollection *handlers.CandleCollection, symbol types.Symbol) (bool, error) {
	return hasSignal(s.SellPaths, candleCollection, symbol)
}

func hasSignal(paths []*Path, candleCollection *handlers.CandleCollection, symbol types.Symbol) (bool, error) {
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
