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

func (s *CareRecordService) Create(ctx context.Context, record *domain.CareRecord) error {
	return s.repo.Create(ctx, record)
}

func (s *CareRecordService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.CareRecord, error) {
	return s.repo.FindById(ctx, id)
}

func (s *CareRecordService) GetList(ctx context.Context, customerName string, skip, limit int64) ([]*domain.CareRecord, int64, error) {
	filter := bson.M{}
	if customerName != "" {
		filter["customer_name"] = primitive.Regex{Pattern: customerName, Options: "i"}
	}
	return s.repo.FindList(ctx, filter, skip, limit)
}

func (s *CareRecordService) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	return s.repo.Update(ctx, id, updates)
}

func (s *CareRecordService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}
