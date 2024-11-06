package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// UserID             primitive.ObjectID `bson:"userId" json:"userId" validate:"required"`
	RestaurantID       primitive.ObjectID `bson:"restaurantId" json:"restaurantId" validate:"required"`
	OrderID            primitive.ObjectID `bson:"orderId" json:"orderId" validate:"required"`
	Rating             int                `bson:"rating" json:"rating" validate:"required,min=1,max=5"`
	Comment            string             `bson:"comment" json:"comment" validate:"required,min=10,max=1000"`
	RestaurantResponse struct {
		Comment   string    `bson:"comment" json:"comment"`
		CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	} `bson:"restaurantResponse" json:"restaurantResponse"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
