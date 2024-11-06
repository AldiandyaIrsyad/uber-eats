package routes

import (
	"github.com/aldiandyaIrsyad/uber-eats/controllers"
	"github.com/aldiandyaIrsyad/uber-eats/repos"
	"github.com/aldiandyaIrsyad/uber-eats/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type RouteHandler struct {
	restaurantController *controllers.RestaurantController
	itemController       *controllers.ItemController
	reviewController     *controllers.ReviewController
}

func NewRouteHandler(client *mongo.Client) *RouteHandler {
	// Initialize repositories
	restaurantRepo := repos.NewRestaurantRepository(client)
	itemRepo := repos.NewItemRepository(client)
	reviewRepo := repos.NewReviewRepository(client)

	// Initialize services
	restaurantService := services.NewRestaurantService(restaurantRepo, reviewRepo)
	itemService := services.NewItemService(itemRepo)
	reviewService := services.NewReviewService(reviewRepo)

	// Initialize controllers
	restaurantController := controllers.NewRestaurantController(restaurantService)
	itemController := controllers.NewItemController(itemService)
	reviewController := controllers.NewReviewController(reviewService)

	return &RouteHandler{
		restaurantController: restaurantController,
		itemController:       itemController,
		reviewController:     reviewController,
	}
}

func (rh *RouteHandler) SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Restaurant routes
		restaurants := api.Group("/restaurants")
		{
			restaurants.POST("", rh.restaurantController.CreateRestaurant)
			restaurants.GET("/:id", rh.restaurantController.GetRestaurantByID)
			restaurants.PUT("/:id", rh.restaurantController.UpdateRestaurant)
			restaurants.DELETE("/:id", rh.restaurantController.DeleteRestaurant)
			restaurants.GET("/:id/rating", rh.restaurantController.GetAverageRating)
			restaurants.GET("", rh.restaurantController.GetRestaurants)
		}

		// Item routes
		items := api.Group("/items")
		{
			items.POST("", rh.itemController.CreateItem)
			items.GET("", rh.itemController.GetItems)
			items.GET("/:id", rh.itemController.GetItemByID)
			items.PUT("/:id", rh.itemController.UpdateItem)
			items.DELETE("/:id", rh.itemController.DeleteItem)
		}

		// Review routes
		reviews := api.Group("/reviews")
		{
			reviews.POST("", rh.reviewController.CreateReview)
			reviews.GET("/restaurant/:restaurantID", rh.reviewController.GetReviewsByRestaurantID)
			reviews.GET("/restaurant/:restaurantID/rating", rh.reviewController.GetAverageRatingByRestaurantID)
		}
	}
}
