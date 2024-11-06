package repos

import (
	"context"

	"github.com/aldiandyaIrsyad/uber-eats/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RestaurantRepository struct {
	collection *mongo.Collection
}

func NewRestaurantRepository(client *mongo.Client) *RestaurantRepository {
	collection := client.Database("testing").Collection("restaurants")
	return &RestaurantRepository{collection: collection}
}

func (r *RestaurantRepository) CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) error {
	_, err := r.collection.InsertOne(ctx, restaurant)
	return err
}

func (r *RestaurantRepository) GetRestaurantByID(ctx context.Context, id primitive.ObjectID) (*models.Restaurant, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "items"},
		{Key: "localField", Value: "_id"},
		{Key: "foreignField", Value: "restaurantId"},
		{Key: "as", Value: "items"},
	}}}
	limitStage := bson.D{{Key: "$limit", Value: 10}}

	pipeline := mongo.Pipeline{matchStage, lookupStage, limitStage}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var restaurants []models.Restaurant
	if err = cursor.All(ctx, &restaurants); err != nil {
		return nil, err
	}

	if len(restaurants) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &restaurants[0], nil
}

func (r *RestaurantRepository) UpdateRestaurant(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *RestaurantRepository) DeleteRestaurant(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *RestaurantRepository) GetAllRestaurants(ctx context.Context, filter bson.M, pagination *models.Pagination) ([]models.Restaurant, error) {
	findOptions := options.Find()
	findOptions.SetSkip(pagination.GetSkip())
	findOptions.SetLimit(pagination.GetLimit())

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var restaurants []models.Restaurant
	for cursor.Next(ctx) {
		var restaurant models.Restaurant
		if err := cursor.Decode(&restaurant); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}
