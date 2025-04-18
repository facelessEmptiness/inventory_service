package repository

import (
	"context"
	"errors"
	"github.com/facelessEmptiness/AP2Assignment1/internal/domain/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

// ProductRepositoryImpl реализует интерфейс ProductRepository
type ProductRepositoryImpl struct {
	collection *mongo.Collection
}

// NewProductRepository создает новый экземпляр ProductRepositoryImpl
func NewProductRepository(db *mongo.Database) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{
		collection: db.Collection("products"),
	}
}

// Create создает новый продукт в базе данных
func (r *ProductRepositoryImpl) Create(ctx context.Context, product *entity.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	return err
}

// GetByID получает продукт по ID
func (r *ProductRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	var product entity.Product
	filter := bson.M{"id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

// Update обновляет информацию о продукте
func (r *ProductRepositoryImpl) Update(ctx context.Context, product *entity.Product) error {
	product.UpdatedAt = time.Now()

	filter := bson.M{"id": product.ID}
	update := bson.M{"$set": product}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("product not found")
	}

	return nil
}

// Delete удаляет продукт по ID
func (r *ProductRepositoryImpl) Delete(ctx context.Context, id string) error {
	filter := bson.M{"id": id}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("product not found")
	}

	return nil
}

// List возвращает список продуктов с поддержкой пагинации
func (r *ProductRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*entity.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// GetByCategory возвращает продукты по ID категории
func (r *ProductRepositoryImpl) GetByCategory(ctx context.Context, categoryID string, limit, offset int) ([]*entity.Product, error) {
	filter := bson.M{"category_id": categoryID}

	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*entity.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// UpdateStock обновляет количество товара на складе
func (r *ProductRepositoryImpl) UpdateStock(ctx context.Context, id string, count int) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{
		"stock_count": count,
		"updated_at":  time.Now(),
	}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("product not found")
	}

	return nil
}
