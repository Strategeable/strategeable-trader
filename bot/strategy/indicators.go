package strategy

import (
	"reflect"

	"github.com/Stratomicl/Trader/indicators"
)

var (
	INDICATOR_REGISTRY = map[string]reflect.Type{
		"RSI":                   reflect.TypeOf(indicators.RsiIndicator{}),
		"SMA":                   reflect.TypeOf(indicators.SmaIndicator{}),
		"EMA":                   reflect.TypeOf(indicators.EmaIndicator{}),
		"NUMBER":                reflect.TypeOf(indicators.NumberIndicator{}),
		"CANDLE_POSITION_VALUE": reflect.TypeOf(indicators.CandlePositionValueIndicator{}),
		"BOLLINGER_BAND":        reflect.TypeOf(indicators.BollingerBandIndicator{}),
	}
)
