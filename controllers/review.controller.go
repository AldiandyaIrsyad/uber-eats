package controllers

import (
	"context"
	"net/http"

	"github.com/aldiandyaIrsyad/uber-eats/models"
	"github.com/aldiandyaIrsyad/uber-eats/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewController struct {
	reviewService *services.ReviewService
}

func NewReviewController(reviewService *services.ReviewService) *ReviewController {
	return &ReviewController{
		reviewService: reviewService,
	}
}

func (c *ReviewController) CreateReview(ctx *gin.Context) {
	var review models.Review
	if err := ctx.ShouldBindJSON(&review); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.reviewService.CreateReview(context.Background(), &review); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, review)
}

func (c *ReviewController) GetReviewsByRestaurantID(ctx *gin.Context) {
	id := ctx.Param("restaurantID")
	restaurantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	reviews, err := c.reviewService.GetReviewsByRestaurantID(context.Background(), restaurantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}

func (c *ReviewController) GetAverageRatingByRestaurantID(ctx *gin.Context) {
	id := ctx.Param("restaurantID")
	restaurantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	averageRating, err := c.reviewService.GetAverageRatingByRestaurantID(context.Background(), restaurantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"averageRating": averageRating})
}
