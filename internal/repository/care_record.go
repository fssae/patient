package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CareRecordRepository interface {
	Create(ctx context.Context, record *domain.CareRecords) error
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.CareRecords, error)
	FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CareRecords, int64, error)
	FindRecordsByCustomerId(ctx context.Context, customerID primitive.ObjectID, startTime, endTime *time.Time, skip, limit int64) (*domain.CareRecords, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	AppendRecord(ctx context.Context, id primitive.ObjectID, item domain.RecordItems) error
}

type careRecordRepository struct {
	dao *dao.CareRecordDAO
}

func NewCareRecordRepository(dao *dao.CareRecordDAO) CareRecordRepository {
	return &careRecordRepository{
		dao: dao,
	}
}

func (r *careRecordRepository) Create(ctx context.Context, record *domain.CareRecords) error {
	return r.dao.Create(ctx, record)
}

func (r *careRecordRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.CareRecords, error) {
	return r.dao.FindById(ctx, id)
}

func (r *careRecordRepository) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CareRecords, int64, error) {
	return r.dao.FindList(ctx, filter, skip, limit)
}

func (r *careRecordRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	return r.dao.Update(ctx, id, bson.M(updates))
}

func (r *careRecordRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

func (r *careRecordRepository) AppendRecord(ctx context.Context, id primitive.ObjectID, item domain.RecordItems) error {
	return r.dao.AppendRecord(ctx, id, item)
}

func (r *careRecordRepository) FindRecordsByCustomerId(ctx context.Context, customerID primitive.ObjectID, startTime, endTime *time.Time, skip, limit int64) (*domain.CareRecords, int64, error) {
	return r.dao.FindRecordsByCustomerId(ctx, customerID, startTime, endTime, skip, limit)
}
