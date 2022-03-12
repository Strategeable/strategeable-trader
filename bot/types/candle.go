package types

import (
	"errors"
	"fmt"
	"time"
)

type TimeFrame string

const (
	M1  TimeFrame = "1m"
	M3  TimeFrame = "3m"
	M5  TimeFrame = "5m"
	M15 TimeFrame = "15m"
	M30 TimeFrame = "30m"
	H1  TimeFrame = "1h"
	H2  TimeFrame = "2h"
	H4  TimeFrame = "4h"
	D1  TimeFrame = "1d"
	W1  TimeFrame = "1w"
)

var (
	CandleDurations = map[TimeFrame]time.Duration{
		M1:  time.Minute,
		M3:  3 * time.Minute,
		M5:  5 * time.Minute,
		M15: 15 * time.Minute,
		M30: 30 * time.Minute,
		H1:  time.Hour,
		H2:  2 * time.Hour,
		H4:  4 * time.Hour,
		D1:  24 * time.Hour,
		W1:  7 * 24 * time.Hour,
	}
)

type Candle struct {
	Symbol    string    `bson:"s"`
	Exchange  Exchange  `bson:"e"`
	OpenTime  time.Time `bson:"oT"`
	CloseTime time.Time `bson:"cT"`
	Open      float64   `bson:"o"`
	High      float64   `bson:"h"`
	Low       float64   `bson:"l"`
	Close     float64   `bson:"c"`
	Volume    float64   `bson:"v"`
}

func NewCandle(exchange Exchange, symbol string, openTime time.Time, closeTime time.Time, open float64, high float64, low float64, close float64, volume float64) *Candle {
	return &Candle{
		Symbol:    symbol,
		Exchange:  exchange,
		OpenTime:  openTime,
		CloseTime: closeTime,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    volume,
	}
}

func (c Candle) String() string {
	return fmt.Sprintf("O: %.4f, H: %.4f, L: %.4f, C: %.4f, V: %.4f", c.Open, c.High, c.Low, c.Close, c.Volume)
}

func (c *Candle) addTrade(rate float64, volume float64, time time.Time) error {
	if time.After(c.CloseTime) {
		return errors.New("candle already closed")
	}

	c.Close = rate
	c.Volume += volume

	if rate > c.High {
		c.High = rate
	} else if rate < c.Low {
		c.Low = rate
	}

	return nil
}
