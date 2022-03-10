package impl

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Stratomicl/Trader/types"

	"github.com/adshao/go-binance/v2"
)

type binanceExchangeImpl struct {
	client *binance.Client

	exchangeInfoService *binance.ExchangeInfoService
	klinesService       *binance.KlinesService
	bookTickerService   *binance.ListBookTickersService
}

func NewBinanceExchangeImpl() *binanceExchangeImpl {
	return &binanceExchangeImpl{}
}

func (b *binanceExchangeImpl) Init() error {
	b.client = binance.NewClient("", "")

	b.exchangeInfoService = b.client.NewExchangeInfoService()
	b.klinesService = b.client.NewKlinesService()
	b.bookTickerService = b.client.NewListBookTickersService()
	return nil
}

func (b *binanceExchangeImpl) GetExchange() types.Exchange {
	return types.BINANCE
}

func (b *binanceExchangeImpl) GetAvailableTimeFrames() []types.TimeFrame {
	return make([]types.TimeFrame, 0)
}

func (b *binanceExchangeImpl) FormatTimeFrame(timeFrame types.TimeFrame) string {
	return string(timeFrame)
}

func (b *binanceExchangeImpl) GetSymbols() ([]types.Symbol, error) {
	symbols := make([]types.Symbol, 0)

	symbolsResp, err := b.exchangeInfoService.Do(context.Background())
	if err != nil {
		return nil, err
	}

	for _, symbol := range symbolsResp.Symbols {
		symbols = append(symbols, types.NewSymbol(symbol.BaseAsset, symbol.BaseAssetPrecision, symbol.QuoteAsset, symbol.QuoteAssetPrecision))
	}

	return symbols, nil
}

func (b *binanceExchangeImpl) GetUniqueQuoteAssets() ([]string, error) {
	symbols, err := b.GetSymbols()
	if err != nil {
		return nil, err
	}

	uniqueQuoteAssets := make([]string, 0)
	existing := make(map[string]bool)

	for _, symbol := range symbols {
		if !existing[symbol.QuoteAsset] {
			uniqueQuoteAssets = append(uniqueQuoteAssets, symbol.QuoteAsset)
			existing[symbol.QuoteAsset] = true
		}
	}

	return uniqueQuoteAssets, nil

}

func (b *binanceExchangeImpl) FormatSymbol(symbol types.Symbol) string {
	return fmt.Sprintf("%s%s", symbol.BaseAsset, symbol.QuoteAsset)
}

func (b *binanceExchangeImpl) GetCandles(symbol types.Symbol, timeFrame types.TimeFrame, limit int) ([]*types.Candle, error) {
	klines, err := b.klinesService.Interval(b.FormatTimeFrame(timeFrame)).Limit(limit).Symbol(b.FormatSymbol(symbol)).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return klinesToCandles(klines), nil
}

func (b *binanceExchangeImpl) GetHistoricalCandles(symbol types.Symbol, timeFrame types.TimeFrame, from time.Time, to time.Time) ([]*types.Candle, error) {
	candles := make([]*types.Candle, 0)

	batchSize := 1000
	timeFromNow := from

	for {
		timeUntil := timeFromNow.Add(time.Duration(batchSize) * types.CandleDurations[timeFrame])

		if timeUntil.After(to) {
			timeUntil = to
		}

		klines, err := b.klinesService.Limit(1000).Interval(b.FormatTimeFrame(timeFrame)).Symbol(b.FormatSymbol(symbol)).StartTime(timeFromNow.UnixNano() / 1000000).EndTime(timeUntil.UnixNano() / 1000000).Do(context.Background())
		if err != nil {
			return nil, err
		}

		candles = append(candles, klinesToCandles(klines)...)

		if timeUntil.Equal(to) {
			break
		}

		timeFromNow = timeUntil
	}

	return candles, nil
}

func (b *binanceExchangeImpl) GetTicker(symbol types.Symbol) (types.Ticker, error) {
	tickers, err := b.bookTickerService.Symbol(b.FormatSymbol(symbol)).Do(context.Background())
	if err != nil {
		return types.Ticker{}, err
	}
	if len(tickers) == 0 {
		return types.Ticker{}, errors.New("ticker not found, list empty")
	}
	ticker := tickers[0]

	return binanceTickerToTicker(ticker), nil
}

// Real-time
func (b *binanceExchangeImpl) WatchTrades(symbols []types.Symbol, handleTrade func(trade types.Trade), handleClose func(error)) (func(), error) {
	symbolMapping := make(map[string]types.Symbol)

	allSymbols := make([]string, 0)
	for _, symbol := range symbols {
		binanceSymbol := b.FormatSymbol(symbol)
		allSymbols = append(allSymbols, binanceSymbol)

		symbolMapping[binanceSymbol] = symbol
	}

	_, stopCh, err := binance.WsCombinedAggTradeServe(allSymbols, func(event *binance.WsAggTradeEvent) {
		handleTrade(types.Trade{
			Symbol:   symbolMapping[event.Symbol],
			Time:     time.Unix(0, event.TradeTime*int64(time.Millisecond)),
			TradeId:  fmt.Sprintf("%d", event.AggTradeID),
			Price:    parseFloatUnsafe(event.Price),
			Quantity: parseFloatUnsafe(event.Quantity),
		})
	}, handleClose)
	if err != nil {
		return nil, err
	}

	return func() {
		stopCh <- struct{}{}
	}, nil
}

// Helper functions
func binanceTickerToTicker(ticker *binance.BookTicker) types.Ticker {
	return types.NewTicker(
		parseFloatUnsafe(ticker.BidPrice),
		parseFloatUnsafe(ticker.AskPrice),
		parseFloatUnsafe(ticker.BidQuantity),
		parseFloatUnsafe(ticker.AskQuantity),
	)
}

func klineToCandle(kline *binance.Kline) *types.Candle {
	return types.NewCandle(
		time.Unix(0, kline.OpenTime*int64(time.Millisecond)),
		time.Unix(0, kline.CloseTime*int64(time.Millisecond)),
		parseFloatUnsafe(kline.Open),
		parseFloatUnsafe(kline.High),
		parseFloatUnsafe(kline.Low),
		parseFloatUnsafe(kline.Close),
		parseFloatUnsafe(kline.Volume),
	)
}
func klinesToCandles(klines []*binance.Kline) []*types.Candle {
	candles := make([]*types.Candle, 0)

	for _, kline := range klines {
		candles = append(candles, klineToCandle(kline))
	}

	return candles
}

func parseFloatUnsafe(entry string) float64 {
	res, _ := strconv.ParseFloat(entry, 32)

	return res
}
