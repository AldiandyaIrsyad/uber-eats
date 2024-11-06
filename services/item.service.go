package services

import (
	"context"

	"github.com/aldiandyaIrsyad/uber-eats/models"
	"github.com/aldiandyaIrsyad/uber-eats/repos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemService struct {
	itemRepo *repos.ItemRepository
}

func NewItemService(itemRepo *repos.ItemRepository) *ItemService {
	return &ItemService{
		itemRepo: itemRepo,
	}
}

func (s *ItemService) CreateItem(ctx context.Context, item *models.Item) error {
	return s.itemRepo.CreateItem(ctx, item)
}

func (s *ItemService) GetItemByID(ctx context.Context, id primitive.ObjectID) (*models.Item, error) {
	return s.itemRepo.GetItemByID(ctx, id)
}

func (s *ItemService) UpdateItem(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error {
	updateBson := bson.M{"$set": update}
	return s.itemRepo.UpdateItem(ctx, id, updateBson)
}

func (s *ItemService) DeleteItem(ctx context.Context, id primitive.ObjectID) error {
	return s.itemRepo.DeleteItem(ctx, id)
}

func (s *ItemService) GetItems(ctx context.Context, queryOpts models.QueryOptions) (*models.Pagination, error) {
	// Set default values if not provided
	if queryOpts.Pagination == nil {
		queryOpts.Pagination = &models.PaginationOptions{
			Page:     1,
			PageSize: 10,
		}
	}

	return s.itemRepo.FindWithOptions(ctx, queryOpts)
}
