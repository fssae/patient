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

type CareLevelService struct {
	repo repository.CareLevelRepository
}

func NewCareLevelService(repo repository.CareLevelRepository) *CareLevelService {
	return &CareLevelService{
		repo: repo,
	}
}

// Create 创建护理级别
func (s *CareLevelService) Create(ctx context.Context, req *domain.CareLevel) error {
	level := &domain.CareLevel{
		Name:        req.Name,
		Level:       req.Level,
		Description: req.Description,
		Content:     req.Content,
		Price:       req.Price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.repo.Create(ctx, level)
}

// GetById 根据ID获取护理级别
func (s *CareLevelService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.CareLevel, error) {
	return s.repo.FindById(ctx, id)
}

// GetList 获取护理级别列表
func (s *CareLevelService) GetList(ctx context.Context, skip, limit int64) ([]*domain.CareLevel, int64, error) {
	filter := bson.M{}
	return s.repo.FindList(ctx, filter, skip, limit)
}

// Update 更新护理级别信息
func (s *CareLevelService) Update(ctx context.Context, id primitive.ObjectID, req *domain.CareLevel) error {
	level, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if level == nil {
		return errors.New("护理级别不存在")
	}

	level.Name = req.Name
	level.Level = req.Level
	level.Description = req.Description
	level.Content = req.Content
	level.Price = req.Price
	level.UpdatedAt = time.Now()

	return s.repo.Update(ctx, level)
}

// Delete 删除护理级别
func (s *CareLevelService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

