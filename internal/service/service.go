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

type ServiceService struct {
	serviceRepo        repository.ServiceRepository
	customerServiceRepo repository.CustomerServiceRepository
	customerRepo       repository.CustomerRepository
}

func NewServiceService(serviceRepo repository.ServiceRepository, customerServiceRepo repository.CustomerServiceRepository, customerRepo repository.CustomerRepository) *ServiceService {
	return &ServiceService{
		serviceRepo:         serviceRepo,
		customerServiceRepo: customerServiceRepo,
		customerRepo:        customerRepo,
	}
}

// CreateService 创建服务项目
func (s *ServiceService) CreateService(ctx context.Context, req *domain.Service) error {
	service := &domain.Service{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Unit:        req.Unit,
		Status:      "启用",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return s.serviceRepo.Create(ctx, service)
}

// GetServiceById 根据ID获取服务项目
func (s *ServiceService) GetServiceById(ctx context.Context, id primitive.ObjectID) (*domain.Service, error) {
	return s.serviceRepo.FindById(ctx, id)
}

// GetServiceList 获取服务项目列表
func (s *ServiceService) GetServiceList(ctx context.Context, category, status string, skip, limit int64) ([]*domain.Service, int64, error) {
	filter := bson.M{}
	if category != "" {
		filter["category"] = category
	}
	if status != "" {
		filter["status"] = status
	}
	return s.serviceRepo.FindList(ctx, filter, skip, limit)
}

// UpdateService 更新服务项目
func (s *ServiceService) UpdateService(ctx context.Context, id primitive.ObjectID, req *domain.Service) error {
	service, err := s.serviceRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if service == nil {
		return errors.New("服务项目不存在")
	}

	service.Name = req.Name
	service.Description = req.Description
	service.Category = req.Category
	service.Price = req.Price
	service.Unit = req.Unit
	service.Status = req.Status
	service.UpdatedAt = time.Now()

	return s.serviceRepo.Update(ctx, service)
}

// DeleteService 删除服务项目
func (s *ServiceService) DeleteService(ctx context.Context, id primitive.ObjectID) error {
	return s.serviceRepo.Delete(ctx, id)
}

// PurchaseService 客户购买服务
func (s *ServiceService) PurchaseService(ctx context.Context, customerID, serviceID primitive.ObjectID, startDate time.Time) error {
	// 验证客户是否存在
	customer, err := s.customerRepo.FindById(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户不存在")
	}

	// 验证服务是否存在
	service, err := s.serviceRepo.FindById(ctx, serviceID)
	if err != nil {
		return err
	}
	if service == nil {
		return errors.New("服务项目不存在")
	}

	cs := &domain.CustomerService{
		CustomerID:  customerID,
		ServiceID:   serviceID,
		ServiceName: service.Name,
		StartDate:   startDate,
		Status:      "进行中",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.customerServiceRepo.Create(ctx, cs)
}

// GetCustomerServices 获取客户购买的服务列表
func (s *ServiceService) GetCustomerServices(ctx context.Context, customerID primitive.ObjectID) ([]*domain.CustomerService, error) {
	return s.customerServiceRepo.FindByCustomerID(ctx, customerID)
}

// EndService 结束客户服务
func (s *ServiceService) EndService(ctx context.Context, customerServiceID primitive.ObjectID, endDate time.Time) error {
	cs, err := s.customerServiceRepo.FindById(ctx, customerServiceID)
	if err != nil {
		return err
	}
	if cs == nil {
		return errors.New("客户服务记录不存在")
	}

	cs.EndDate = endDate
	cs.Status = "已结束"
	cs.UpdatedAt = time.Now()

	return s.customerServiceRepo.Update(ctx, cs)
}

