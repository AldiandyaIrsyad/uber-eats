# Uber-eats

A demonstration of implementation mongodb in Go language

## How to run

1. Clone the repository
2. Run the command `docker-compose up --build` in the root directory of the project

## Requirement

### Showcase Aggregation

To showcase aggregation, we're normalizing our rating. We're storing the rating in the range of 1-5. We're calculating the average rating of a restaurant using the aggregation pipeline.

```Go
// review.repo.go
func (r *ReviewRepository) GetAverageRatingByRestaurantID(ctx context.Context, restaurantID primitive.ObjectID) (float64, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "restaurantId", Value: restaurantID}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: "$restaurantId"},
		{Key: "averageRating", Value: bson.D{{Key: "$avg", Value: "$rating"}}},
	}}}

	cursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		AverageRating float64 `bson:"averageRating"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	return result.AverageRating, nil
}

```

This is used to get the average rating of a restaurant

```Go
// restaurant.service.go
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

```

```Bash
curl --location 'http://localhost:8080/api/restaurants/'
```

I wouldn't use this as normal implementation as using loops in database are really slow and considered bad practices. Rather add `rating` and `ratingCount` in the restaurant collection and update them whenever a new review is added or updated.

### Showcase sorting & Limit

to showcase sorting and limit, we're using a pagination that can be sorted

```Go
// item.repo.go

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


```

```Bash
curl --location 'http://localhost:8080/api/items?page=2&pageSize=5&restaurantID=672bd1e53c51c50425934960&sortField=price&sortOrder=desc'
```

### Join

```Go
// restaurant.repo.go
func (r *RestaurantRepository) GetRestaurantByID(ctx context.Context, id primitive.ObjectID) (*models.Restaurant, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "items"},
		{Key: "localField", Value: "_id"},
		{Key: "foreignField", Value: "restaurantId"},
		{Key: "as", Value: "items"},
	}}}
	limitStage := bson.D{{Key: "$limit", Value: 10}}

	pipeline := mongo.Pipeline{matchStage, lookupStage, limitStage}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var restaurants []models.Restaurant
	if err = cursor.All(ctx, &restaurants); err != nil {
		return nil, err
	}

	if len(restaurants) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &restaurants[0], nil
}
```

```Bash
curl --location 'http://localhost:8080/api/restaurants/672bd1e53c51c50425934960'
```

### CRUD

For demonstration we're going to create, read, update and delete a restaurant

Create new restaurant

```Bash
curl --location 'http://localhost:8080/api/restaurants/' \
--header 'Content-Type: application/json' \
--data '{
           "name": "Burger Palace",
           "description": "Best burgers in town",
           "address": "123 Main St",
           "imageUrl": "https://example.com/burger.jpg",
           "location": {
             "type": "Point",
             "coordinates": [-73.935242, 40.730610]
           },
           "operatingHours": [
             {"day": "Monday", "openTime": "09:00", "closeTime": "22:00"},
             {"day": "Tuesday", "openTime": "09:00", "closeTime": "22:00"},
             {"day": "Wednesday", "openTime": "09:00", "closeTime": "22:00"},
             {"day": "Thursday", "openTime": "09:00", "closeTime": "22:00"},
             {"day": "Friday", "openTime": "09:00", "closeTime": "23:00"},
             {"day": "Saturday", "openTime": "10:00", "closeTime": "23:00"},
             {"day": "Sunday", "openTime": "10:00", "closeTime": "22:00"}
           ]
         }'
```

Read by ID

```Bash
curl --location 'http://localhost:8080/api/restaurants/672be0b125a2a7b9cd92e136'
```

Update

```Bash
curl --location --request PUT 'http://localhost:8080/api/restaurants/672be0b125a2a7b9cd92e136' \
--header 'Content-Type: application/json' \
--data '{
           "name": "Updated Burger Palace",
           "description": "Updated description"
}'
```

Delete

```Bash
curl --location --request DELETE 'http://localhost:8080/api/restaurants/672be0b125a2a7b9cd92e136'
```
