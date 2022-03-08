package main

import (
	"cex-bot/handlers"
	"cex-bot/helpers"
	"cex-bot/impl"
	"cex-bot/indicators"
	"cex-bot/strategy"
	"cex-bot/types"
	"fmt"
	"time"
)

func main() {
	var exchangeImpl types.ExchangeImplementation

	exchangeImpl = impl.NewBinanceExchangeImpl()

	err := exchangeImpl.Init()
	if err != nil {
		panic(err)
	}

	eth := types.Symbol{BaseAsset: "ETH", QuoteAsset: "USDT"}
	// btc := types.Symbol{BaseAsset: "BTC", QuoteAsset: "USDT"}

	symbols := make([]types.Symbol, 0)
	symbols = append(symbols, eth)
	// symbols = append(symbols, btc)

	buyPath := &strategy.Path{
		Tiles: []strategy.Tile{
			&strategy.SignalTile{
				IndicatorA: strategy.IndicatorSettings{
					Indicator: &indicators.RsiIndicator{
						Config: indicators.RsiIndicatorConfig{
							CandlePosition: helpers.CLOSE,
							Period:         14,
						},
					},
					RealTime:    true,
					CandlesBack: 0,
					TimeFrame:   types.M1,
				},
				IndicatorB: strategy.IndicatorSettings{
					Indicator: &indicators.NumberIndicator{
						Config: indicators.NumberIndicatorConfig{
							Number: 47,
						},
					},
					RealTime:    true,
					CandlesBack: 0,
					TimeFrame:   types.M1,
				},
				Operand:     strategy.GREATER_THAN,
				Persistence: 1,
			},
		},
	}
	sellPath := &strategy.Path{
		Tiles: []strategy.Tile{
			&strategy.SignalTile{
				IndicatorA: strategy.IndicatorSettings{
					Indicator: &indicators.RsiIndicator{
						Config: indicators.RsiIndicatorConfig{
							CandlePosition: helpers.CLOSE,
							Period:         14,
						},
					},
					RealTime:    true,
					CandlesBack: 0,
					TimeFrame:   types.M1,
				},
				IndicatorB: strategy.IndicatorSettings{
					Indicator: &indicators.NumberIndicator{
						Config: indicators.NumberIndicatorConfig{
							Number: 47,
						},
					},
					RealTime:    true,
					CandlesBack: 0,
					TimeFrame:   types.M1,
				},
				Operand:     strategy.LOWER_THAN,
				Persistence: 1,
			},
		},
	}

	strategy := strategy.Strategy{
		Exchange:    types.BINANCE,
		Symbols:     symbols,
		BuyPaths:    []*strategy.Path{buyPath},
		SellPaths:   []*strategy.Path{sellPath},
		BuyCooldown: 60 * time.Second,
	}

	from, _ := time.Parse("2006-01-02 15:04", "2022-03-08 22:00")
	to, _ := time.Parse("2006-01-02 15:04", "2022-03-08 22:30")

	marketDataProvider := impl.NewHistoricalMarketDataProvider(exchangeImpl, from, to, symbols)

	positionHandler := impl.NewSimulatedPositionHandler(1000, make([]*types.Position, 0))

	eventCh := make(chan types.PositionHandlerEvent, 5)
	positionHandler.SubscribeEvents(eventCh)

	engine := handlers.NewEngine(strategy, marketDataProvider, positionHandler)

	err = engine.Start()
	if err != nil {
		panic(err)
	}

	for event := range eventCh {
		fmt.Println(event)
	}

	// _, err = exchangeImpl.WatchTrades(symbols, func(trade types.Trade) {
	// 	candleCollection.AddTrade(types.BINANCE, trade.Symbol, trade)

	// 	// cache := candleCollection.GetCache(trade.Symbol, types.M1)

	// 	// signal, err := path.HasSignal(candleCollection, trade.Symbol)
	// 	// if err != nil {
	// 	// 	fmt.Println(err)
	// 	// 	return
	// 	// }

	// 	// fmt.Println(trade.Symbol.String(), cache.GetCurrentRate(), signal)

	// 	bot.TradeCh <- trade
	// }, func(err error) {
	// 	fmt.Println(err)
	// })
	// if err != nil {
	// 	panic(err)
	// }
}
