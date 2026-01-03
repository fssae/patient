package service

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CareRecordService struct {
	repo repository.CareRecordRepository
}

func NewCareRecordService(repo repository.CareRecordRepository) *CareRecordService {
	return &CareRecordService{
		repo: repo,
	}
}

// Create 创建护理记录
func (s *CareRecordService) Create(ctx context.Context, record *domain.CareRecords) error {
	return s.repo.Create(ctx, record)
}

// GetById 根据ID获取护理记录
func (s *CareRecordService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.CareRecords, error) {
	return s.repo.FindById(ctx, id)
}

// GetList 获取护理记录列表
func (s *CareRecordService) GetList(ctx context.Context, customerName string, skip, limit int64) ([]*domain.CareRecords, int64, error) {
	filter := bson.M{}
	if customerName != "" {
		filter["customer_name"] = primitive.Regex{Pattern: customerName, Options: "i"}
	}
	return s.repo.FindList(ctx, filter, skip, limit)
}

// Update 更新护理记录信息
func (s *CareRecordService) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	return s.repo.Update(ctx, id, updates)
}

// Delete 删除护理记录
func (s *CareRecordService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

// AddRecord 向护理记录中追加子项
func (s *CareRecordService) AddRecord(ctx context.Context, id primitive.ObjectID, item domain.RecordItems) error {
	return s.repo.AppendRecord(ctx, id, item)
}
