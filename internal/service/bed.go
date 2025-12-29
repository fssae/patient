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

type BedService struct {
	repo     repository.BedRepository
	roomRepo repository.RoomRepository
}

func NewBedService(repo repository.BedRepository, roomRepo repository.RoomRepository) *BedService {
	return &BedService{
		repo:     repo,
		roomRepo: roomRepo,
	}
}

// Create 创建床位
func (s *BedService) Create(ctx context.Context, req *domain.CreateBedRequest) error {
	// 验证房间是否存在
	room, err := s.roomRepo.FindById(ctx, req.BedID)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("房间不存在")
	}
	bed := &domain.Bed{
		RoomID:    req.BedID,
		Number:    room.Number,
		Status:    "空闲",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.repo.Create(ctx, bed)
}

// GetById 根据ID获取床位
func (s *BedService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Bed, error) {
	return s.repo.FindById(ctx, id)
}

// GetByRoomID 根据房间ID获取床位列表
func (s *BedService) GetByRoomID(ctx context.Context, roomID primitive.ObjectID) ([]*domain.Bed, error) {
	return s.repo.FindByRoomID(ctx, roomID)
}

// GetList 获取床位列表
func (s *BedService) GetList(ctx context.Context, roomID primitive.ObjectID, status string, skip, limit int64) ([]*domain.Bed, int64, error) {
	filter := bson.M{}
	if !roomID.IsZero() {
		filter["room_id"] = roomID
	}
	if status != "" {
		filter["status"] = status
	}
	return s.repo.FindList(ctx, filter, skip, limit)
}

// AssignToCustomer 分配床位给客户
func (s *BedService) AssignToCustomer(ctx context.Context, bedID, customerID primitive.ObjectID) error {
	bed, err := s.repo.FindById(ctx, bedID)
	if err != nil {
		return err
	}
	if bed == nil {
		return errors.New("床位不存在")
	}
	if bed.Status == "占用" {
		return errors.New("床位已被占用")
	}

	return s.repo.AssignToCustomer(ctx, bedID, customerID)
}

// Release 释放床位
func (s *BedService) Release(ctx context.Context, bedID primitive.ObjectID) error {
	return s.repo.Release(ctx, bedID)
}

// Update 更新床位信息
func (s *BedService) Update(ctx context.Context, id primitive.ObjectID, req *domain.Bed) error {
	bed, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if bed == nil {
		return errors.New("床位不存在")
	}

	bed.Number = req.Number
	bed.Status = req.Status
	bed.UpdatedAt = time.Now()

	return s.repo.Update(ctx, bed)
}

// Delete 删除床位
func (s *BedService) Delete(ctx context.Context, id primitive.ObjectID) error {
	bed, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if bed == nil {
		return errors.New("床位不存在")
	}
	if bed.Status == "占用" {
		return errors.New("床位已被占用，无法删除")
	}
	return s.repo.Delete(ctx, id)
}
