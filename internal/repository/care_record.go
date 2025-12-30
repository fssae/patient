package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CareRecordRepository interface {
	Create(ctx context.Context, record *domain.CareRecord) error
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.CareRecord, error)
	FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CareRecord, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type careRecordRepository struct {
	dao *dao.CareRecordDAO
}

func NewCareRecordRepository(dao *dao.CareRecordDAO) CareRecordRepository {
	return &careRecordRepository{
		dao: dao,
	}
}

func (r *careRecordRepository) Create(ctx context.Context, record *domain.CareRecord) error {
	return r.dao.Create(ctx, record)
}

func (r *careRecordRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.CareRecord, error) {
	return r.dao.FindById(ctx, id)
}

func (r *careRecordRepository) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CareRecord, int64, error) {
	return r.dao.FindList(ctx, filter, skip, limit)
}

func (r *careRecordRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	return r.dao.Update(ctx, id, bson.M(updates))
}

func (r *careRecordRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return r.dao.Delete(ctx, id)
}
