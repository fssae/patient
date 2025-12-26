package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomRepository interface {
	Create(ctx context.Context, room *domain.Room) error
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.Room, error)
	FindByNumber(ctx context.Context, number string) (*domain.Room, error)
	FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Room, int64, error)
	Update(ctx context.Context, room *domain.Room) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type roomRepository struct {
	dao *dao.RoomDAO
}

func NewRoomRepository(dao *dao.RoomDAO) RoomRepository {
	return &roomRepository{
		dao: dao,
	}
}

func (r *roomRepository) Create(ctx context.Context, room *domain.Room) error {
	return r.dao.Create(ctx, room)
}

func (r *roomRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Room, error) {
	return r.dao.FindById(ctx, id)
}

func (r *roomRepository) FindByNumber(ctx context.Context, number string) (*domain.Room, error) {
	return r.dao.FindByNumber(ctx, number)
}

func (r *roomRepository) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Room, int64, error) {
	return r.dao.FindList(ctx, filter, skip, limit)
}

func (r *roomRepository) Update(ctx context.Context, room *domain.Room) error {
	return r.dao.Update(ctx, room)
}

func (r *roomRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

