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
func (s *HealthManagerService) GetList(ctx context.Context, query map[string]interface{}, skip, limit int64) ([]*domain.HealthManager, int64, error) {
	filter := bson.M{}
	for k, v := range query {
		if v == "" || v == nil {
			continue
		}
		switch k {
		case "id":
			if str, ok := v.(string); ok {
				if oid, err := primitive.ObjectIDFromHex(str); err == nil {
					filter["_id"] = oid
				}
			}
		case "name", "phone", "id_card", "email", "specialty", "status": // 模糊查询
			if str, ok := v.(string); ok {
				filter[k] = primitive.Regex{Pattern: str, Options: "i"}
			}
		case "created_at", "updated_at":
			if dateStr, ok := v.(string); ok {
				if date, err := time.Parse("2006-01-02", dateStr); err == nil {
					startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
					endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)
					filter[k] = bson.M{
						"$gte": startOfDay,
						"$lte": endOfDay,
					}
				}
			}
		default:
			filter[k] = v
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

	// 验证是否有客户关联

	return s.repo.Delete(ctx, id)
}
