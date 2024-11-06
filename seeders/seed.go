package seeders

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aldiandyaIrsyad/uber-eats/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	numRestaurants = 2
	numItems       = 15
	numReviews     = 15
)

func createRestaurants() []models.Restaurant {
	restaurants := make([]models.Restaurant, numRestaurants)

	for i := 0; i < numRestaurants; i++ {
		restaurants[i] = models.Restaurant{
			Name:        fmt.Sprintf("Restaurant %d", i+1),
			Description: fmt.Sprintf("Description for Restaurant %d", i+1),
			Address:     fmt.Sprintf("%d Main Street", (i+1)*100),
			ImageURL:    fmt.Sprintf("https://example.com/restaurant-%d.jpg", i+1),
			Location: struct {
				Type        string    `bson:"type" json:"type"`
				Coordinates []float64 `bson:"coordinates" json:"coordinates"`
			}{
				Type:        "Point",
				Coordinates: []float64{-73.935242 + float64(i)*0.01, 40.730610}, // Slight variation in coordinates
			},
			OperatingHours: []models.OperatingHours{
				{Day: "Monday", OpenTime: "09:00", CloseTime: "22:00"},
				{Day: "Tuesday", OpenTime: "09:00", CloseTime: "22:00"},
				{Day: "Wednesday", OpenTime: "09:00", CloseTime: "22:00"},
				{Day: "Thursday", OpenTime: "09:00", CloseTime: "22:00"},
				{Day: "Friday", OpenTime: "09:00", CloseTime: "23:00"},
				{Day: "Saturday", OpenTime: "10:00", CloseTime: "23:00"},
				{Day: "Sunday", OpenTime: "10:00", CloseTime: "22:00"},
			},
		}
	}
	return restaurants
}

func createItems(restaurantID primitive.ObjectID, restaurantIndex int) []models.Item {
	items := make([]models.Item, numItems)

	for i := 0; i < numItems; i++ {
		items[i] = models.Item{
			RestaurantID: restaurantID,
			Name:         fmt.Sprintf("Item %d (Restaurant %d)", i+1, restaurantIndex+1),
			Description:  fmt.Sprintf("Description for Item %d from Restaurant %d", i+1, restaurantIndex+1),
			Price:        float64(5+(i*2)) + 0.99, // Prices from 5.99 to 23.99
			ImageURL:     fmt.Sprintf("https://example.com/restaurant-%d/item-%d.jpg", restaurantIndex+1, i+1),
			Status:       "available",
		}
	}
	return items
}

func createReviews(restaurantID, orderID primitive.ObjectID, restaurantIndex int) []models.Review {
	reviews := make([]models.Review, numReviews)

	for i := 0; i < numReviews; i++ {
		rating := 3 + (i % 3) // Ratings from 3 to 5
		reviews[i] = models.Review{
			RestaurantID: restaurantID,
			OrderID:      orderID,
			Rating:       rating,
			Comment:      fmt.Sprintf("Review %d for Restaurant %d - %d stars", i+1, restaurantIndex+1, rating),
			RestaurantResponse: struct {
				Comment   string    `bson:"comment" json:"comment"`
				CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
			}{
				Comment:   fmt.Sprintf("Thank you for your %d-star review!", rating),
				CreatedAt: time.Now(),
			},
		}
	}
	return reviews
}

func isDatabaseEmpty(client *mongo.Client) (bool, error) {
	ctx := context.Background()
	db := client.Database("testing")

	// Check each collection
	collections := []string{"restaurants", "items", "reviews"}
	for _, collName := range collections {
		count, err := db.Collection(collName).CountDocuments(ctx, bson.M{})
		if err != nil {
			return false, err
		}
		if count > 0 {
			return false, nil
		}
	}
	return true, nil
}

func SeedDatabase(client *mongo.Client) error {
	isEmpty, err := isDatabaseEmpty(client)
	if err != nil {
		return fmt.Errorf("error checking database: %v", err)
	}

	if !isEmpty {
		log.Println("Database is not empty, skipping seed")
		return nil
	}
	ctx := context.Background()

	// Create restaurants
	restaurants := createRestaurants()
	restaurantCollection := client.Database("testing").Collection("restaurants")
	itemCollection := client.Database("testing").Collection("items")
	reviewCollection := client.Database("testing").Collection("reviews")

	for i, restaurant := range restaurants {
		// Seed restaurant
		restaurant.ID = primitive.NewObjectID()
		restaurant.CreatedAt = time.Now()
		restaurant.UpdatedAt = time.Now()
		_, err := restaurantCollection.InsertOne(ctx, restaurant)
		if err != nil {
			return fmt.Errorf("error seeding restaurant %d: %v", i+1, err)
		}

		// Seed items
		items := createItems(restaurant.ID, i)
		for _, item := range items {
			item.ID = primitive.NewObjectID()
			item.CreatedAt = time.Now()
			item.UpdatedAt = time.Now()
			_, err := itemCollection.InsertOne(ctx, item)
			if err != nil {
				return fmt.Errorf("error seeding item for restaurant %d: %v", i+1, err)
			}
		}

		// Seed reviews
		orderID := primitive.NewObjectID()
		reviews := createReviews(restaurant.ID, orderID, i)
		for _, review := range reviews {
			review.ID = primitive.NewObjectID()
			review.CreatedAt = time.Now()
			review.UpdatedAt = time.Now()
			_, err := reviewCollection.InsertOne(ctx, review)
			if err != nil {
				return fmt.Errorf("error seeding review for restaurant %d: %v", i+1, err)
			}
		}
	}

	log.Printf("Seeded %d restaurants with %d items and %d reviews each", numRestaurants, numItems, numReviews)
	return nil
}
