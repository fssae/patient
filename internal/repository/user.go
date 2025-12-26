package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByPhone(ctx context.Context, phone string) (*domain.User, error)
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type userRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(userDAO *dao.UserDAO) UserRepository {
	return &userRepository{
		dao: userDAO,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.dao.Create(ctx, user)
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	return r.dao.FindByPhone(ctx, phone)
}

func (r *userRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	return r.dao.FindById(ctx, id)
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.dao.Update(ctx, user)
}

