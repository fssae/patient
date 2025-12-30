package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HealthManagerRepository interface {
	Create(ctx context.Context, manager *domain.HealthManager) error
	FindByPhone(ctx context.Context, phone string) (*domain.HealthManager, error)
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.HealthManager, error)
	FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.HealthManager, int64, error)
	Update(ctx context.Context, manager *domain.HealthManager) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type healthManagerRepository struct {
	dao *dao.HealthManagerDAO
}

func NewHealthManagerRepository(dao *dao.HealthManagerDAO) HealthManagerRepository {
	return &healthManagerRepository{
		dao: dao,
	}
}

func (r *healthManagerRepository) Create(ctx context.Context, manager *domain.HealthManager) error {
	return r.dao.Create(ctx, manager)
}

func (r *healthManagerRepository) FindByPhone(ctx context.Context, phone string) (*domain.HealthManager, error) {
	return r.dao.FindByPhone(ctx, phone)
}

func (r *healthManagerRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.HealthManager, error) {
	return r.dao.FindById(ctx, id)
}

func (r *healthManagerRepository) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.HealthManager, int64, error) {
	return r.dao.FindList(ctx, filter, skip, limit)
}

func (r *healthManagerRepository) Update(ctx context.Context, manager *domain.HealthManager) error {
	return r.dao.Update(ctx, manager)
}

func (r *healthManagerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return r.dao.Delete(ctx, id)
}
