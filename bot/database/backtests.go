package database

import (
	"context"
	"errors"

	"github.com/Stratomicl/Trader/strategy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *DatabaseHandler) GetBacktestById(id string) (*strategy.Backtest, error) {
	collection := d.database.Collection("backtest")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := collection.FindOne(context.Background(), bson.M{
		"_id": objId,
	})

	if result == nil {
		return nil, errors.New("backtest not found")
	}

	backtest := &strategy.Backtest{}
	err = result.Decode(backtest)

	return backtest, err
}
