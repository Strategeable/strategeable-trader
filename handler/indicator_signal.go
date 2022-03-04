package handler

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
	Symbol      string
}

type IndicatorSignal struct {
	IndicatorA  IndicatorSettings
	IndicatorB  IndicatorSettings
	Operand     Operand
	Persistence int
}

func (i *IndicatorSignal) HasSignal(candles []*types.Candle) (bool, error) {
	candlesA := getCandles(candles, i.IndicatorA)
	candlesB := getCandles(candles, i.IndicatorB)

	valuesA := i.IndicatorA.Indicator.Calculate(candlesA)
	valuesB := i.IndicatorB.Indicator.Calculate(candlesB)

	if len(valuesA) < i.Persistence || len(valuesB) < i.Persistence {
		return false, errors.New("not enough candle values to cover persistence")
	}

	for j := 0; j < i.Persistence; j++ {
		valueA := valuesA[len(valuesA)-i.Persistence-1+j]
		valueB := valuesB[len(valuesB)-i.Persistence-1+j]

		isFinalIndex := j == i.Persistence-1

		switch i.Operand {
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

func getCandles(candles []*types.Candle, settings IndicatorSettings) []*types.Candle {
	finalCandles := candles[:len(candles)-settings.CandlesBack]
	if !settings.RealTime {
		finalCandles = finalCandles[:len(finalCandles)-1]
	}
	return finalCandles
}
