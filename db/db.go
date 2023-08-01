package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB() (*mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MongoURI"))

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	db := client.Database("player_score_management")
	collection := db.Collection("players")

	fmt.Println("db running")

	return collection, nil
}
