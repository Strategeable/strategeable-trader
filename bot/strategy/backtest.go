package strategy

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BacktestPositionValue struct {
	Date      time.Time `bson:"date"`
	Rate      float64   `bson:"rate"`
	BaseSize  float64   `bson:"baseSize"`
	QuoteFees float64   `bson:"quoteFees"`
}

type BacktestPosition struct {
	OpenedAt   time.Time             `bson:"openedAt"`
	ClosedAt   time.Time             `bson:"closedAt"`
	Symbol     string                `bson:"symbol"`
	EntryValue BacktestPositionValue `bson:"entryValue"`
	ExitValue  BacktestPositionValue `bson:"exitValue"`
}

type Backtest struct {
	Id           primitive.ObjectID `bson:"_id"`
	StrategyId   primitive.ObjectID `bson:"strategyId"`
	Strategy     rawStrategy        `bson:"strategy"`
	FromDate     primitive.DateTime `bson:"fromDate"`
	ToDate       primitive.DateTime `bson:"toDate"`
	StartBalance float64            `bson:"startBalance"`
	EndBalance   float64            `bson:"endBalance"`
	Finished     bool               `bson:"finished"`
	Positions    []BacktestPosition `bson:"positions"`
}
