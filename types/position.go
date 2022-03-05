package types

import "time"

type Position struct {
	OpenTime  time.Time
	CloseTime time.Time
	Symbol    Symbol
}
