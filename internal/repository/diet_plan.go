package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DietPlanRepository interface {
	Create(ctx context.Context, plan *domain.DietPlan) error
	FindById(ctx context.Context, id primitive.ObjectID) (*domain.DietPlan, error)
	FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.DietPlan, int64, error)
	Update(ctx context.Context, plan *domain.DietPlan) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type dietPlanRepository struct {
	dao *dao.DietPlanDAO
}

func NewDietPlanRepository(dao *dao.DietPlanDAO) DietPlanRepository {
	return &dietPlanRepository{
		dao: dao,
	}
}

func (r *dietPlanRepository) Create(ctx context.Context, plan *domain.DietPlan) error {
	return r.dao.Create(ctx, plan)
}

func (r *dietPlanRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.DietPlan, error) {
	return r.dao.FindById(ctx, id)
}

func (r *dietPlanRepository) FindList(ctx context.Context, filter bson.M, skip, limit int64) ([]*domain.DietPlan, int64, error) {
	return r.dao.FindList(ctx, filter, skip, limit)
}

func (r *dietPlanRepository) Update(ctx context.Context, plan *domain.DietPlan) error {
	return r.dao.Update(ctx, plan)
}

func (r *dietPlanRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

