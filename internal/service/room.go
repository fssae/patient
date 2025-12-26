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

type RoomService struct {
	repo repository.RoomRepository
}

func NewRoomService(repo repository.RoomRepository) *RoomService {
	return &RoomService{
		repo: repo,
	}
}

// Create 创建房间
func (s *RoomService) Create(ctx context.Context, req *domain.Room) error {
	// 检查房间号是否已存在
	existing, err := s.repo.FindByNumber(ctx, req.Number)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("房间号已存在")
	}

	room := &domain.Room{
		Number:      req.Number,
		Floor:       req.Floor,
		Type:         req.Type,
		Capacity:     req.Capacity,
		Status:       "可用",
		Description:  req.Description,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return s.repo.Create(ctx, room)
}

// GetById 根据ID获取房间
func (s *RoomService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Room, error) {
	return s.repo.FindById(ctx, id)
}

// GetList 获取房间列表
func (s *RoomService) GetList(ctx context.Context, status string, floor int, skip, limit int64) ([]*domain.Room, int64, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	if floor > 0 {
		filter["floor"] = floor
	}
	return s.repo.FindList(ctx, filter, skip, limit)
}

// Update 更新房间信息
func (s *RoomService) Update(ctx context.Context, id primitive.ObjectID, req *domain.Room) error {
	room, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("房间不存在")
	}

	room.Number = req.Number
	room.Floor = req.Floor
	room.Type = req.Type
	room.Capacity = req.Capacity
	room.Status = req.Status
	room.Description = req.Description
	room.UpdatedAt = time.Now()

	return s.repo.Update(ctx, room)
}

// Delete 删除房间
func (s *RoomService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

