package impl

import (
	"math/rand"
	"time"

	"github.com/Stratomicl/Trader/types"
)

type historicalMarketDataProvider struct {
	types.BaseMarketDataProvider

	exchangeImpl types.ExchangeImplementation
	from         time.Time
	until        time.Time

	symbols []types.Symbol

	fullCandleCollection *types.CandleCollection
	tradeCh              chan types.Trade
}

func NewHistoricalMarketDataProvider(
	exchangeImpl types.ExchangeImplementation, from time.Time, until time.Time, symbols []types.Symbol,
) *historicalMarketDataProvider {
	provider := &historicalMarketDataProvider{
		exchangeImpl:         exchangeImpl,
		from:                 from,
		until:                until,
		symbols:              symbols,
		fullCandleCollection: types.NewCandleCollection(),
		tradeCh:              make(chan types.Trade),
	}
	provider.InitCandleCollection()
	return provider
}

func (h *historicalMarketDataProvider) Init() error {
	for _, symbol := range h.symbols {
		candles, err := h.exchangeImpl.GetHistoricalCandles(symbol, types.M1, h.from, h.until)
		if err != nil {
			return err
		}

		h.fullCandleCollection.InitializeTimeFrame(h.exchangeImpl.GetExchange(), symbol, types.M1, candles)
		h.GetCandleCollection().InitializeTimeFrame(h.exchangeImpl.GetExchange(), symbol, types.M1, make([]*types.Candle, 0))
	}

	go func() {
		currentTime := h.from
		for currentTime.Before(h.until) || currentTime.Equal(h.until) {
			trade := types.Trade{
				Symbol:   h.symbols[0],
				Time:     currentTime,
				Price:    float64(rand.Intn(1000)),
				Quantity: 1,
			}

			h.GetCandleCollection().AddTrade(h.exchangeImpl.GetExchange(), h.symbols[0], trade)
			h.tradeCh <- trade

			currentTime = currentTime.Add(1 * time.Minute)
		}
	}()
	return nil
}

func (h *historicalMarketDataProvider) GetTradeCh() chan types.Trade {
	return h.tradeCh
}
