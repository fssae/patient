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

type DietPlanDAO struct {
	collection *mongo.Collection
}

func NewDietPlanDAO(db *mongo.Database) *DietPlanDAO {
	return &DietPlanDAO{
		collection: db.Collection("diet_plans"),
	}
}

// Create 创建膳食计划
func (dao *DietPlanDAO) Create(ctx context.Context, plan *domain.DietPlan) error {
	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, plan)
	if err != nil {
		return err
	}
	plan.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindById 根据ID查找膳食计划
func (dao *DietPlanDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.DietPlan, error) {
	var plan domain.DietPlan
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&plan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

// FindList 查询膳食计划列表
func (dao *DietPlanDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.DietPlan, int64, error) {
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

	var results []*domain.DietPlan
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update 更新膳食计划信息
func (dao *DietPlanDAO) Update(ctx context.Context, plan *domain.DietPlan) error {
	plan.UpdatedAt = time.Now()
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": plan.ID},
		bson.M{"$set": plan},
	)
	return err
}

// Delete 删除膳食计划
func (dao *DietPlanDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

