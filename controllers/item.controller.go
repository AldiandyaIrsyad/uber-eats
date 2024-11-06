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

type ItemController struct {
	itemService *services.ItemService
}

func NewItemController(itemService *services.ItemService) *ItemController {
	return &ItemController{
		itemService: itemService,
	}
}

func (c *ItemController) CreateItem(ctx *gin.Context) {
	var item models.Item
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.itemService.CreateItem(context.Background(), &item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

func (c *ItemController) GetItemByID(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	item, err := c.itemService.GetItemByID(context.Background(), itemID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *ItemController) UpdateItem(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var update map[string]interface{}
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.itemService.UpdateItem(context.Background(), itemID, update); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

func (c *ItemController) DeleteItem(ctx *gin.Context) {
	id := ctx.Param("id")
	itemID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.itemService.DeleteItem(context.Background(), itemID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

func (c *ItemController) GetItems(ctx *gin.Context) {
	// Parse query parameters
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("pageSize", "10")
	sortField := ctx.DefaultQuery("sortField", "")
	sortOrder := ctx.DefaultQuery("sortOrder", "asc")
	restaurantID := ctx.Query("restaurantID")

	// Convert string parameters to appropriate types
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	pageSizeInt, _ := strconv.ParseInt(pageSize, 10, 64)

	// Build query options
	queryOpts := models.QueryOptions{
		Pagination: &models.PaginationOptions{
			Page:     pageInt,
			PageSize: pageSizeInt,
		},
	}

	// Add sorting if specified
	if sortField != "" {
		order := 1
		if sortOrder == "desc" {
			order = -1
		}
		queryOpts.Sort = &models.SortOptions{
			Field: sortField,
			Order: order,
		}
	}

	// Add filters
	if restaurantID != "" {
		restaurantObjID, err := primitive.ObjectIDFromHex(restaurantID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
			return
		}
		queryOpts.Filter = map[string]interface{}{
			"restaurantId": restaurantObjID,
		}
	}

	// Get items
	result, err := c.itemService.GetItems(context.Background(), queryOpts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
