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
}

func NewCustomerService(customerRepo repository.CustomerRepository, userRepo repository.UserRepository, bedRepo repository.BedRepository) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		userRepo:     userRepo,
		bedRepo:      bedRepo,
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
		UserID:          req.UserID,
		Name:            req.Name,
		Age:             req.Age,
		Gender:          req.Gender,
		Phone:           req.Phone,
		IDCard:          req.IDCard,
		Status:          "未入住",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	return s.customerRepo.Create(ctx, customer)
}

// GetById 根据ID获取客户
func (s *CustomerService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Customer, error) {
	return s.customerRepo.FindById(ctx, id)
}

// GetList 获取客户列表
func (s *CustomerService) GetList(ctx context.Context, status string, skip, limit int64) ([]*domain.Customer, int64, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	return s.customerRepo.FindList(ctx, filter, skip, limit)
}

// Update 更新客户信息
func (s *CustomerService) Update(ctx context.Context, id primitive.ObjectID, req *domain.Customer) error {
	customer, err := s.customerRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户不存在")
	}

	customer.Name = req.Name
	customer.Age = req.Age
	customer.Gender = req.Gender
	customer.Phone = req.Phone
	customer.IDCard = req.IDCard
	customer.UpdatedAt = time.Now()

	return s.customerRepo.Update(ctx, customer)
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

	customer.HealthManagerID = managerID
	customer.HealthManager = managerName
	customer.UpdatedAt = time.Now()

	return s.customerRepo.Update(ctx, customer)
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

	customer.BedID = bedID
	customer.UpdatedAt = time.Now()

	return s.customerRepo.Update(ctx, customer)
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

	customer.DietPlanID = dietPlanID
	customer.UpdatedAt = time.Now()

	return s.customerRepo.Update(ctx, customer)
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

	customer.CareLevelID = careLevelID
	customer.UpdatedAt = time.Now()

	return s.customerRepo.Update(ctx, customer)
}

// Delete 删除客户
func (s *CustomerService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.customerRepo.Delete(ctx, id)
}

