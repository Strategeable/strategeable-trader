package strategy

import (
	"time"

	"github.com/Strategeable/Trader/types"
)

type Strategy struct {
	Exchange         types.Exchange
	BuyPaths         []*Path
	SellPaths        []*Path
	Symbols          []types.Symbol
	DefaultTimeFrame types.TimeFrame
	BuyCooldown      time.Duration
	BuySize          float64
}

func (s *Strategy) GetTimeFrames() []types.TimeFrame {
	timeFrameMap := make(map[types.TimeFrame]bool)

	paths := make([]*Path, 0)
	paths = append(paths, s.BuyPaths...)
	paths = append(paths, s.SellPaths...)

	for _, path := range paths {
		for _, timeFrame := range path.GetTimeFrames() {
			timeFrameMap[timeFrame] = true
		}
	}

	timeFrames := make([]types.TimeFrame, 0)
	for timeFrame := range timeFrameMap {
		timeFrames = append(timeFrames, timeFrame)
	}
	return timeFrames
}

func (s *Strategy) GetQuoteAsset() string {
	return s.Symbols[0].QuoteAsset
}

func (s *Strategy) HasBuySignal(candleCollection *types.CandleCollection, symbol types.Symbol) (bool, error) {
	return s.hasSignal(s.BuyPaths, candleCollection, symbol, nil)
}

func (s *Strategy) HasSellSignal(candleCollection *types.CandleCollection, symbol types.Symbol, position *types.Position) (bool, error) {
	return s.hasSignal(s.SellPaths, candleCollection, symbol, position)
}

func (s *Strategy) hasSignal(paths []*Path, candleCollection *types.CandleCollection, symbol types.Symbol, position *types.Position) (bool, error) {
	var pathError error

	for _, path := range paths {
		signal, err := path.HasSignal(candleCollection, symbol, s.Exchange, position)
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
