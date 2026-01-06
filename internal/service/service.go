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

type ServerService struct {
	serviceRepo         repository.ServiceRepository
	customerServiceRepo repository.CustomerServiceRepository
	customerRepo        repository.CustomerRepository
	bedRepo             repository.BedRepository
	roomRepo            repository.RoomRepository
	careLevelRepo       repository.CareLevelRepository
}

func NewServiceService(
	serviceRepo repository.ServiceRepository,
	customerServiceRepo repository.CustomerServiceRepository,
	customerRepo repository.CustomerRepository,
	bedRepo repository.BedRepository,
	roomRepo repository.RoomRepository,
	careLevelRepo repository.CareLevelRepository,
) *ServerService {
	return &ServerService{
		serviceRepo:         serviceRepo,
		customerServiceRepo: customerServiceRepo,
		customerRepo:        customerRepo,
		bedRepo:             bedRepo,
		roomRepo:            roomRepo,
		careLevelRepo:       careLevelRepo,
	}
}

// GetCustomerServiceByGroup 获取客户服务列表（包含用户名称和服务名称）
// 分页基于用户数量，skip=1 表示跳过1个用户，limit=10 表示返回10个用户的所有服务记录
func (s *ServerService) GetCustomerServiceByGroup(ctx context.Context, status string, skip, limit int64) ([]*domain.CustomerServiceResponse, int64, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	// 先获取所有符合条件的唯一用户ID列表（用于分页）
	allCustomerIDs, err := s.customerServiceRepo.FindDistinctCustomerIDs(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 用户总数
	totalCustomers := int64(len(allCustomerIDs))
	if totalCustomers == 0 {
		return []*domain.CustomerServiceResponse{}, 0, nil
	}

	// 对用户ID列表进行分页
	startIdx := skip
	if startIdx < 0 {
		startIdx = 0
	}
	if startIdx >= totalCustomers {
		return []*domain.CustomerServiceResponse{}, totalCustomers, nil
	}

	endIdx := startIdx + limit
	if limit <= 0 || endIdx > totalCustomers {
		endIdx = totalCustomers
	}

	// 获取当前页的用户ID
	pagedCustomerIDs := allCustomerIDs[startIdx:endIdx]

	// 查询这些用户的所有服务记录
	customerIDFilter := bson.M{"customer_id": bson.M{"$in": pagedCustomerIDs}}
	if status != "" {
		customerIDFilter["status"] = status
	}
	list, _, err := s.customerServiceRepo.FindList(ctx, customerIDFilter, 0, 0)
	if err != nil {
		return nil, 0, err
	}

	// 收集所有需要查询的 ServiceID
	serviceIDs := make(map[primitive.ObjectID]bool)
	for _, cs := range list {
		serviceIDs[cs.ServiceID] = true
	}

	// 批量查询客户信息（包含完整信息以获取 BedID 和 CareLevelID）
	customerMap := make(map[primitive.ObjectID]*domain.Customer)
	bedIDs := make(map[primitive.ObjectID]bool)
	careLevelIDs := make(map[primitive.ObjectID]bool)
	for _, id := range pagedCustomerIDs {
		customer, err := s.customerRepo.FindById(ctx, id)
		if err == nil && customer != nil {
			customerMap[id] = customer
			if !customer.BedID.IsZero() {
				bedIDs[customer.BedID] = true
			}
			if !customer.CareLevelID.IsZero() {
				careLevelIDs[customer.CareLevelID] = true
			}
		}
	}

	// 批量查询服务信息
	serviceMap := make(map[primitive.ObjectID]*domain.Service)
	for id := range serviceIDs {
		service, err := s.serviceRepo.FindById(ctx, id)
		if err == nil && service != nil {
			serviceMap[id] = service
		}
	}

	// 批量查询床位信息
	bedMap := make(map[primitive.ObjectID]*domain.Bed)
	for id := range bedIDs {
		bed, err := s.bedRepo.FindById(ctx, id)
		if err == nil && bed != nil {
			bedMap[id] = bed
		}
	}

	// 批量查询护理级别信息
	careLevelMap := make(map[primitive.ObjectID]*domain.CareLevel)
	for id := range careLevelIDs {
		careLevel, err := s.careLevelRepo.FindById(ctx, id)
		if err == nil && careLevel != nil {
			careLevelMap[id] = careLevel
		}
	}

	// 组装响应数据
	result := make([]*domain.CustomerServiceResponse, 0, len(list))
	for _, cs := range list {
		resp := &domain.CustomerServiceResponse{
			ID:          cs.ID,
			CustomerID:  cs.CustomerID,
			ServiceID:   cs.ServiceID,
			ServiceName: cs.ServiceName,
			StartDate:   cs.StartDate,
			EndDate:     cs.EndDate,
			Status:      cs.Status,
			CreatedAt:   cs.CreatedAt,
			UpdatedAt:   cs.UpdatedAt,
		}

		// 填充客户相关信息（姓名、床位号、护理级别）
		if customer, ok := customerMap[cs.CustomerID]; ok {
			resp.CustomerName = customer.Name
			// 填充床位号
			if bed, bedOk := bedMap[customer.BedID]; bedOk {
				resp.BedNumber = bed.Number
			}
			// 填充护理级别
			if careLevel, clOk := careLevelMap[customer.CareLevelID]; clOk {
				resp.CareLevelstr = careLevel.Name
			}
		}

		// 填充服务详情
		if svc, ok := serviceMap[cs.ServiceID]; ok {
			resp.ServiceName = svc.Name
			resp.ServiceDesc = svc.Description
			resp.Category = svc.Category
			resp.Price = svc.Price
			resp.Unit = svc.Unit
		}
		result = append(result, resp)
	}

	return result, totalCustomers, nil
}

// CreateService 创建服务项目
func (s *ServerService) CreateService(ctx context.Context, req *domain.Service) error {
	service := &domain.Service{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Unit:        req.Unit,
		Status:      "启用",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.serviceRepo.Create(ctx, service)
}

// GetServiceById 根据ID获取服务项目
func (s *ServerService) GetServiceById(ctx context.Context, id primitive.ObjectID) (*domain.Service, error) {
	return s.serviceRepo.FindById(ctx, id)
}

// GetServiceList 获取服务项目列表
func (s *ServerService) GetServiceList(ctx context.Context, category, status string, skip, limit int64) ([]*domain.Service, int64, error) {
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
func (s *ServerService) UpdateService(ctx context.Context, id primitive.ObjectID, req *domain.Service) error {
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
func (s *ServerService) DeleteService(ctx context.Context, id primitive.ObjectID) error {
	return s.serviceRepo.Delete(ctx, id)
}

// TODO
// PurchaseService 客户购买服务
func (s *ServerService) PurchaseService(ctx context.Context, customerID, serviceID primitive.ObjectID, startDate time.Time) error {
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
func (s *ServerService) GetCustomerServices(ctx context.Context, customerID primitive.ObjectID) ([]*domain.CustomerService, error) {
	return s.customerServiceRepo.FindByCustomerID(ctx, customerID)
}

// GetCustomerServiceList 获取客户服务列表（支持分页和筛选，包含关联信息）
func (s *ServerService) GetCustomerServiceList(ctx context.Context, customerID, serviceID primitive.ObjectID, status string, skip, limit int64) ([]*domain.CustomerServiceResponse, int64, error) {
	filter := bson.M{}
	if !customerID.IsZero() {
		filter["customer_id"] = customerID
	}
	if !serviceID.IsZero() {
		filter["service_id"] = serviceID
	}
	if status != "" {
		filter["status"] = status
	}

	list, total, err := s.customerServiceRepo.FindList(ctx, filter, skip, limit)
	if err != nil {
		return nil, 0, err
	}

	// 收集所有需要查询的 CustomerID 和 ServiceID
	customerIDs := make(map[primitive.ObjectID]bool)
	serviceIDs := make(map[primitive.ObjectID]bool)
	for _, cs := range list {
		customerIDs[cs.CustomerID] = true
		serviceIDs[cs.ServiceID] = true
	}

	// 批量查询客户信息（包含完整信息以获取 BedID 和 CareLevelID）
	customerMap := make(map[primitive.ObjectID]*domain.Customer)
	bedIDs := make(map[primitive.ObjectID]bool)
	careLevelIDs := make(map[primitive.ObjectID]bool)
	for id := range customerIDs {
		customer, err := s.customerRepo.FindById(ctx, id)
		if err == nil && customer != nil {
			customerMap[id] = customer
			if !customer.BedID.IsZero() {
				bedIDs[customer.BedID] = true
			}
			if !customer.CareLevelID.IsZero() {
				careLevelIDs[customer.CareLevelID] = true
			}
		}
	}

	// 批量查询服务信息
	serviceMap := make(map[primitive.ObjectID]*domain.Service)
	for id := range serviceIDs {
		service, err := s.serviceRepo.FindById(ctx, id)
		if err == nil && service != nil {
			serviceMap[id] = service
		}
	}

	// 批量查询床位信息
	bedMap := make(map[primitive.ObjectID]*domain.Bed)
	for id := range bedIDs {
		bed, err := s.bedRepo.FindById(ctx, id)
		if err == nil && bed != nil {
			bedMap[id] = bed
		}
	}

	// 批量查询护理级别信息
	careLevelMap := make(map[primitive.ObjectID]*domain.CareLevel)
	for id := range careLevelIDs {
		careLevel, err := s.careLevelRepo.FindById(ctx, id)
		if err == nil && careLevel != nil {
			careLevelMap[id] = careLevel
		}
	}

	// 组装响应数据
	result := make([]*domain.CustomerServiceResponse, 0, len(list))
	for _, cs := range list {
		resp := &domain.CustomerServiceResponse{
			ID:          cs.ID,
			CustomerID:  cs.CustomerID,
			ServiceID:   cs.ServiceID,
			ServiceName: cs.ServiceName,
			StartDate:   cs.StartDate,
			EndDate:     cs.EndDate,
			Status:      cs.Status,
			CreatedAt:   cs.CreatedAt,
			UpdatedAt:   cs.UpdatedAt,
		}

		// 填充客户相关信息（姓名、床位号、护理级别）
		if customer, ok := customerMap[cs.CustomerID]; ok {
			resp.CustomerName = customer.Name
			// 填充床位号
			if bed, bedOk := bedMap[customer.BedID]; bedOk {
				resp.BedNumber = bed.Number
			}
			// 填充护理级别
			if careLevel, clOk := careLevelMap[customer.CareLevelID]; clOk {
				resp.CareLevelstr = careLevel.Name
			}
		}

		// 填充服务详情
		if svc, ok := serviceMap[cs.ServiceID]; ok {
			resp.ServiceName = svc.Name
			resp.ServiceDesc = svc.Description
			resp.Category = svc.Category
			resp.Price = svc.Price
			resp.Unit = svc.Unit
		}
		result = append(result, resp)
	}

	return result, total, nil
}

// EndService 结束客户服务
func (s *ServerService) EndService(ctx context.Context, customerServiceID primitive.ObjectID, endDate time.Time) error {
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

// UpdateCustomerServiceEndDate 修改客户服务结束时间（不改变状态）
func (s *ServerService) UpdateCustomerServiceEndDate(ctx context.Context, customerServiceID primitive.ObjectID, endDate time.Time) error {
	cs, err := s.customerServiceRepo.FindById(ctx, customerServiceID)
	if err != nil {
		return err
	}
	if cs == nil {
		return errors.New("客户服务记录不存在")
	}

	cs.EndDate = endDate
	cs.UpdatedAt = time.Now()

	return s.customerServiceRepo.Update(ctx, cs)
}

// CancelCustomerService 取消客户单一服务
func (s *ServerService) CancelCustomerService(ctx context.Context, customerServiceID primitive.ObjectID) error {
	cs, err := s.customerServiceRepo.FindById(ctx, customerServiceID)
	if err != nil {
		return err
	}
	if cs == nil {
		return errors.New("客户服务记录不存在")
	}

	cs.Status = "已取消"
	cs.UpdatedAt = time.Now()

	return s.customerServiceRepo.Update(ctx, cs)
}
