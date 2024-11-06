package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OperatingHours struct {
	Day       string `bson:"day" json:"day"`
	OpenTime  string `bson:"openTime" json:"openTime"`
	CloseTime string `bson:"closeTime" json:"closeTime"`
}

type Restaurant struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// OwnerID     primitive.ObjectID `bson:"ownerId" json:"ownerId" validate:"required"`
	Name        string `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Description string `bson:"description" json:"description"`
	Address     string `bson:"address" json:"address" validate:"required"`
	ImageURL    string `bson:"imageUrl" json:"imageUrl"`
	// Normally I would do something like this
	// Rating      float64            `bson:"rating" json:"rating" validate:"min=0,max=5"`
	// RatingCount int                `bson:"ratingCount" json:"ratingCount"`
	Location struct {
		Type        string    `bson:"type" json:"type"`
		Coordinates []float64 `bson:"coordinates" json:"coordinates"`
	} `bson:"location" json:"location"`
	OperatingHours []OperatingHours `bson:"operatingHours" json:"operatingHours"`
	CreatedAt      time.Time        `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time        `bson:"updatedAt" json:"updatedAt"`

	AverageRating float64 `bson:"averageRating,omitempty" json:"averageRating,omitempty"`
	Items         []Item  `bson:"items,omitempty" json:"items,omitempty"`
}
