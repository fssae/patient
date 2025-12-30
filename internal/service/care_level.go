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
func (s *CareLevelService) GetList(ctx context.Context, level int64, skip, limit int64) ([]*domain.CareLevel, int64, error) {
	filter := bson.M{}
	if level > 0 {
		filter["level"] = level
	}
	return s.repo.FindList(ctx, filter, skip, limit)
}

// Update 更新护理级别信息
func (s *CareLevelService) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	level, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if level == nil {
		return errors.New("护理级别不存在")
	}
	return s.repo.Update(ctx, id, updates)
}

// Delete 删除护理级别
func (s *CareLevelService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

// GetMenu 获取护理级别菜单（用于下拉选择）
func (s *CareLevelService) GetMenu(ctx context.Context) ([]*domain.CareLevel, error) {
	// 菜单通常需要全部数据，或者根据业务过滤掉已停用的
	list, _, err := s.repo.FindList(ctx, bson.M{}, 0, 1000) // 获取前1000条，基本覆盖所有级别
	return list, err
}
