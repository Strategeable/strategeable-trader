package strategy

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BacktestPositionValue struct {
	Date      time.Time `bson:"date" json:"date"`
	Rate      float64   `bson:"rate" json:"rate"`
	BaseSize  float64   `bson:"baseSize" json:"baseSize"`
	QuoteFees float64   `bson:"quoteFees" json:"quoteFees"`
}

type BacktestPosition struct {
	OpenedAt   time.Time             `bson:"openedAt" json:"openedAt"`
	ClosedAt   time.Time             `bson:"closedAt" json:"closedAt"`
	Symbol     string                `bson:"symbol" json:"symbol"`
	EntryValue BacktestPositionValue `bson:"entryValue" json:"entryValue"`
	ExitValue  BacktestPositionValue `bson:"exitValue" json:"exitValue"`
}

type Backtest struct {
	Id           primitive.ObjectID `bson:"_id"`
	Status       string             `bson:"status"`
	StrategyId   primitive.ObjectID `bson:"strategyId"`
	Strategy     rawStrategy        `bson:"strategy"`
	FromDate     primitive.DateTime `bson:"fromDate"`
	ToDate       primitive.DateTime `bson:"toDate"`
	StartBalance float64            `bson:"startBalance"`
	EndBalance   float64            `bson:"endBalance"`
	Finished     bool               `bson:"finished"`
	Positions    []BacktestPosition `bson:"positions"`
}
