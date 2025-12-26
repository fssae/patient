package service

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DietPlanService struct {
	repo repository.DietPlanRepository
}

func NewDietPlanService(repo repository.DietPlanRepository) *DietPlanService {
	return &DietPlanService{
		repo: repo,
	}
}

// Create 创建膳食计划
func (s *DietPlanService) Create(ctx context.Context, req *domain.DietPlan) error {
	plan := &domain.DietPlan{
		Name:        req.Name,
		Description: req.Description,
		WeekMenu:    req.WeekMenu,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.repo.Create(ctx, plan)
}

// GetById 根据ID获取膳食计划
func (s *DietPlanService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.DietPlan, error) {
	return s.repo.FindById(ctx, id)
}

// GetList 获取膳食计划列表
func (s *DietPlanService) GetList(ctx context.Context, skip, limit int64) ([]*domain.DietPlan, int64, error) {
	filter := bson.M{}
	return s.repo.FindList(ctx, filter, skip, limit)
}

// Update 更新膳食计划信息
func (s *DietPlanService) Update(ctx context.Context, id primitive.ObjectID, req *domain.DietPlan) error {
	plan, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if plan == nil {
		return errors.New("膳食计划不存在")
	}

	plan.Name = req.Name
	plan.Description = req.Description
	plan.WeekMenu = req.WeekMenu
	plan.UpdatedAt = time.Now()

	return s.repo.Update(ctx, plan)
}

// Delete 删除膳食计划
func (s *DietPlanService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

