package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseHandler struct {
	client   *mongo.Client
	database *mongo.Database
}

func (d *DatabaseHandler) Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DATABASE_URL")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	d.client = client
	d.database = client.Database(os.Getenv("DATABASE_NAME"))

	return err
}
