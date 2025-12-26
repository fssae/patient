package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecordRepository interface {
	Create(ctx context.Context, record *domain.Record) error
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.Record, error)
	FindByCustomerID(ctx context.Context, customerID primitive.ObjectID) ([]*domain.Record, error)
	FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Record, int64, error)
	Update(ctx context.Context, record *domain.Record) error
}

type recordRepository struct {
	dao *dao.RecordDAO
}

func NewRecordRepository(dao *dao.RecordDAO) RecordRepository {
	return &recordRepository{
		dao: dao,
	}
}

func (r *recordRepository) Create(ctx context.Context, record *domain.Record) error {
	return r.dao.Create(ctx, record)
}

func (r *recordRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Record, error) {
	return r.dao.FindById(ctx, id)
}

func (r *recordRepository) FindByCustomerID(ctx context.Context, customerID primitive.ObjectID) ([]*domain.Record, error) {
	return r.dao.FindByCustomerID(ctx, customerID)
}

func (r *recordRepository) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Record, int64, error) {
	return r.dao.FindList(ctx, filter, skip, limit)
}

func (r *recordRepository) Update(ctx context.Context, record *domain.Record) error {
	return r.dao.Update(ctx, record)
}

