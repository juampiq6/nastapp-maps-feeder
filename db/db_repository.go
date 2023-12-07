package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertManyIntoCollection(colName string, items []interface{}) {
	client := GetMongoClientInstance()
	opts := options.InsertMany().SetOrdered(false)
	_, err := client.Database(dbName).Collection(colName).InsertMany(context.TODO(), items, opts)
	if err != nil {
		panic(err)
	}
}
