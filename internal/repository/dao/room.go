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

type RoomDAO struct {
	collection *mongo.Collection
}

func NewRoomDAO(db *mongo.Database) *RoomDAO {
	return &RoomDAO{
		collection: db.Collection("rooms"),
	}
}

// Create 创建房间
func (dao *RoomDAO) Create(ctx context.Context, room *domain.Room) error {
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, room)
	if err != nil {
		return err
	}
	room.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindById 根据ID查找房间
func (dao *RoomDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Room, error) {
	var room domain.Room
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

// FindByNumber 根据房间号查找
func (dao *RoomDAO) FindByNumber(ctx context.Context, number string) (*domain.Room, error) {
	var room domain.Room
	err := dao.collection.FindOne(ctx, bson.M{"number": number}).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

// FindList 查询房间列表
func (dao *RoomDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Room, int64, error) {
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

	var results []*domain.Room
	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := dao.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// Update 更新房间信息
func (dao *RoomDAO) Update(ctx context.Context, room *domain.Room) error {
	room.UpdatedAt = time.Now()
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": room.ID},
		bson.M{"$set": room},
	)
	return err
}

// Delete 删除房间
func (dao *RoomDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := dao.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

