package dao

import (
	"classroom-analysis/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomDAO struct {
	BaseDAO[domain.Room]
}

func NewRoomDAO(db *mongo.Database) *RoomDAO {
	return &RoomDAO{
		BaseDAO: NewBaseDAO[domain.Room](db, "rooms"),
	}
}

// Create 创建房间
func (dao *RoomDAO) Create(ctx context.Context, room *domain.Room) error {
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()
	result, err := dao.Coll.InsertOne(ctx, room)
	if err != nil {
		return err
	}
	room.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByNumber 根据房间号查找
func (dao *RoomDAO) FindByNumber(ctx context.Context, number string) (*domain.Room, error) {
	return dao.FindOne(ctx, bson.M{"number": number})
}

// FindList 查询房间列表
func (dao *RoomDAO) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Room, int64, error) {
	return dao.BaseDAO.FindList(ctx, filter, skip, limit, nil)
}

// Update 更新房间信息
func (dao *RoomDAO) Update(ctx context.Context, room *domain.Room) error {
	room.UpdatedAt = time.Now()
	_, err := dao.Coll.UpdateOne(
		ctx,
		bson.M{"_id": room.ID},
		bson.M{"$set": room},
	)
	return err
}

// Delete 删除房间
func (dao *RoomDAO) Delete(ctx context.Context, id primitive.ObjectID) error {
	return dao.DeleteOne(ctx, id)
}
