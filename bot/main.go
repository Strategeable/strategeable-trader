package main

import (
	"fmt"
	"time"

	"github.com/Stratomicl/Trader/handlers"
	"github.com/Stratomicl/Trader/helpers"
	"github.com/Stratomicl/Trader/impl"
	"github.com/Stratomicl/Trader/indicators"
	"github.com/Stratomicl/Trader/strategy"
	"github.com/Stratomicl/Trader/types"
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
			&strategy.AnySignalTile{
				SignalTiles: []strategy.SignalTile{
					{
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
				Amount: 1,
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
		BuySize:     100,
	}

	from, _ := time.Parse("2006-01-02 15:04", "2022-03-06 22:00")
	to, _ := time.Parse("2006-01-02 15:04", "2022-03-09 15:30")

	marketDataProvider := impl.NewHistoricalMarketDataProvider(exchangeImpl, from, to, symbols, strategy.GetTimeFrames())

	positionHandler := impl.NewSimulatedPositionHandler(1000, make([]*types.Position, 0))

	eventCh := make(chan types.PositionHandlerEvent, 5)
	positionHandler.SubscribeEvents(eventCh)

	engine := handlers.NewEngine(strategy, marketDataProvider, positionHandler)

	go func() {
		err = engine.Start()
		if err != nil {
			panic(err)
		}
	}()

	for event := range eventCh {
		switch event.Type {
		case types.POSITION_CREATED:
			fmt.Println(event.Time, "Opened position")
		case types.POSITION_CLOSED:
			fmt.Println(event.Time, "Closed position")
		case types.TOTAL_BALANCE_CHANGED:
			fmt.Printf("New balance: %.2f\n", event.Data.(float64))
		}
	}

	keepaliveCh := make(chan string)
	<-keepaliveCh

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
