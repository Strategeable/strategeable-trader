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

	candleCollection := types.NewCandleCollection()

	for _, timeFrame := range []types.TimeFrame{types.M1, types.M5, types.H1} {
		for _, symbol := range symbols {
			candleCollection.RegisterSymbol(types.BINANCE, symbol)

			candles, err := exchangeImpl.GetCandles(symbol, timeFrame, 1000)
			if err != nil {
				panic(err)
			}

			candleCollection.InitializeTimeFrame(types.BINANCE, symbol, timeFrame, candles)
		}
	}

	keepaliveCh := make(chan string)

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

	bot := handlers.NewBot("test bot", handlers.ExchangeDetails{}, strategy.Strategy{
		Exchange:    types.BINANCE,
		Symbols:     symbols,
		BuyPaths:    []*strategy.Path{buyPath},
		SellPaths:   []*strategy.Path{sellPath},
		BuyCooldown: 60 * time.Second,
	}, make([]*types.Position, 0), make([]*types.Position, 0), candleCollection)

	go bot.RunLoop()

	_, err = exchangeImpl.WatchTrades(symbols, func(trade types.Trade) {
		candleCollection.AddTrade(types.BINANCE, trade.Symbol, trade)

		// cache := candleCollection.GetCache(trade.Symbol, types.M1)

		// signal, err := path.HasSignal(candleCollection, trade.Symbol)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

		// fmt.Println(trade.Symbol.String(), cache.GetCurrentRate(), signal)

		bot.TradeCh <- trade
	}, func(err error) {
		fmt.Println(err)
	})
	if err != nil {
		panic(err)
	}

	<-keepaliveCh
}
