package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RestaurantID primitive.ObjectID `bson:"restaurantId" json:"restaurantId" validate:"required"`
	Name         string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	Description  string             `bson:"description" json:"description" validate:"required,min=10,max=500"`
	Price        float64            `bson:"price" json:"price" validate:"required,gt=0"`
	ImageURL     string             `bson:"imageUrl" json:"imageUrl" validate:"required,url"`
	Status       string             `bson:"status" json:"status" validate:"required,oneof=available unavailable"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
}
