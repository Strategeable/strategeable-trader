package types

type Indicator interface {
	Calculate(input []*Candle, position *Position) []float64
}
