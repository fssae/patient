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

type CareLevelDAO struct {
	collection *mongo.Collection
}

func NewCareLevelDAO(db *mongo.Database) *CareLevelDAO {
	return &CareLevelDAO{
		collection: db.Collection("care_levels"),
	}
}

// Create 创建护理级别
func (dao *CareLevelDAO) Create(ctx context.Context, level *domain.CareLevel) error {
	level.CreatedAt = time.Now()
	level.UpdatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, level)
	if err != nil {
		return err
	}
	level.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindById 根据ID查找护理级别
func (dao *CareLevelDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.CareLevel, error) {
	var level domain.CareLevel
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&level)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &level, nil
}

// FindList 查询护理级别列表
func (dao *CareLevelDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CareLevel, int64, error) {
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

	var results []*domain.CareLevel
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update 更新护理级别信息
func (dao *CareLevelDAO) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	return err
}

// Delete 删除护理级别
func (dao *CareLevelDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
