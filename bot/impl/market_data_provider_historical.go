package impl

import (
	"time"

	"github.com/Stratomicl/Trader/types"
)

type historicalMarketDataProvider struct {
	types.BaseMarketDataProvider

	exchangeImpl types.ExchangeImplementation
	from         time.Time
	until        time.Time

	symbols    []types.Symbol
	timeFrames []types.TimeFrame

	fullCandleCollection *types.CandleCollection
	tradeCh              chan types.Trade
	ackCh                chan string
}

func NewHistoricalMarketDataProvider(
	exchangeImpl types.ExchangeImplementation, from time.Time, until time.Time, symbols []types.Symbol, timeFrames []types.TimeFrame,
) *historicalMarketDataProvider {
	provider := &historicalMarketDataProvider{
		exchangeImpl:         exchangeImpl,
		from:                 from,
		until:                until,
		symbols:              symbols,
		timeFrames:           timeFrames,
		fullCandleCollection: types.NewCandleCollection(),
		tradeCh:              make(chan types.Trade),
		ackCh:                make(chan string),
	}
	provider.InitCandleCollection()
	return provider
}

func (h *historicalMarketDataProvider) Init() error {
	candleMapping := make(map[string]map[int64]*types.Candle)

	for _, symbol := range h.symbols {
		candles, err := h.exchangeImpl.GetHistoricalCandles(symbol, types.M1, h.from, h.until)
		if err != nil {
			return err
		}

		mapping := make(map[int64]*types.Candle)

		for _, candle := range candles {
			mapping[candle.OpenTime.Unix()] = candle
		}

		candleMapping[symbol.String()] = mapping

		h.fullCandleCollection.InitializeTimeFrame(h.exchangeImpl.GetExchange(), symbol, types.M1, candles)

		for _, timeFrame := range h.timeFrames {
			h.GetCandleCollection().InitializeTimeFrame(h.exchangeImpl.GetExchange(), symbol, timeFrame, make([]*types.Candle, 0))
		}
	}

	go func() {
		currentTime := h.from
		for currentTime.Before(h.until) || currentTime.Equal(h.until) {
			for _, symbol := range h.symbols {
				candle := candleMapping[symbol.String()][currentTime.Unix()]
				if candle == nil {
					continue
				}

				trade := types.Trade{
					Symbol:   symbol,
					Time:     currentTime,
					Price:    candle.Close,
					Quantity: candle.Volume,
				}

				h.GetCandleCollection().AddTrade(h.exchangeImpl.GetExchange(), symbol, trade)
				h.tradeCh <- trade
				<-h.ackCh
			}

			currentTime = currentTime.Add(1 * time.Minute)
		}
	}()
	return nil
}

func (h *historicalMarketDataProvider) GetTradeCh() chan types.Trade {
	return h.tradeCh
}

func (h *historicalMarketDataProvider) RequiresAcks() bool {
	return true
}

func (h *historicalMarketDataProvider) GetAckCh() chan string {
	return h.ackCh
}
