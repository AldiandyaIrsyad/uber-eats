package repos

import (
	"context"
	"time"

	"github.com/aldiandyaIrsyad/uber-eats/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewRepository struct {
	collection *mongo.Collection
}

func NewReviewRepository(client *mongo.Client) *ReviewRepository {
	collection := client.Database("testing").Collection("reviews")
	return &ReviewRepository{collection: collection}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, review *models.Review) error {
	review.ID = primitive.NewObjectID()
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, review)
	return err
}

func (r *ReviewRepository) GetReviewsByRestaurantID(ctx context.Context, restaurantID primitive.ObjectID) ([]models.Review, error) {
	filter := bson.M{"restaurantId": restaurantID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	for cursor.Next(ctx) {
		var review models.Review
		if err := cursor.Decode(&review); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r *ReviewRepository) GetAverageRatingByRestaurantID(ctx context.Context, restaurantID primitive.ObjectID) (float64, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "restaurantId", Value: restaurantID}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$restaurantId"},
		{Key: "averageRating", Value: bson.D{{Key: "$avg", Value: "$rating"}}},
	}}}

	cursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		AverageRating float64 `bson:"averageRating"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	return result.AverageRating, nil
}
