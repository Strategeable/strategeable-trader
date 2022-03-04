package types

import (
	"errors"
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
)

var (
	CandleDurations = map[TimeFrame]time.Duration{
		M1:  time.Minute,
		M3:  3 * time.Minute,
		M5:  5 * time.Minute,
		M15: 15 * time.Minute,
		M30: 30 * time.Minute,
		H1:  time.Hour,
	}
)

type Candle struct {
	OpenTime  time.Time
	CloseTime time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

func NewCandle(openTime time.Time, closeTime time.Time, open float64, high float64, low float64, close float64, volume float64) *Candle {
	return &Candle{
		OpenTime:  openTime,
		CloseTime: closeTime,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    volume,
	}
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
