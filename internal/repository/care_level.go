package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CareLevelRepository interface {
	Create(ctx context.Context, level *domain.CareLevel) error
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.CareLevel, error)
	FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CareLevel, int64, error)
	Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type careLevelRepository struct {
	dao *dao.CareLevelDAO
}

func NewCareLevelRepository(dao *dao.CareLevelDAO) CareLevelRepository {
	return &careLevelRepository{
		dao: dao,
	}
}

func (r *careLevelRepository) Create(ctx context.Context, level *domain.CareLevel) error {
	return r.dao.Create(ctx, level)
}

func (r *careLevelRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.CareLevel, error) {
	return r.dao.FindById(ctx, id)
}

func (r *careLevelRepository) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CareLevel, int64, error) {
	return r.dao.FindList(ctx, filter, skip, limit)
}

func (r *careLevelRepository) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	return r.dao.Update(ctx, id, updates)
}

func (r *careLevelRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return r.dao.Delete(ctx, id)
}
