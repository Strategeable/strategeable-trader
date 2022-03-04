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
}

type SignalTile struct {
	IndicatorA  IndicatorSettings
	IndicatorB  IndicatorSettings
	Operand     Operand
	Persistence int
}

func (s *SignalTile) HasSignal(candleCollection *types.CandleCollection, symbol types.Symbol) (bool, error) {
	candlesA := getCandles(candleCollection, symbol, s.IndicatorA)
	candlesB := getCandles(candleCollection, symbol, s.IndicatorB)

	valuesA := s.IndicatorA.Indicator.Calculate(candlesA)
	valuesB := s.IndicatorB.Indicator.Calculate(candlesB)

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

func getCandles(candleCollection *types.CandleCollection, symbol types.Symbol, settings IndicatorSettings) []*types.Candle {
	if settings.Symbol != nil {
		symbol = *settings.Symbol
	}

	candles := candleCollection.GetCache(symbol, settings.TimeFrame).GetCandles()

	finalCandles := candles[:len(candles)-settings.CandlesBack]
	if !settings.RealTime {
		finalCandles = finalCandles[:len(finalCandles)-1]
	}
	return finalCandles
}
