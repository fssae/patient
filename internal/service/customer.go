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

type CustomerService struct {
	customerRepo repository.CustomerRepository
	userRepo     repository.UserRepository
	bedRepo      repository.BedRepository
	careRepo     repository.CareLevelRepository
	dietRepo     repository.DietPlanRepository
}

func NewCustomerService(customerRepo repository.CustomerRepository, userRepo repository.UserRepository, bedRepo repository.BedRepository, careRepo repository.CareLevelRepository, dietRepo repository.DietPlanRepository) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		userRepo:     userRepo,
		bedRepo:      bedRepo,
		careRepo:     careRepo,
		dietRepo:     dietRepo,
	}
}

// Create 创建客户
func (s *CustomerService) Create(ctx context.Context, req *domain.Customer) error {
	// 验证用户是否存在
	if !req.UserID.IsZero() {
		user, err := s.userRepo.FindById(ctx, req.UserID)
		if err != nil {
			return err
		}
		if user == nil {
			return errors.New("用户不存在")
		}
	}

	customer := &domain.Customer{
		UserID:    req.UserID,
		Name:      req.Name,
		Age:       req.Age,
		Gender:    req.Gender,
		Phone:     req.Phone,
		IDCard:    req.IDCard,
		Status:    "未入住",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.customerRepo.Create(ctx, customer)
}

// GetById 根据ID获取客户
func (s *CustomerService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Customer, error) {
	return s.customerRepo.FindById(ctx, id)
}

// GetList 获取客户列表
func (s *CustomerService) GetList(ctx context.Context, query domain.CustomerQuery, skip, limit int64) ([]*domain.Customer, int64, error) {
	filter := bson.M{}

	// 精确匹配
	if query.Status != "" {
		filter["status"] = query.Status
	}

	// ObjectID 转换处理
	idFields := map[string]string{
		"bed_id":            query.BedID,
		"care_level_id":     query.CareLevelID,
		"diet_plan_id":      query.DietPlanID,
		"health_manager_id": query.HealthManagerID,
	}

	for field, val := range idFields {
		if val != "" {
			if oid, err := primitive.ObjectIDFromHex(val); err == nil {
				filter[field] = oid
			}
		}
	}

	// 模糊查询
	if query.Name != "" {
		filter["name"] = primitive.Regex{Pattern: query.Name, Options: "i"}
	}
	if query.Phone != "" {
		filter["phone"] = primitive.Regex{Pattern: query.Phone, Options: "i"}
	}
	if query.IDCard != "" {
		filter["id_card"] = primitive.Regex{Pattern: query.IDCard, Options: "i"}
	}

	// 数值范围查询
	if query.MinAge > 0 || query.MaxAge > 0 {
		ageFilter := bson.M{}
		if query.MinAge > 0 {
			ageFilter["$gte"] = query.MinAge
		}
		if query.MaxAge > 0 {
			ageFilter["$lte"] = query.MaxAge
		}
		filter["age"] = ageFilter
	}

	// 复合搜索逻辑
	if query.SearchKey != "" {
		filter["$or"] = []bson.M{
			{"name": primitive.Regex{Pattern: query.SearchKey, Options: "i"}},
			{"phone": primitive.Regex{Pattern: query.SearchKey, Options: "i"}},
		}
	}

	return s.customerRepo.FindList(ctx, filter, skip, limit)
}

// Update 更新客户信息 (支持部分字段更新)
func (s *CustomerService) Update(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	// 验证客户是否存在
	customer, err := s.customerRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户不存在")
	}

	// 业务逻辑处理：检查是否修改了床位
	if val, hasBedID := updates["bed_id"]; hasBedID {
		var newBedID primitive.ObjectID
		var err error

		// 解析 bed_id，支持字符串映射或 ObjectID
		switch v := val.(type) {
		case string:
			if v != "" {
				newBedID, err = primitive.ObjectIDFromHex(v)
				if err != nil {
					return errors.New("无效的床位ID格式")
				}
			}
		case primitive.ObjectID:
			newBedID = v
		}

		// 如果床位发生了变更
		if customer.BedID != newBedID {
			// 1. 释放原有的旧床位
			if !customer.BedID.IsZero() {
				if err := s.bedRepo.Release(ctx, customer.BedID); err != nil {
					return err
				}
			}

			// 2. 分配并占用新床位
			if !newBedID.IsZero() {
				// 检查新床位是否存在且处于空闲状态
				bed, err := s.bedRepo.FindById(ctx, newBedID)
				if err != nil {
					return err
				}
				if bed == nil {
					return errors.New("床位不存在")
				}
				if bed.Status == "占用" {
					return errors.New("该床位已被占用")
				}

				// 占用该新床位
				if err := s.bedRepo.AssignToCustomer(ctx, newBedID, id); err != nil {
					return err
				}
			}
		}

		updates["bed_id"] = newBedID
	}

	return s.customerRepo.Update(ctx, id, updates)
}

// SetHealthManager 设置健康管家
func (s *CustomerService) SetHealthManager(ctx context.Context, customerID, managerID primitive.ObjectID, managerName string) error {
	customer, err := s.customerRepo.FindById(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户不存在")
	}

	updates := map[string]interface{}{
		"health_manager_id": managerID,
		"health_manager":    managerName,
	}

	return s.customerRepo.Update(ctx, customer.ID, updates)
}

// SetBed 设置床位
func (s *CustomerService) SetBed(ctx context.Context, customerID, bedID primitive.ObjectID) error {
	customer, err := s.customerRepo.FindById(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户不存在")
	}

	// 如果之前有床位，先释放
	if !customer.BedID.IsZero() {
		err = s.bedRepo.Release(ctx, customer.BedID)
		if err != nil {
			return err
		}
	}

	// 分配新床位
	err = s.bedRepo.AssignToCustomer(ctx, bedID, customerID)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"bed_id": bedID,
	}

	return s.customerRepo.Update(ctx, customerID, updates)
}

// SetDietPlan 设置膳食计划
func (s *CustomerService) SetDietPlan(ctx context.Context, customerID, dietPlanID primitive.ObjectID) error {
	customer, err := s.customerRepo.FindById(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户不存在")
	}
	_, err = s.dietRepo.FindById(ctx, dietPlanID)
	if err != nil {
		return errors.New("膳食计划不存在")
	}
	updates := map[string]interface{}{
		"diet_plan_id": dietPlanID,
	}
	return s.customerRepo.Update(ctx, customerID, updates)
}

// SetCareLevel 设置护理级别
func (s *CustomerService) SetCareLevel(ctx context.Context, customerID, careLevelID primitive.ObjectID) error {
	customer, err := s.customerRepo.FindById(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户不存在")
	}
	_, err = s.careRepo.FindById(ctx, careLevelID)
	if err != nil {
		return errors.New("护理级别不存在")
	}
	updates := map[string]interface{}{
		"care_level_id": careLevelID,
	}
	return s.customerRepo.Update(ctx, customerID, updates)
}

// Delete 删除客户
func (s *CustomerService) Delete(ctx context.Context, id primitive.ObjectID) error {
	customer, err := s.customerRepo.FindById(ctx, id)
	if err != nil {
		return errors.New("患者不存在")
	}
	if customer.Status != "退住" {
		return errors.New("请先退住，无法删除")
	}
	return s.customerRepo.Delete(ctx, id)
}
