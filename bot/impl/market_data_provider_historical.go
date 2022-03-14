package impl

import (
	"fmt"
	"time"

	"github.com/Strategeable/Trader/database"
	"github.com/Strategeable/Trader/types"
)

type historicalMarketDataProvider struct {
	types.BaseMarketDataProvider

	databaseHandler *database.DatabaseHandler

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

	stopped bool
}

func NewHistoricalMarketDataProvider(
	exchangeImpl types.ExchangeImplementation, from time.Time, until time.Time, symbols []types.Symbol, timeFrames []types.TimeFrame,
	mainCandleCollection *types.CandleCollection, databaseHandler *database.DatabaseHandler,
) *historicalMarketDataProvider {
	provider := &historicalMarketDataProvider{
		exchangeImpl:         exchangeImpl,
		from:                 from,
		until:                until,
		symbols:              symbols,
		timeFrames:           timeFrames,
		mainCandleCollection: mainCandleCollection,
		fullCandleCollection: types.NewCandleCollection(-1),
		tradeCh:              make(chan types.Trade),
		ackCh:                make(chan string),
		closeCh:              make(chan string),
		databaseHandler:      databaseHandler,
	}
	provider.InitCandleCollection()
	return provider
}

func (h *historicalMarketDataProvider) Init() error {
	candleMapping := make(map[string]map[int64]*types.Candle)

	for _, symbol := range h.symbols {
		if h.stopped {
			return nil
		}

		cache := h.mainCandleCollection.GetCache(h.exchangeImpl.GetExchange(), symbol, types.M1)

		if cache == nil {
			candles, err := h.databaseHandler.GetCandles(h.exchangeImpl.GetExchange(), symbol)
			if err != nil {
				return err
			}

			firstOpenTime, err := h.exchangeImpl.GetFirstCandleTime(symbol)
			if err != nil {
				return err
			}

			latestCandleTime := firstOpenTime
			if len(candles) > 0 {
				latestCandleTime = candles[len(candles)-1].OpenTime.Add(1 * time.Minute)
			}

			if latestCandleTime.Before(time.Now().Add(-1 * time.Hour)) {
				// Load all candles and save to database
				fmt.Printf("Updating DB cache for %s.\n", symbol.String())

				fmt.Printf("Loading all 1m candles for %s, starting at %s. This may take a while.\n", symbol.String(), latestCandleTime.Format(time.RFC822))

				candleCh := make(chan []*types.Candle)

				go func() {
					for candles := range candleCh {
						if len(candles) == 0 {
							continue
						}
						fmt.Printf("Saving %d candles for %s.\n", len(candles), symbol.String())
						err := h.databaseHandler.SaveCandles(candles)
						if err != nil {
							panic(err)
						}
					}
				}()

				addedCandles, err := h.exchangeImpl.GetHistoricalCandles(symbol, types.M1, latestCandleTime, time.Now(), candleCh)
				if err != nil {
					return err
				}

				candles = append(candles, addedCandles...)
			}

			h.mainCandleCollection.InitializeTimeFrame(h.exchangeImpl.GetExchange(), symbol, types.M1, candles)
			cache = h.mainCandleCollection.GetCache(h.exchangeImpl.GetExchange(), symbol, types.M1)
		}

		candles := make([]*types.Candle, 0)
		currentTime := h.from
		for currentTime.Before(h.until) || currentTime.Equal(h.until) {
			candle := cache.GetCandleAt(currentTime)
			currentTime = currentTime.Add(time.Minute)
			if candle == nil {
				continue
			}
			candles = append(candles, candle)
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
		for !h.stopped && (currentTime.Before(h.until) || currentTime.Equal(h.until)) {
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

func (h *historicalMarketDataProvider) Close() {
	h.stopped = true
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
