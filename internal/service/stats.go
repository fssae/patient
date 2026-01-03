package service

import (
	"context"
)

type SystemStats struct {
	TotalPatients int `json:"totalPatients"`
	TodayAlerts   int `json:"todayAlerts"`
	HighRiskCount int `json:"highRiskCount"`
	OnlineStaff   int `json:"onlineStaff"`
}

type TrendData struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type StatsService interface {
	GetSummary(ctx context.Context) (*SystemStats, error)
	GetTrend(ctx context.Context) ([]TrendData, error)
}

type statsService struct {
	// 可以在这里注入 repository
}

func NewStatsService() StatsService {
	return &statsService{}
}

func (s *statsService) GetSummary(ctx context.Context) (*SystemStats, error) {
	// TODO: 从数据库统计真实数据
	// 目前返回模拟数据以符合接口规范
	return &SystemStats{
		TotalPatients: 128,
		TodayAlerts:   42,
		HighRiskCount: 5,
		OnlineStaff:   12,
	}, nil
}

func (s *statsService) GetTrend(ctx context.Context) ([]TrendData, error) {
	// TODO: 从数据库统计过去 24 小时趋势
	// 目前返回模拟数据以符合接口规范
	return []TrendData{
		{Name: "00:00", Count: 2},
		{Name: "04:00", Count: 1},
		{Name: "08:00", Count: 5},
		{Name: "12:00", Count: 8},
		{Name: "16:00", Count: 4},
		{Name: "20:00", Count: 6},
		{Name: "23:59", Count: 3},
	}, nil
}
