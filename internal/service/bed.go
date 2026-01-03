package service

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository"
	"context"
	"errors"
	"strconv"
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
func (s *BedService) Create(ctx context.Context, req *domain.CreateBed) error {
	// 验证房间是否存在
	RId, err := primitive.ObjectIDFromHex(req.RoomId)
	if err != nil {
		return errors.New("id格式错误")
	}
	room, err := s.roomRepo.FindById(ctx, RId)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("房间不存在")
	}
	bed := &domain.Bed{
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
func (s *BedService) GetList(ctx context.Context, roomNumber, bedNumber, status string, skip, limit int64) ([]*domain.Bed, int64, error) {
	filter := bson.M{}

	if roomNumber != "" {
		// 先查询房间ID
		room, err := s.roomRepo.FindByNumber(ctx, roomNumber)
		if err != nil {
			return nil, 0, err
		}
		if room != nil {
			filter["room_id"] = room.ID
		} else {
			// 房间不存在，直接返回空结果
			return []*domain.Bed{}, 0, nil
		}
	}

	if bedNumber != "" {
		filter["number"] = bson.M{"$regex": bedNumber, "$options": "i"}
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

// GetBedOptions 获取床位递归选项列表 (房间 -> 床位)
func (s *BedService) GetBedOptions(ctx context.Context) ([]map[string]interface{}, error) {
	// 获取所有房间
	var rooms []*domain.Room
	rooms, _, err := s.roomRepo.FindList(ctx, bson.M{}, 0, 0)
	if err != nil {
		return nil, err
	}

	// 获取所有床位
	beds, _, err := s.repo.FindList(ctx, bson.M{}, 0, 0)
	if err != nil {
		return nil, err
	}

	// 组装结果
	var options []map[string]interface{}
	for _, room := range rooms {
		roomOption := map[string]interface{}{
			"value": room.ID.Hex(),
			"label": room.Number,
		}

		var children []map[string]interface{}
		for _, bed := range beds {
			if bed.RoomID == room.ID {
				children = append(children, map[string]interface{}{
					"value":  bed.ID.Hex(),
					"label":  bed.Number,
					"status": bed.Status,
				})
			}
		}

		if len(children) > 0 {
			roomOption["children"] = children
			options = append(options, roomOption)
		}
	}
	return options, nil
}

// GetRoomOptions 获取房间递归选项列表 (楼层 -> 房间)
func (s *BedService) GetRoomOptions(ctx context.Context) ([]map[string]interface{}, error) {
	// 获取所有房间
	rooms, _, err := s.roomRepo.FindList(ctx, bson.M{}, 0, 0)
	if err != nil {
		return nil, err
	}

	// 按楼层分组
	floorMap := make(map[int][]*domain.Room)
	for _, room := range rooms {
		floorMap[room.Floor] = append(floorMap[room.Floor], room)
	}

	// 组装结果
	var options []map[string]interface{}
	for floor, floorRooms := range floorMap {
		floorOption := map[string]interface{}{
			"value": floor,
			"label": strconv.Itoa(floor) + "层",
		}

		var children []map[string]interface{}
		for _, room := range floorRooms {
			children = append(children, map[string]interface{}{
				"value":  room.ID.Hex(),
				"label":  room.Number,
				"status": room.Status,
			})
		}

		if len(children) > 0 {
			floorOption["children"] = children
			options = append(options, floorOption)
		}
	}

	return options, nil
}
