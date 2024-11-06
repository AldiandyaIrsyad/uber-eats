package repos

import (
	"context"

	"github.com/aldiandyaIrsyad/uber-eats/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ItemRepository struct {
	collection *mongo.Collection
}

func NewItemRepository(client *mongo.Client) *ItemRepository {
	collection := client.Database("testing").Collection("items")
	return &ItemRepository{collection: collection}
}

func (r *ItemRepository) CreateItem(ctx context.Context, item *models.Item) error {
	_, err := r.collection.InsertOne(ctx, item)
	return err
}

func (r *ItemRepository) GetItemByID(ctx context.Context, id primitive.ObjectID) (*models.Item, error) {
	var item models.Item
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	return &item, err
}

func (r *ItemRepository) UpdateItem(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *ItemRepository) DeleteItem(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ItemRepository) FindWithOptions(ctx context.Context, queryOpts models.QueryOptions) (*models.Pagination, error) {
	filter := bson.M{}

	// Merge custom filters
	if queryOpts.Filter != nil {
		for k, v := range queryOpts.Filter {
			filter[k] = v
		}
	}

	// Set up options
	findOptions := options.Find()

	// Handle pagination
	if queryOpts.Pagination != nil {
		findOptions.SetSkip((queryOpts.Pagination.Page - 1) * queryOpts.Pagination.PageSize)
		findOptions.SetLimit(queryOpts.Pagination.PageSize)
	}

	// Handle sorting
	if queryOpts.Sort != nil {
		findOptions.SetSort(bson.D{{Key: queryOpts.Sort.Field, Value: queryOpts.Sort.Order}})
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []models.Item
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	// Create pagination result
	pagination := models.NewPagination(queryOpts.Pagination.Page, queryOpts.Pagination.PageSize)
	pagination.Total = total
	pagination.Data = make([]interface{}, len(items))
	for i, item := range items {
		pagination.Data[i] = item
	}

	return pagination, nil
}
