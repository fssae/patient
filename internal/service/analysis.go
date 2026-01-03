package service

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository"
	"context"
	"time"
)

// AnalysisService 定义了分析相关的业务接口
type AnalysisService interface {
	// RecordAnalysis 记录分析结果
	RecordAnalysis(ctx context.Context, msg *domain.AnalysisLog) error
	// GetAnalysisLogs 分页获取分析日志
	GetAnalysisLogs(ctx context.Context, page, size int64) ([]*domain.AnalysisLog, int64, error)
}

type analysisService struct {
	repo repository.AnalysisRepository
}

func NewAnalysisService(repo repository.AnalysisRepository) AnalysisService {
	return &analysisService{repo: repo}
}

func (s *analysisService) RecordAnalysis(ctx context.Context, log *domain.AnalysisLog) error {
	if log.CreatedAt == 0 {
		log.CreatedAt = time.Now().Unix()
	}
	return s.repo.Save(ctx, log)
}

func (s *analysisService) GetAnalysisLogs(ctx context.Context, page, size int64) ([]*domain.AnalysisLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 20
	}
	return s.repo.GetList(ctx, page, size)
}
