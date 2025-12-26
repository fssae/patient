package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.Customer, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID) (*domain.Customer, error)
	FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Customer, int64, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type customerRepository struct {
	dao *dao.CustomerDAO
}

func NewCustomerRepository(dao *dao.CustomerDAO) CustomerRepository {
	return &customerRepository{
		dao: dao,
	}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	return r.dao.Create(ctx, customer)
}

func (r *customerRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Customer, error) {
	return r.dao.FindById(ctx, id)
}

func (r *customerRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*domain.Customer, error) {
	return r.dao.FindByUserID(ctx, userID)
}

func (r *customerRepository) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.Customer, int64, error) {
	return r.dao.FindList(ctx, filter, skip, limit)
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	return r.dao.Update(ctx, customer)
}

func (r *customerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

