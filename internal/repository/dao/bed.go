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

type BedDAO struct {
	collection *mongo.Collection
}

func NewBedDAO(db *mongo.Database) *BedDAO {
	return &BedDAO{
		collection: db.Collection("beds"),
	}
}

// Create 创建床位
func (dao *BedDAO) Create(ctx context.Context, bed *domain.Bed) error {
	bed.CreatedAt = time.Now()
	bed.UpdatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, bed)
	if err != nil {
		return err
	}
	bed.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindById 根据ID查找床位
func (dao *BedDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Bed, error) {
	var bed domain.Bed
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bed)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &bed, nil
}

// FindByRoomID 根据房间ID查找床位列表
func (dao *BedDAO) FindByRoomID(ctx context.Context, roomID primitive.ObjectID) ([]*domain.Bed, error) {
	cursor, err := dao.collection.Find(ctx, bson.M{"room_id": roomID})
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
	var bed domain.Bed
	err := dao.collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&bed)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &bed, nil
}

// FindList 查询床位列表
func (dao *BedDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Bed, int64, error) {
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

	var results []*domain.Bed
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update 更新床位信息
func (dao *BedDAO) Update(ctx context.Context, bed *domain.Bed) error {
	bed.UpdatedAt = time.Now()
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": bed.ID},
		bson.M{"$set": bed},
	)
	return err
}

// AssignToCustomer 分配床位给客户
func (dao *BedDAO) AssignToCustomer(ctx context.Context, bedID, customerID primitive.ObjectID) error {
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": bedID},
		bson.M{"$set": bson.M{
			"customer_id": customerID,
			"status":       "占用",
			"updated_at":   time.Now(),
		}},
	)
	return err
}

// Release 释放床位
func (dao *BedDAO) Release(ctx context.Context, bedID primitive.ObjectID) error {
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": bedID},
		bson.M{"$set": bson.M{
			"customer_id": primitive.NilObjectID,
			"status":       "空闲",
			"updated_at":   time.Now(),
		}},
	)
	return err
}

