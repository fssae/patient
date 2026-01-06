package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerServiceRepository interface {
	Create(ctx context.Context, cs *domain.CustomerService) error
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.CustomerService, error)
	FindByCustomerID(ctx context.Context, customerID primitive.ObjectID) ([]*domain.CustomerService, error)
	FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CustomerService, int64, error)
	FindDistinctCustomerIDs(ctx context.Context, filter bson.M) ([]primitive.ObjectID, error)
	Update(ctx context.Context, cs *domain.CustomerService) error
}

type customerServiceRepository struct {
	dao *dao.CustomerServiceDAO
}

func NewCustomerServiceRepository(dao *dao.CustomerServiceDAO) CustomerServiceRepository {
	return &customerServiceRepository{
		dao: dao,
	}
}

func (r *customerServiceRepository) Create(ctx context.Context, cs *domain.CustomerService) error {
	return r.dao.Create(ctx, cs)
}

func (r *customerServiceRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.CustomerService, error) {
	return r.dao.FindById(ctx, id)
}

func (r *customerServiceRepository) FindByCustomerID(ctx context.Context, customerID primitive.ObjectID) ([]*domain.CustomerService, error) {
	return r.dao.FindByCustomerID(ctx, customerID)
}

func (r *customerServiceRepository) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.CustomerService, int64, error) {
	return r.dao.FindList(ctx, filter, skip, limit)
}

func (r *customerServiceRepository) FindDistinctCustomerIDs(ctx context.Context, filter bson.M) ([]primitive.ObjectID, error) {
	return r.dao.FindDistinctCustomerIDs(ctx, filter)
}

func (r *customerServiceRepository) Update(ctx context.Context, cs *domain.CustomerService) error {
	return r.dao.Update(ctx, cs)
}
