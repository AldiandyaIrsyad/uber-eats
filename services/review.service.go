package services

import (
	"context"

	"github.com/aldiandyaIrsyad/uber-eats/models"
	"github.com/aldiandyaIrsyad/uber-eats/repos"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewService struct {
	reviewRepo *repos.ReviewRepository
}

func NewReviewService(reviewRepo *repos.ReviewRepository) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, review *models.Review) error {
	return s.reviewRepo.CreateReview(ctx, review)
}

func (s *ReviewService) GetReviewsByRestaurantID(ctx context.Context, restaurantID primitive.ObjectID) ([]models.Review, error) {
	return s.reviewRepo.GetReviewsByRestaurantID(ctx, restaurantID)
}

func (s *ReviewService) GetAverageRatingByRestaurantID(ctx context.Context, restaurantID primitive.ObjectID) (float64, error) {
	return s.reviewRepo.GetAverageRatingByRestaurantID(ctx, restaurantID)
}
