package service

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository"
	"context"
	"time"

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

// GetList 获取单个用户的护理记录列表
// 根据 CustomerId 查询，按 care_time 筛选时间范围，返回的 total 为符合条件的记录条数
func (s *CareRecordService) GetList(ctx context.Context, customerName string, skip, limit int64, startDate, endDate string, customerID string) (*domain.CareRecords, int64, error) {
	// 如果没有 customerID，返回空结果
	if customerID == "" {
		return nil, 0, nil
	}

	// 解析 customerID
	objID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		return nil, 0, err
	}

	// 解析时间范围
	var startTime, endTime *time.Time
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			startTime = &t
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			// 结束日期包含当天，加一天
			t = t.Add(24 * time.Hour)
			endTime = &t
		}
	}

	return s.repo.FindRecordsByCustomerId(ctx, objID, startTime, endTime, skip, limit)
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
