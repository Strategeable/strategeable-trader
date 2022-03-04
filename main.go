package main

import (
	"cex-bot/helpers"
	"cex-bot/impl"
	"cex-bot/indicators"
	"cex-bot/trader/types"
	"fmt"
)

func main() {
	var exchangeImpl types.ExchangeImplementation

	exchangeImpl = impl.NewBinanceExchangeImpl()

	err := exchangeImpl.Init()
	if err != nil {
		panic(err)
	}

	candles, err := exchangeImpl.GetCandles(types.Symbol{BaseAsset: "BTC", QuoteAsset: "USDT"}, types.M15, 500)
	if err != nil {
		panic(err)
	}

	heikinAshiCandles := helpers.CandlesCopyToHeikinAshi(candles)

	emaIndicator := indicators.EmaIndicator{
		Config: indicators.EmaIndicatorConfig{
			CandlePosition: helpers.CLOSE,
			Period:         25,
		},
	}

	fmt.Println(emaIndicator.Calculate(heikinAshiCandles))
}
