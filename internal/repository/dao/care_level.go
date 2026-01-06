package dao

import (
	"classroom-analysis/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CareLevelDAO struct {
	BaseDAO[domain.CareLevel]
}

func NewCareLevelDAO(db *mongo.Database) *CareLevelDAO {
	return &CareLevelDAO{
		BaseDAO: NewBaseDAO[domain.CareLevel](db, "care_levels"),
	}
}

// Create 创建护理级别
func (dao *CareLevelDAO) Create(ctx context.Context, level *domain.CareLevel) error {
	level.CreatedAt = time.Now()
	level.UpdatedAt = time.Now()
	result, err := dao.Coll.InsertOne(ctx, level)
	if err != nil {
		return err
	}
	level.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindList 查询护理级别列表
func (dao *CareLevelDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CareLevel, int64, error) {
	return dao.BaseDAO.FindList(ctx, filter, skip, limit, nil)
}

// Update 更新护理级别信息
func (dao *CareLevelDAO) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	_, err := dao.Coll.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	return err
}

// Delete 删除护理级别
func (dao *CareLevelDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	return dao.DeleteOne(ctx, id)
}
