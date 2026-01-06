package dao

import (
	"classroom-analysis/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerDAO struct {
	BaseDAO[domain.Customer]
}

func NewCustomerDAO(db *mongo.Database) *CustomerDAO {
	return &CustomerDAO{
		BaseDAO: NewBaseDAO[domain.Customer](db, "customers"),
	}
}

// Create 创建客户
func (dao *CustomerDAO) Create(ctx context.Context, customer *domain.Customer) error {
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	result, err := dao.Coll.InsertOne(ctx, customer)
	if err != nil {
		return err
	}
	customer.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByUserID 根据用户ID查找客户
func (dao *CustomerDAO) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*domain.Customer, error) {
	return dao.FindOne(ctx, bson.M{"user_id": userID})
}

// FindList 查询客户列表
func (dao *CustomerDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Customer, int64, error) {
	return dao.BaseDAO.FindList(ctx, filter, skip, limit, nil)
}

// Update 更新客户信息
func (dao *CustomerDAO) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	_, err := dao.Coll.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	return err
}

// Delete 删除客户
func (dao *CustomerDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	return dao.DeleteOne(ctx, id)
}
