package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *domain.Service) error
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.Service, error)
	FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Service, int64, error)
	Update(ctx context.Context, service *domain.Service) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type serviceRepository struct {
	dao *dao.ServiceDAO
}

func NewServiceRepository(dao *dao.ServiceDAO) ServiceRepository {
	return &serviceRepository{
		dao: dao,
	}
}

func (r *serviceRepository) Create(ctx context.Context, service *domain.Service) error {
	return r.dao.Create(ctx, service)
}

func (r *serviceRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Service, error) {
	return r.dao.FindById(ctx, id)
}

func (r *serviceRepository) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Service, int64, error) {
	return r.dao.FindList(ctx, filter, skip, limit)
}

func (r *serviceRepository) Update(ctx context.Context, service *domain.Service) error {
	return r.dao.Update(ctx, service)
}

func (r *serviceRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

