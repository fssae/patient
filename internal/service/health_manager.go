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

func NewHealthManagerService(
	repo repository.HealthManagerRepository) *HealthManagerService {
	return &HealthManagerService{
		repo: repo,
	}
}

// Create 创建健康管家
func (s *HealthManagerService) Create(ctx context.Context, req *domain.HealthManager) error {
	// 验证手机号是否已注册
	existing, err := s.repo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("手机号已注册")
	}

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
func (s *HealthManagerService) GetList(ctx context.Context, query domain.HealthManagerQuery, skip, limit int64) ([]*domain.HealthManager, int64, error) {
	// 构建查询条件
	filter := bson.M{"status": "在职"} // 默认只查询在职人员

	// 姓名模糊查询
	if query.Name != "" {
		filter["name"] = bson.M{
			"$regex":   query.Name,
			"$options": "i",
		}
	}

	// 手机号模糊查询
	if query.Phone != "" {
		filter["phone"] = query.Phone
	}

	// 邮箱模糊查询
	if query.Email != "" {
		filter["email"] = query.Email
	}

	// 专业模糊查询
	if query.Specialty != "" {
		filter["specialty"] = bson.M{
			"$regex":   query.Specialty,
			"$options": "i",
		}
	}

	return s.repo.FindList(ctx, filter, skip, limit)
}

// Update 更新健康管家信息
func (s *HealthManagerService) Update(ctx context.Context, id primitive.ObjectID, req *domain.HealthManager) error {
	// 验证健康管家是否存在
	manager, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if manager == nil {
		return errors.New("健康管家不存在")
	}

	//如果手机号修改了，检查是否已存在
	if req.Phone != "" && req.Phone != manager.Phone {
		existing, err := s.repo.FindByPhone(ctx, req.Phone)
		if err != nil {
			return err
		}
		if existing != nil {
			return errors.New("手机号已注册")
		}
	}

	// 更新健康管家信息
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
	// 验证健康管家是否存在
	manager, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if manager == nil {
		return errors.New("健康管家不存在")
	}

	return s.repo.Delete(ctx, id)
}
