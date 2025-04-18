package repository

import (
	"context"
	"github.com/facelessEmptiness/inventory_service/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoProductRepo struct {
	coll *mongo.Collection
}

func NewMongoProductRepository(db *mongo.Database) ProductRepository {
	return &mongoProductRepo{coll: db.Collection("products")}
}

func (r *mongoProductRepo) Create(p *domain.Product) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := bson.M{
		"name":        p.Name,
		"description": p.Description,
		"price":       p.Price,
		"stock":       p.Stock,
		"category_id": p.CategoryID,
	}
	res, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}
	oid := res.InsertedID.(primitive.ObjectID).Hex()
	return oid, nil
}

func (r *mongoProductRepo) GetByID(id string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, _ := primitive.ObjectIDFromHex(id)
	var p domain.Product
	if err := r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}
