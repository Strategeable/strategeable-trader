package database

import (
	"context"
	"errors"

	"github.com/Strategeable/Trader/strategy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *DatabaseHandler) GetBacktestById(id string) (*strategy.Backtest, error) {
	collection := d.database.Collection("backtests")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := collection.FindOne(context.Background(), bson.M{
		"_id": objId,
	})

	if result.Err() != nil {
		return nil, result.Err()
	}

	backtest := &strategy.Backtest{}
	err = result.Decode(backtest)

	return backtest, err
}

func (d *DatabaseHandler) UpdateBacktestStatus(id string, status string) error {
	collection := d.database.Collection("backtests")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateResult, err := collection.UpdateByID(context.Background(), objectId, bson.M{
		"$set": bson.M{
			"status": status,
		},
	})
	if err != nil {
		return err
	}

	if updateResult.MatchedCount < 1 {
		return errors.New("backtest not found")
	}

	return nil
}

func (d *DatabaseHandler) SaveBacktest(backtest *strategy.Backtest) error {
	collection := d.database.Collection("backtests")

	updateResult, err := collection.UpdateByID(context.Background(), backtest.Id, bson.M{
		"$set": bson.M{
			"endBalance": backtest.EndBalance,
			"finished":   backtest.Finished,
			"positions":  backtest.Positions,
		},
	})
	if err != nil {
		return err
	}

	if updateResult.MatchedCount < 1 {
		return errors.New("backtest not found")
	}

	return nil
}
