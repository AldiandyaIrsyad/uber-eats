package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() (*mongo.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	clientOptions := options.Client().
		ApplyURI("mongodb://mongodb:27017").
		SetMaxPoolSize(100).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(5 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Create indexes
	createIndexes(client)

	return client, ctx, cancel
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("testing").Collection(collectionName)
	return collection
}

func createIndexes(client *mongo.Client) {
	reviewCollection := GetCollection(client, "reviews")
	reviewIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "restaurantId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "orderId", Value: 1}},
		},
	}
	_, err := reviewCollection.Indexes().CreateMany(context.Background(), reviewIndexes)
	if err != nil {
		log.Fatal(err)
	}

	restaurantCollection := GetCollection(client, "restaurants")
	restaurantIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "location", Value: "2dsphere"}},
		},
	}
	_, err = restaurantCollection.Indexes().CreateMany(context.Background(), restaurantIndexes)
	if err != nil {
		log.Fatal(err)
	}

	itemCollection := GetCollection(client, "items")
	itemIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "restaurantId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
	}
	_, err = itemCollection.Indexes().CreateMany(context.Background(), itemIndexes)
	if err != nil {
		log.Fatal(err)
	}
}
