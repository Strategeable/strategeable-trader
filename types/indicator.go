package types

type Indicator interface {
	Calculate(input []*Candle) []float64
}
