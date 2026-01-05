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

// GetNameAndID 获取客户名称和ID列表
func (s *CustomerService) GetListNameAndID(ctx context.Context) ([]*domain.CustomerNameID, error) {
	// 获取客户列表
	customers, total, err := s.customerRepo.FindList(ctx, bson.M{}, 0, 0) // 这里假设我们想要获取所有客户，所以skip和limit都设置为0
	if err != nil {
		return nil, err
	}

	// 构建名称和ID的切片
	nameIDList := make([]*domain.CustomerNameID, 0, total)
	for _, customer := range customers {
		nameID := &domain.CustomerNameID{
			ID:   customer.ID,
			Name: customer.Name,
		}
		nameIDList = append(nameIDList, nameID)
	}

	return nameIDList, nil
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
func (s *CustomerService) GetList(ctx context.Context, query domain.CustomerQuery, skip, limit int64) ([]*domain.CustomerResponse, int64, error) {
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

	// 1. 查询客户列表
	list, total, err := s.customerRepo.FindList(ctx, filter, skip, limit)
	if err != nil {
		return nil, 0, err
	}

	if len(list) == 0 {
		return []*domain.CustomerResponse{}, total, nil
	}

	// 2. 收集所有关联 ID（去重）
	bedIDs := make([]primitive.ObjectID, 0)
	careIDs := make([]primitive.ObjectID, 0)
	dietIDs := make([]primitive.ObjectID, 0)

	for _, v := range list {
		if !v.BedID.IsZero() {
			bedIDs = append(bedIDs, v.BedID)
		}
		if !v.CareLevelID.IsZero() {
			careIDs = append(careIDs, v.CareLevelID)
		}
		if !v.DietPlanID.IsZero() {
			dietIDs = append(dietIDs, v.DietPlanID)
		}
	}

	// 3. 批量查询关联表并构建 Map
	bedMap := s.getBedMap(ctx, bedIDs)
	careMap := s.getCareMap(ctx, careIDs)
	dietMap := s.getDietMap(ctx, dietIDs)

	// 4. 组装结果
	listRes := make([]*domain.CustomerResponse, 0, len(list))
	for _, v := range list {
		res := &domain.CustomerResponse{
			ID:              v.ID,
			UserID:          v.UserID,
			Name:            v.Name,
			Age:             v.Age,
			Gender:          v.Gender,
			Phone:           v.Phone,
			IDCard:          v.IDCard,
			BedID:           v.BedID,
			Bed:             bedMap[v.BedID],
			DietPlanID:      v.DietPlanID,
			DietPlan:        dietMap[v.DietPlanID],
			CareLevelID:     v.CareLevelID,
			CareLevel:       careMap[v.CareLevelID],
			HealthManagerID: v.HealthManagerID,
			HealthManager:   v.HealthManager,
			Status:          v.Status,
			CheckInDate:     v.CheckInDate,
			CheckOutDate:    v.CheckOutDate,
			HealthLevel:     v.HealthLevel,
			MedicalHistory:  v.MedicalHistory,
			Medication:      v.Medication,
			AllergyHistory:  v.AllergyHistory,
			ContactName:     v.ContactName,
			Relationship:    v.Relationship,
			ContactPhone:    v.ContactPhone,
			ContactAddress:  v.ContactAddress,
			Remarks:         v.Remarks,
			CreatedAt:       v.CreatedAt,
			UpdatedAt:       v.UpdatedAt,
		}
		listRes = append(listRes, res)
	}

	return listRes, total, nil
}

// getBedMap 批量查询床位信息并构建 ID->Name 映射
func (s *CustomerService) getBedMap(ctx context.Context, ids []primitive.ObjectID) map[primitive.ObjectID]string {
	result := make(map[primitive.ObjectID]string)
	if len(ids) == 0 {
		return result
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	beds, _, err := s.bedRepo.FindList(ctx, filter, 0, int64(len(ids)))
	if err != nil {
		return result
	}

	for _, bed := range beds {
		result[bed.ID] = bed.Number
	}
	return result
}

// getCareMap 批量查询护理级别信息并构建 ID->Name 映射
func (s *CustomerService) getCareMap(ctx context.Context, ids []primitive.ObjectID) map[primitive.ObjectID]string {
	result := make(map[primitive.ObjectID]string)
	if len(ids) == 0 {
		return result
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	careLevels, _, err := s.careRepo.FindList(ctx, filter, 0, int64(len(ids)))
	if err != nil {
		return result
	}

	for _, care := range careLevels {
		result[care.ID] = care.Name
	}
	return result
}

// getDietMap 批量查询膳食计划信息并构建 ID->Name 映射
func (s *CustomerService) getDietMap(ctx context.Context, ids []primitive.ObjectID) map[primitive.ObjectID]string {
	result := make(map[primitive.ObjectID]string)
	if len(ids) == 0 {
		return result
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	dietPlans, _, err := s.dietRepo.FindList(ctx, filter, 0, int64(len(ids)))
	if err != nil {
		return result
	}

	for _, diet := range dietPlans {
		result[diet.ID] = diet.Name
	}
	return result
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
					//继续占用原来的床位
					if err := s.bedRepo.AssignToCustomer(ctx, customer.BedID, customer.ID); err != nil {
						return err
					}
					return errors.New("该床位已被占用")
				}

				// 占用该新床位
				if err := s.bedRepo.AssignToCustomer(ctx, newBedID, customer.ID); err != nil {
					return err
				}

			}

		}

		updates["bed_id"] = newBedID
	}

	// 业务逻辑处理：检查是否修改了护理级别

	if val, hasCareLevelID := updates["care_level_id"]; hasCareLevelID {
		var newCareLevelID primitive.ObjectID
		var err error

		// 解析 care_level_id，支持字符串映射或 ObjectID
		switch v := val.(type) {
		case string:
			if v != "" {
				newCareLevelID, err = primitive.ObjectIDFromHex(v)
				if err != nil {
					return errors.New("无效的护理级别ID格式")
				}
			}
		case primitive.ObjectID:
			newCareLevelID = v
		}

		updates["care_level_id"] = newCareLevelID
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
