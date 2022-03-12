package database

import (
	"context"
	"errors"

	"github.com/Stratomicl/Trader/types"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *DatabaseHandler) GetCandles(exchange types.Exchange, symbol types.Symbol) ([]*types.Candle, error) {
	collection := d.database.Collection("candles")

	cur, err := collection.Find(context.Background(), bson.M{
		"e": exchange,
		"s": symbol.String(),
	})
	if err != nil {
		return nil, err
	}

	var candles []*types.Candle

	err = cur.All(context.Background(), &candles)
	if err != nil {
		return nil, err
	}

	return candles, nil
}

func (d *DatabaseHandler) SaveCandles(candles []*types.Candle) error {
	collection := d.database.Collection("candles")

	docs := make([]interface{}, 0)

	for _, candle := range candles {
		docs = append(docs, candle)
	}

	result, err := collection.InsertMany(context.Background(), docs)
	if err != nil {
		return err
	}

	if len(result.InsertedIDs) != len(docs) {
		return errors.New("did not insert enough new candles")
	}

	return nil
}
