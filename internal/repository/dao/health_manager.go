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

type HealthManagerDAO struct {
	collection *mongo.Collection
}

func NewHealthManagerDAO(db *mongo.Database) *HealthManagerDAO {
	return &HealthManagerDAO{
		collection: db.Collection("health_managers"),
	}
}

// Create 创建健康管家
func (dao *HealthManagerDAO) Create(ctx context.Context, manager *domain.HealthManager) error {
	manager.CreatedAt = time.Now()
	manager.UpdatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, manager)
	if err != nil {
		return err
	}
	manager.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindById 根据ID查找健康管家
func (dao *HealthManagerDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.HealthManager, error) {
	var manager domain.HealthManager
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&manager)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &manager, nil
}

// FindList 查询健康管家列表
func (dao *HealthManagerDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.HealthManager, int64, error) {
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

	var results []*domain.HealthManager
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update 更新健康管家信息
func (dao *HealthManagerDAO) Update(ctx context.Context, manager *domain.HealthManager) error {
	manager.UpdatedAt = time.Now()
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": manager.ID},
		bson.M{"$set": manager},
	)
	return err
}

// Delete 删除健康管家
func (dao *HealthManagerDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

