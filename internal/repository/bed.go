package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BedRepository interface {
	Create(ctx context.Context, bed *domain.Bed) error
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.Bed, error)
	FindByRoomID(ctx context.Context, roomID primitive.ObjectID) ([]*domain.Bed, error)
	FindByCustomerID(ctx context.Context, customerID primitive.ObjectID) (*domain.Bed, error)
	FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Bed, int64, error)
	Update(ctx context.Context, bed *domain.Bed) error
	AssignToCustomer(ctx context.Context, bedID, customerID primitive.ObjectID) error
	Release(ctx context.Context, bedID primitive.ObjectID) error
}

type bedRepository struct {
	dao *dao.BedDAO
}

func NewBedRepository(dao *dao.BedDAO) BedRepository {
	return &bedRepository{
		dao: dao,
	}
}

func (r *bedRepository) Create(ctx context.Context, bed *domain.Bed) error {
	return r.dao.Create(ctx, bed)
}

func (r *bedRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Bed, error) {
	return r.dao.FindById(ctx, id)
}

func (r *bedRepository) FindByRoomID(ctx context.Context, roomID primitive.ObjectID) ([]*domain.Bed, error) {
	return r.dao.FindByRoomID(ctx, roomID)
}

func (r *bedRepository) FindByCustomerID(ctx context.Context, customerID primitive.ObjectID) (*domain.Bed, error) {
	return r.dao.FindByCustomerID(ctx, customerID)
}

func (r *bedRepository) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Bed, int64, error) {
	return r.dao.FindList(ctx, filter, skip, limit)
}

func (r *bedRepository) Update(ctx context.Context, bed *domain.Bed) error {
	return r.dao.Update(ctx, bed)
}

func (r *bedRepository) AssignToCustomer(ctx context.Context, bedID, customerID primitive.ObjectID) error {
	return r.dao.AssignToCustomer(ctx, bedID, customerID)
}

func (r *bedRepository) Release(ctx context.Context, bedID primitive.ObjectID) error {
	return r.dao.Release(ctx, bedID)
}

