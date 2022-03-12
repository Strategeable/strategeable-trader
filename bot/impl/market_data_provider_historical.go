package impl

import (
	"fmt"
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

	mainCandleCollection *types.CandleCollection
	fullCandleCollection *types.CandleCollection

	tradeCh chan types.Trade
	ackCh   chan string
	closeCh chan string
}

func NewHistoricalMarketDataProvider(
	exchangeImpl types.ExchangeImplementation, from time.Time, until time.Time, symbols []types.Symbol, timeFrames []types.TimeFrame,
	mainCandleCollection *types.CandleCollection,
) *historicalMarketDataProvider {
	provider := &historicalMarketDataProvider{
		exchangeImpl:         exchangeImpl,
		from:                 from,
		until:                until,
		symbols:              symbols,
		timeFrames:           timeFrames,
		mainCandleCollection: mainCandleCollection,
		fullCandleCollection: types.NewCandleCollection(1000000),
		tradeCh:              make(chan types.Trade),
		ackCh:                make(chan string),
		closeCh:              make(chan string),
	}
	provider.InitCandleCollection()
	return provider
}

func (h *historicalMarketDataProvider) Init() error {
	candleMapping := make(map[string]map[int64]*types.Candle)

	for _, symbol := range h.symbols {
		cache := h.mainCandleCollection.GetCache(h.exchangeImpl.GetExchange(), symbol, types.M1)

		dateRanges := make([]types.DateRange, 0)

		if cache != nil {
			var currentDateRange *types.DateRange

			currentTime := h.from
			for currentTime.Before(h.until) || currentTime.Equal(h.until) {
				if cache.GetCandleAt(currentTime) == nil {
					if currentDateRange == nil {
						currentDateRange = &types.DateRange{
							From: currentTime,
						}
					}
				} else {
					if currentDateRange != nil {
						currentDateRange.To = currentTime.Add(-1 * time.Minute)
						dateRanges = append(dateRanges, *currentDateRange)
						currentDateRange = nil
					}
				}

				currentTime = currentTime.Add(time.Minute)
			}
			if currentDateRange != nil {
				currentDateRange.To = currentTime
				dateRanges = append(dateRanges, *currentDateRange)
			}
		} else {
			dateRanges = append(dateRanges, types.DateRange{
				From: h.from,
				To:   h.until,
			})
		}

		fmt.Println(len(dateRanges))

		var candles []*types.Candle
		if len(dateRanges) == 0 {
			candles = make([]*types.Candle, 0)
			currentTime := h.from
			for currentTime.Before(h.until) || currentTime.Equal(h.until) {
				candles = append(candles, cache.GetCandleAt(currentTime))
				currentTime = currentTime.Add(time.Minute)
			}
		} else {
			c, err := h.exchangeImpl.GetHistoricalCandles(symbol, types.M1, h.from, h.until)
			if err != nil {
				return err
			}
			candles = c
		}

		if cache == nil {
			h.mainCandleCollection.InitializeTimeFrame(h.exchangeImpl.GetExchange(), symbol, types.M1, candles)
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

		close(h.closeCh)
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

func (h *historicalMarketDataProvider) GetCloseCh() chan string {
	return h.closeCh
}
