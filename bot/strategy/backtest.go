package strategy

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Backtest struct {
	Id           primitive.ObjectID `bson:"_id"`
	StrategyId   primitive.ObjectID
	Strategy     rawStrategy
	FromDate     primitive.DateTime
	ToDate       primitive.DateTime
	StartBalance float64
}
