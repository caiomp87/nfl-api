package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ConnectionString string
	dbName           string
	dbCollection     string
)

func Connect() (*mongo.Collection, context.Context, error) {
	dbName = os.Getenv("DB_NAME")
	dbCollection = os.Getenv("DB_COLLECTION")

	ConnectionString = fmt.Sprintf("mongodb://%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(ConnectionString))
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}

	mongodb := client.Database(dbName).Collection(dbCollection)
	return mongodb, context.Background(), nil
}
