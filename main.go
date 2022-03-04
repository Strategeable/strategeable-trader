package main

import (
	"cex-bot/helpers"
	"cex-bot/impl"
	"cex-bot/indicators"
	"cex-bot/strategy"
	"cex-bot/types"
	"fmt"
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
			candleCollection.RegisterSymbol(symbol)

			candles, err := exchangeImpl.GetCandles(symbol, timeFrame, 1000)
			if err != nil {
				panic(err)
			}

			candleCollection.InitializeTimeFrame(symbol, timeFrame, candles)
		}
	}

	keepaliveCh := make(chan string)

	path := &strategy.Path{
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
							Number: 32,
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

	_, err = exchangeImpl.WatchTrades(symbols, func(trade types.Trade) {
		candleCollection.AddTrade(trade.Symbol, trade)

		cache := candleCollection.GetCache(trade.Symbol, types.M1)

		signal, err := path.HasSignal(candleCollection, trade.Symbol)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(trade.Symbol.String(), cache.GetCurrentRate(), signal)
	}, func(err error) {
		fmt.Println(err)
	})
	if err != nil {
		panic(err)
	}

	<-keepaliveCh
}
