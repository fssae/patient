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

type ServiceDAO struct {
	collection *mongo.Collection
}

func NewServiceDAO(db *mongo.Database) *ServiceDAO {
	return &ServiceDAO{
		collection: db.Collection("services"),
	}
}

// Create 创建服务项目
func (dao *ServiceDAO) Create(ctx context.Context, service *domain.Service) error {
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, service)
	if err != nil {
		return err
	}
	service.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindById 根据ID查找服务项目
func (dao *ServiceDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Service, error) {
	var service domain.Service
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &service, nil
}

// FindList 查询服务项目列表
func (dao *ServiceDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Service, int64, error) {
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

	var results []*domain.Service
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update 更新服务项目信息
func (dao *ServiceDAO) Update(ctx context.Context, service *domain.Service) error {
	service.UpdatedAt = time.Now()
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": service.ID},
		bson.M{"$set": service},
	)
	return err
}

// Delete 删除服务项目
func (dao *ServiceDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

