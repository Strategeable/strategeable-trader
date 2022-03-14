package main

import (
	"github.com/Strategeable/Trader/database"
	"github.com/Strategeable/Trader/rpcserver"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("error loading .env file")
	}

	databaseHandler := &database.DatabaseHandler{}
	err = databaseHandler.Connect()
	if err != nil {
		panic(err)
	}

	server := rpcserver.NewRpcServer(databaseHandler)
	server.Start()

	// for event := range eventCh {
	// 	switch event.Type {
	// 	case types.POSITION_CREATED:
	// 		fmt.Println(event.Time, "Opened position", event.Data.(*types.Position).Symbol().String())
	// 	case types.POSITION_CLOSED:
	// 		fmt.Println(event.Time, "Closed position", event.Data.(*types.Position).Symbol().String())
	// 	case types.TOTAL_BALANCE_CHANGED:
	// 		fmt.Printf("New balance: %.2f\n", event.Data.(float64))
	// 	}
	// }

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
