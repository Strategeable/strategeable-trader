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

	eth := types.Symbol{BaseAsset: "ETH", QuoteAsset: "USDT"}
	btc := types.Symbol{BaseAsset: "BTC", QuoteAsset: "USDT"}

	symbols := make([]types.Symbol, 0)
	symbols = append(symbols, eth)
	symbols = append(symbols, btc)

	caches := make(map[string]map[types.TimeFrame]*types.CandleCache)

	for _, timeFrame := range []types.TimeFrame{types.M1, types.M5, types.H1} {
		for _, symbol := range symbols {
			candles, err := exchangeImpl.GetCandles(symbol, timeFrame, 500)
			if err != nil {
				panic(err)
			}

			cache := types.NewCandleCache(candles, timeFrame, 500)

			if caches[symbol.String()] == nil {
				caches[symbol.String()] = make(map[types.TimeFrame]*types.CandleCache)
			}
			caches[symbol.String()][timeFrame] = cache
		}
	}

	keepaliveCh := make(chan string)

	_, err = exchangeImpl.WatchTrades(symbols, func(trade types.Trade) {
		timeFrameCaches := caches[trade.Symbol.String()]

		for timeFrame, cache := range timeFrameCaches {
			newCandle := cache.AddTrade(trade.Price, trade.Quantity, trade.Time)

			fmt.Println(trade.Symbol.String(), newCandle, cache.GetCurrentRate(), timeFrame)
		}
	}, func(err error) {
		fmt.Println(err)
	})
	if err != nil {
		panic(err)
	}

	<-keepaliveCh
}
