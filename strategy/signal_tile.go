package strategy

import (
	"cex-bot/types"
	"errors"
)

type Operand string

const (
	GREATER_THAN          Operand = "GREATER_THAN"
	LOWER_THAN            Operand = "LOWER_THAN"
	GREATER_THAN_OR_EQUAL Operand = "GREATER_THAN_OR_EQUAL"
	LOWER_THAN_OR_EQUAL   Operand = "LOWER_THAN_OR_EQUAL"
	EQUAL                 Operand = "EQUAL"
	CROSS_ABOVE           Operand = "CROSS_ABOVE"
	CROSS_BELOW           Operand = "CROSS_BELOW"
)

type IndicatorSettings struct {
	Indicator   types.Indicator
	RealTime    bool
	CandlesBack int
	TimeFrame   types.TimeFrame
	Symbol      *types.Symbol
	Exchange    *types.Exchange
}

type SignalTile struct {
	IndicatorA  IndicatorSettings
	IndicatorB  IndicatorSettings
	Operand     Operand
	Persistence int
}

func (s *SignalTile) HasSignal(candleCollection *types.CandleCollection, symbol types.Symbol, exchange types.Exchange, position *types.Position) (bool, error) {
	candlesA, err := getCandles(candleCollection, exchange, symbol, s.IndicatorA)
	if err != nil {
		return false, err
	}
	candlesB, err := getCandles(candleCollection, exchange, symbol, s.IndicatorB)
	if err != nil {
		return false, err
	}

	valuesA := s.IndicatorA.Indicator.Calculate(candlesA, position)
	valuesB := s.IndicatorB.Indicator.Calculate(candlesB, position)

	if len(valuesA) < s.Persistence || len(valuesB) < s.Persistence {
		return false, errors.New("not enough candle values to cover persistence")
	}

	for j := 0; j < s.Persistence; j++ {
		valueA := valuesA[len(valuesA)-s.Persistence+j]
		valueB := valuesB[len(valuesB)-s.Persistence+j]

		isFinalIndex := j == s.Persistence-1

		switch s.Operand {
		case GREATER_THAN:
			if valueA <= valueB {
				return false, nil
			}
		case LOWER_THAN:
			if valueA >= valueB {
				return false, nil
			}
		case GREATER_THAN_OR_EQUAL:
			if valueA < valueB {
				return false, nil
			}
		case LOWER_THAN_OR_EQUAL:
			if valueA > valueB {
				return false, nil
			}
		case EQUAL:
			if valueA != valueB {
				return false, nil
			}
		case CROSS_ABOVE:
			if isFinalIndex {
				if valueA <= valueB {
					return false, nil
				}
			} else {
				if valueA > valueB {
					return false, nil
				}
			}
		case CROSS_BELOW:
			if isFinalIndex {
				if valueA >= valueB {
					return false, nil
				}
			} else {
				if valueA < valueB {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

func getCandles(candleCollection *types.CandleCollection, exchange types.Exchange, symbol types.Symbol, settings IndicatorSettings) ([]*types.Candle, error) {
	if settings.Symbol != nil {
		symbol = *settings.Symbol
	}
	if settings.Exchange != nil {
		exchange = *settings.Exchange
	}

	cache := candleCollection.GetCache(exchange, symbol, settings.TimeFrame)
	if cache == nil {
		return nil, errors.New("cache not found")
	}
	candles := cache.GetCandles()

	finalCandles := candles[:len(candles)-settings.CandlesBack]
	if !settings.RealTime {
		finalCandles = finalCandles[:len(finalCandles)-1]
	}
	return finalCandles, nil
}
