package types

// Indicators are used in Signal and Any Signal tiles.
// They perform math on a slice of candles and return
// values that can be used to base decisions on.
type Indicator interface {
	// Calculate all indicator values based on a slice of candles
	// and a (optional) position object
	Calculate(input []*Candle, position *Position) []float64
}
