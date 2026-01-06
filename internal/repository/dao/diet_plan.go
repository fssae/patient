package dao

import (
	"classroom-analysis/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DietPlanDAO struct {
	BaseDAO[domain.DietPlan]
}

func NewDietPlanDAO(db *mongo.Database) *DietPlanDAO {
	return &DietPlanDAO{
		BaseDAO: NewBaseDAO[domain.DietPlan](db, "diet_plans"),
	}
}

// Create 创建膳食计划
func (dao *DietPlanDAO) Create(ctx context.Context, plan *domain.DietPlan) error {
	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()
	result, err := dao.Coll.InsertOne(ctx, plan)
	if err != nil {
		return err
	}
	plan.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindList 查询膳食计划列表
func (dao *DietPlanDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.DietPlan, int64, error) {
	return dao.BaseDAO.FindList(ctx, filter, skip, limit, nil)
}

// Update 更新膳食计划信息
func (dao *DietPlanDAO) Update(ctx context.Context, plan *domain.DietPlan) error {
	plan.UpdatedAt = time.Now()
	_, err := dao.Coll.UpdateOne(
		ctx,
		bson.M{"_id": plan.ID},
		bson.M{"$set": plan},
	)
	return err
}

// Delete 删除膳食计划
func (dao *DietPlanDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	return dao.DeleteOne(ctx, id)
}
