package dao

import (
	"classroom-analysis/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthManagerDAO struct {
	BaseDAO[domain.HealthManager]
}

func NewHealthManagerDAO(db *mongo.Database) *HealthManagerDAO {
	return &HealthManagerDAO{
		BaseDAO: NewBaseDAO[domain.HealthManager](db, "health_managers"),
	}
}

// Create 创建健康管家
func (dao *HealthManagerDAO) Create(ctx context.Context, manager *domain.HealthManager) error {
	manager.CreatedAt = time.Now()
	manager.UpdatedAt = time.Now()
	result, err := dao.Coll.InsertOne(ctx, manager)
	if err != nil {
		return err
	}
	manager.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByPhone 根据手机号查找健康管家
func (dao *HealthManagerDAO) FindByPhone(ctx context.Context, phone string) (*domain.HealthManager, error) {
	return dao.FindOne(ctx, bson.M{"phone": phone})
}

// FindList 查询健康管家列表
func (dao *HealthManagerDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.HealthManager, int64, error) {
	return dao.BaseDAO.FindList(ctx, filter, skip, limit, nil)
}

// Update 更新健康管家信息
func (dao *HealthManagerDAO) Update(ctx context.Context, manager *domain.HealthManager) error {
	manager.UpdatedAt = time.Now()
	_, err := dao.Coll.UpdateOne(
		ctx,
		bson.M{"_id": manager.ID},
		bson.M{"$set": manager},
	)
	return err
}

// Delete 删除健康管家
func (dao *HealthManagerDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	return dao.DeleteOne(ctx, id)
}
