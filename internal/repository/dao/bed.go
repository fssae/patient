package dao

import (
	"classroom-analysis/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BedDAO struct {
	BaseDAO[domain.Bed]
}

func NewBedDAO(db *mongo.Database) *BedDAO {
	return &BedDAO{
		BaseDAO: NewBaseDAO[domain.Bed](db, "beds"),
	}
}

// Create 创建床位
func (dao *BedDAO) Create(ctx context.Context, bed *domain.Bed) error {
	bed.CreatedAt = time.Now()
	bed.UpdatedAt = time.Now()
	result, err := dao.Coll.InsertOne(ctx, bed)
	if err != nil {
		return err
	}
	bed.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByRoomID 根据房间ID查找床位列表
func (dao *BedDAO) FindByRoomID(ctx context.Context, roomID primitive.ObjectID) ([]*domain.Bed, error) {
	cursor, err := dao.Coll.Find(ctx, bson.M{"room_id": roomID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*domain.Bed
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// FindByCustomerID 根据客户ID查找床位
func (dao *BedDAO) FindByCustomerID(ctx context.Context, customerID primitive.ObjectID) (*domain.Bed, error) {
	return dao.FindOne(ctx, bson.M{"customer_id": customerID})
}

// FindList 查询床位列表（覆盖基类方法以使用默认排序）
func (dao *BedDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Bed, int64, error) {
	return dao.BaseDAO.FindList(ctx, filter, skip, limit, nil)
}

// Update 更新床位信息
func (dao *BedDAO) Update(ctx context.Context, bed *domain.Bed) error {
	bed.UpdatedAt = time.Now()
	_, err := dao.Coll.UpdateOne(
		ctx,
		bson.M{"_id": bed.ID},
		bson.M{"$set": bed},
	)
	return err
}

// AssignToCustomer 分配床位给客户
func (dao *BedDAO) AssignToCustomer(ctx context.Context, bedID, customerID primitive.ObjectID) error {
	_, err := dao.Coll.UpdateOne(
		ctx,
		bson.M{"_id": bedID},
		bson.M{"$set": bson.M{
			"customer_id": customerID,
			"status":      "占用",
			"updated_at":  time.Now(),
		}},
	)
	return err
}

// Release 释放床位
func (dao *BedDAO) Release(ctx context.Context, bedID primitive.ObjectID) error {
	_, err := dao.Coll.UpdateOne(
		ctx,
		bson.M{"_id": bedID},
		bson.M{"$set": bson.M{
			"customer_id": primitive.NilObjectID,
			"status":      "空闲",
			"updated_at":  time.Now(),
		}},
	)
	return err
}

// Delete 删除床位（使用BaseDAO的DeleteOne）
func (dao *BedDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	return dao.DeleteOne(ctx, id)
}
