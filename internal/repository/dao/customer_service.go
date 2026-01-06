package dao

import (
	"classroom-analysis/internal/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomerServiceDAO struct {
	collection *mongo.Collection
}

func NewCustomerServiceDAO(db *mongo.Database) *CustomerServiceDAO {
	return &CustomerServiceDAO{
		collection: db.Collection("customer_services"),
	}
}

// Create 创建客户服务
func (dao *CustomerServiceDAO) Create(ctx context.Context, cs *domain.CustomerService) error {
	cs.CreatedAt = time.Now()
	cs.UpdatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, cs)
	if err != nil {
		return err
	}
	cs.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindById 根据ID查找客户服务
func (dao *CustomerServiceDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.CustomerService, error) {
	var cs domain.CustomerService
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&cs)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &cs, nil
}

// FindByCustomerID 根据客户ID查找服务列表
func (dao *CustomerServiceDAO) FindByCustomerID(ctx context.Context, customerID primitive.ObjectID) ([]*domain.CustomerService, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{"customer_id": customerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*domain.CustomerService
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// FindList 查询客户服务列表
func (dao *CustomerServiceDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CustomerService, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if skip > 0 {
		opts.SetSkip(skip)
	}

	cursor, err := dao.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*domain.CustomerService
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update 更新客户服务信息
func (dao *CustomerServiceDAO) Update(ctx context.Context, cs *domain.CustomerService) error {
	cs.UpdatedAt = time.Now()
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": cs.ID},
		bson.M{"$set": cs},
	)
	return err
}
