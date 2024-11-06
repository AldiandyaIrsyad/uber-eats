package services

import (
	"context"

	"github.com/aldiandyaIrsyad/uber-eats/models"
	"github.com/aldiandyaIrsyad/uber-eats/repos"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestaurantService struct {
	restaurantRepo *repos.RestaurantRepository
	reviewRepo     *repos.ReviewRepository
}

func NewRestaurantService(restaurantRepo *repos.RestaurantRepository, reviewRepo *repos.ReviewRepository) *RestaurantService {
	return &RestaurantService{
		restaurantRepo: restaurantRepo,
		reviewRepo:     reviewRepo,
	}
}

func (s *RestaurantService) CreateRestaurant(ctx context.Context, restaurant *models.Restaurant) error {
	return s.restaurantRepo.CreateRestaurant(ctx, restaurant)
}

func (s *RestaurantService) GetRestaurantByID(ctx context.Context, id primitive.ObjectID) (*models.Restaurant, error) {
	restaurant, err := s.restaurantRepo.GetRestaurantByID(ctx, id)
	if err != nil {
		return nil, err
	}

	averageRating, err := s.reviewRepo.GetAverageRatingByRestaurantID(ctx, id)
	if err != nil {
		return nil, err
	}
	restaurant.AverageRating = averageRating

	return restaurant, nil
}

func (s *RestaurantService) UpdateRestaurant(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error {
	updateBson := map[string]interface{}{
		"$set": update,
	}
	return s.restaurantRepo.UpdateRestaurant(ctx, id, updateBson)
}

func (s *RestaurantService) DeleteRestaurant(ctx context.Context, id primitive.ObjectID) error {
	return s.restaurantRepo.DeleteRestaurant(ctx, id)
}

func (s *RestaurantService) GetAverageRating(ctx context.Context, restaurantID primitive.ObjectID) (float64, error) {
	return s.reviewRepo.GetAverageRatingByRestaurantID(ctx, restaurantID)
}

func (s *RestaurantService) GetRestaurants(ctx context.Context, filter map[string]interface{}, pagination *models.Pagination) ([]models.Restaurant, error) {
	if filter == nil {
		filter = make(map[string]interface{})
	}

	restaurants, err := s.restaurantRepo.GetAllRestaurants(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	for i := range restaurants {
		averageRating, err := s.reviewRepo.GetAverageRatingByRestaurantID(ctx, restaurants[i].ID)
		if err != nil {
			return nil, err
		}
		restaurants[i].AverageRating = averageRating
	}

	return restaurants, nil
}
