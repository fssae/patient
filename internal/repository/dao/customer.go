package dao

import (
	"classroom-analysis/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomerDAO struct {
	collection *mongo.Collection
}

func NewCustomerDAO(db *mongo.Database) *CustomerDAO {
	return &CustomerDAO{
		collection: db.Collection("customers"),
	}
}

// Create 创建客户
func (dao *CustomerDAO) Create(ctx context.Context, customer *domain.Customer) error {
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, customer)
	if err != nil {
		return err
	}
	customer.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindById 根据ID查找客户
func (dao *CustomerDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Customer, error) {
	var customer domain.Customer
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

// FindByUserID 根据用户ID查找客户
func (dao *CustomerDAO) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*domain.Customer, error) {
	var customer domain.Customer
	err := dao.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

// FindList 查询客户列表
func (dao *CustomerDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Customer, int64, error) {
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

	var results []*domain.Customer
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update 更新客户信息
func (dao *CustomerDAO) Update(ctx context.Context, customer *domain.Customer) error {
	customer.UpdatedAt = time.Now()
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": customer.ID},
		bson.M{"$set": customer},
	)
	return err
}

// Delete 删除客户
func (dao *CustomerDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

