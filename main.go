package main

import (
	"cex-bot/impl"
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

	btc := types.Symbol{BaseAsset: "BTC", QuoteAsset: "USDT"}
	candles, err := exchangeImpl.GetCandles(btc, types.M1, 500)
	if err != nil {
		panic(err)
	}

	cache := types.NewCandleCache(candles, types.M1, 500)

	keepaliveCh := make(chan string)

	exchangeImpl.WatchTrades([]types.Symbol{btc}, func(trade types.Trade) {
		fmt.Println(cache.GetSize(), cache.AddTrade(trade.Price, trade.Quantity, trade.Time), cache.GetCurrentCandle())
	}, func(err error) {
		fmt.Println(err)
	})

	<-keepaliveCh
}
