package controllers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aldiandyaIrsyad/uber-eats/models"
	"github.com/aldiandyaIrsyad/uber-eats/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestaurantController struct {
	restaurantService *services.RestaurantService
}

func NewRestaurantController(restaurantService *services.RestaurantService) *RestaurantController {
	return &RestaurantController{
		restaurantService: restaurantService,
	}
}

func (c *RestaurantController) CreateRestaurant(ctx *gin.Context) {
	var restaurant models.Restaurant
	if err := ctx.ShouldBindJSON(&restaurant); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.restaurantService.CreateRestaurant(context.Background(), &restaurant); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, restaurant)
}

func (c *RestaurantController) GetRestaurantByID(ctx *gin.Context) {
	id := ctx.Param("id")
	restaurantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	restaurant, err := c.restaurantService.GetRestaurantByID(context.Background(), restaurantID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	ctx.JSON(http.StatusOK, restaurant)
}

func (c *RestaurantController) UpdateRestaurant(ctx *gin.Context) {
	id := ctx.Param("id")
	restaurantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var update map[string]interface{}
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.restaurantService.UpdateRestaurant(context.Background(), restaurantID, update); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Restaurant updated successfully"})
}

func (c *RestaurantController) DeleteRestaurant(ctx *gin.Context) {
	id := ctx.Param("id")
	restaurantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.restaurantService.DeleteRestaurant(context.Background(), restaurantID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Restaurant deleted successfully"})
}

func (c *RestaurantController) GetAverageRating(ctx *gin.Context) {
	id := ctx.Param("id")
	restaurantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	averageRating, err := c.restaurantService.GetAverageRating(context.Background(), restaurantID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"averageRating": averageRating})
}

func (c *RestaurantController) GetRestaurants(ctx *gin.Context) {
	var filter map[string]interface{}
	if err := ctx.ShouldBindJSON(&filter); err != nil {
		filter = nil // If JSON is empty or invalid, use no filter
	}

	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(ctx.DefaultQuery("pageSize", "10"), 10, 64)

	pagination := &models.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	restaurants, err := c.restaurantService.GetRestaurants(context.Background(), filter, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, restaurants)
}
