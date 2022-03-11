package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/Stratomicl/Trader/handlers"
	"github.com/Stratomicl/Trader/impl"
	"github.com/Stratomicl/Trader/strategy"
	"github.com/Stratomicl/Trader/types"
)

func main() {
	content, err := ioutil.ReadFile("strategy.json")
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)

	strategy, err := strategy.StrategyFromJson(text)
	if err != nil {
		panic(err)
	}

	var exchangeImpl types.ExchangeImplementation

	exchangeImpl = impl.NewBinanceExchangeImpl()

	err = exchangeImpl.Init()
	if err != nil {
		panic(err)
	}

	from, _ := time.Parse("2006-01-02 15:04", "2022-03-06 22:00")
	to, _ := time.Parse("2006-01-02 15:04", "2022-03-10 18:30")

	marketDataProvider := impl.NewHistoricalMarketDataProvider(exchangeImpl, from, to, strategy.Symbols, strategy.GetTimeFrames())

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
			fmt.Println(event.Time, "Opened position", event.Data.(*types.Position).Symbol().String())
		case types.POSITION_CLOSED:
			fmt.Println(event.Time, "Closed position", event.Data.(*types.Position).Symbol().String())
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
