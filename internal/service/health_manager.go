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

type HealthManagerService struct {
	repo repository.HealthManagerRepository
}

func NewHealthManagerService(repo repository.HealthManagerRepository) *HealthManagerService {
	return &HealthManagerService{
		repo: repo,
	}
}

// Create 创建健康管家
func (s *HealthManagerService) Create(ctx context.Context, req *domain.HealthManager) error {
	manager := &domain.HealthManager{
		Name:      req.Name,
		Phone:     req.Phone,
		Email:     req.Email,
		Specialty: req.Specialty,
		Status:    "在职",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.repo.Create(ctx, manager)
}

// GetById 根据ID获取健康管家
func (s *HealthManagerService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.HealthManager, error) {
	return s.repo.FindById(ctx, id)
}

// GetList 获取健康管家列表
func (s *HealthManagerService) GetList(ctx context.Context, status string, skip, limit int64) ([]*domain.HealthManager, int64, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	return s.repo.FindList(ctx, filter, skip, limit)
}

// Update 更新健康管家信息
func (s *HealthManagerService) Update(ctx context.Context, id primitive.ObjectID, req *domain.HealthManager) error {
	manager, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if manager == nil {
		return errors.New("健康管家不存在")
	}

	manager.Name = req.Name
	manager.Phone = req.Phone
	manager.Email = req.Email
	manager.Specialty = req.Specialty
	manager.Status = req.Status
	manager.UpdatedAt = time.Now()

	return s.repo.Update(ctx, manager)
}

// Delete 删除健康管家
func (s *HealthManagerService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

