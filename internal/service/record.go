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

type RecordService struct {
	recordRepo    repository.RecordRepository
	customerRepo  repository.CustomerRepository
	bedRepo       repository.BedRepository
	careLevelRepo repository.CareLevelRepository
}

func NewRecordService(recordRepo repository.RecordRepository, customerRepo repository.CustomerRepository, bedRepo repository.BedRepository, careLevelRepo repository.CareLevelRepository) *RecordService {
	return &RecordService{
		recordRepo:    recordRepo,
		customerRepo:  customerRepo,
		bedRepo:       bedRepo,
		careLevelRepo: careLevelRepo,
	}
}

// CheckIn 入住登记
func (s *RecordService) CheckIn(ctx context.Context, req *domain.ElderlyRegisterRequest) error {
	// 1. 验证床位
	bedID, err := primitive.ObjectIDFromHex(req.BedID)
	if err != nil {
		return errors.New("无效的床位ID")
	}
	// 这里的床位检查逻辑
	bed, err := s.bedRepo.FindById(ctx, bedID)
	if err != nil {
		return err
	}
	if bed == nil {
		return errors.New("床位不存在")
	}
	if bed.Status == "占用" {
		return errors.New("床位已被占用")
	}

	// 2. 解析其他 ID
	// 假设 NursingLevel 和 DietaryType 传的是 ID 字符串
	// 如果是名称，需要额外的查询逻辑。这里暂且尝试解析为 ID，如果失败则作为 Name 存储（如果 Customer 结构支持）
	// Customer 结构中有 ID 字段吗？CareLevelID, DietPlanID.
	var careLevelID, dietPlanID primitive.ObjectID
	if req.NursingLevel != "" {
		if id, err := primitive.ObjectIDFromHex(req.NursingLevel); err == nil {
			careLevelID = id
		}
	}
	if req.DietaryType != "" {
		if id, err := primitive.ObjectIDFromHex(req.DietaryType); err == nil {
			dietPlanID = id
		}
	}

	// 3. 解析日期
	checkInDate := time.Now()
	if req.CheckInDate != "" {
		// 尝试解析 "2006-01-02"
		if t, err := time.Parse("2006-01-02", req.CheckInDate); err == nil {
			checkInDate = t
		}
	}

	// 4. 创建客户
	customer := &domain.Customer{
		Name:           req.Name,
		Age:            req.Age,
		Gender:         req.Gender,
		Phone:          req.PhoneNumber,
		IDCard:         req.IDCard,
		BedID:          bedID,
		CareLevelID:    careLevelID,
		DietPlanID:     dietPlanID,
		Status:         "入住中",
		CheckInDate:    checkInDate,
		HealthLevel:    req.HealthLevel,
		MedicalHistory: req.MedicalHistory,
		Medication:     req.Medication,
		AllergyHistory: req.AllergyHistory,
		ContactName:    req.ContactName,
		Relationship:   req.Relationship,
		ContactPhone:   req.ContactPhone,
		ContactAddress: req.ContactAddress,
		Remarks:        req.Remarks,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = s.customerRepo.Create(ctx, customer)
	if err != nil {
		return err
	}

	// 5. 分配床位
	err = s.bedRepo.AssignToCustomer(ctx, bedID, customer.ID)
	if err != nil {
		// 回滚？这里暂不处理复杂回滚，生产环境需要事务
		return err
	}

	// 6. 创建入住记录
	record := &domain.Record{
		CustomerID: customer.ID,
		Type:       "入住",
		StartTime:  checkInDate,
		Note:       req.Remarks, // 使用备注作为记录 Note
		CreatedBy:  "System",    // 或从 Context 获取? 上下文没传 CreatedBy. 暂时写 System 或 "Admin"
		CreatedAt:  time.Now(),
	}

	return s.recordRepo.Create(ctx, record)
}

// CheckOut 退住登记
func (s *RecordService) CheckOut(ctx context.Context, customerID primitive.ObjectID, note, createdBy string) error {
	// 验证客户是否存在
	customer, err := s.customerRepo.FindById(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户不存在")
	}

	// 释放床位
	if !customer.BedID.IsZero() {
		err = s.bedRepo.Release(ctx, customer.BedID)
		if err != nil {
			return err
		}
		customer.BedID = primitive.NilObjectID
	}

	// 更新客户状态
	var updates = map[string]interface{}{
		"status":         "退住",
		"updated_at":     time.Now(),
		"check_out_date": time.Now(),
	}
	err = s.customerRepo.Update(ctx, customerID, updates)
	if err != nil {
		return err
	}

	// 创建退住记录
	record := &domain.Record{
		CustomerID: customerID,
		Type:       "退住",
		StartTime:  customer.CheckInDate,
		EndTime:    time.Now(),
		Note:       note,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
	}

	return s.recordRepo.Create(ctx, record)
}

// Outgoing 外出登记
func (s *RecordService) Outgoing(ctx context.Context, customerID primitive.ObjectID, note, createdBy string) error {
	// 验证客户是否存在
	customer, err := s.customerRepo.FindById(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户不存在")
	}
	if customer.Status != "入住中" {
		return errors.New("只有入住中的客户才能外出")
	}

	// 更新客户状态
	var updates = map[string]interface{}{
		"status":       "外出中",
		"updated_at":   time.Now(),
		"check_out_at": time.Now(),
	}
	err = s.customerRepo.Update(ctx, customerID, updates)
	if err != nil {
		return err
	}

	// 创建外出记录
	record := &domain.Record{
		CustomerID: customerID,
		Type:       "外出",
		StartTime:  time.Now(),
		Note:       note,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
	}

	return s.recordRepo.Create(ctx, record)
}

// Return 外出返回
func (s *RecordService) Return(ctx context.Context, customerID primitive.ObjectID, note, createdBy string) error {
	// 验证客户是否存在
	customer, err := s.customerRepo.FindById(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户不存在")
	}
	if customer.Status != "外出中" {
		return errors.New("客户当前状态不是外出中")
	}

	// 更新客户状态
	var updates = map[string]interface{}{
		"status":       "入住中",
		"updated_at":   time.Now(),
		"check_out_at": time.Time{},
	}
	err = s.customerRepo.Update(ctx, customerID, updates)
	if err != nil {
		return err
	}

	// 查找最近的外出记录并更新
	records, err := s.recordRepo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return err
	}

	// 找到最近的外出记录
	for i := len(records) - 1; i >= 0; i-- {
		if records[i].Type == "外出" && records[i].EndTime.IsZero() {
			records[i].EndTime = time.Now()
			records[i].Note = note
			return s.recordRepo.Update(ctx, records[i])
		}
	}

	return errors.New("未找到外出记录")
}

// GetList 获取登记记录列表
func (s *RecordService) GetList(ctx context.Context, customerID primitive.ObjectID, recordType string, skip, limit int64) ([]*domain.Record, int64, error) {
	filter := bson.M{}
	if !customerID.IsZero() {
		filter["customer_id"] = customerID
	}
	if recordType != "" {
		filter["type"] = recordType
	}
	return s.recordRepo.FindList(ctx, filter, skip, limit)
}

// GetByCustomerID 获取客户的登记记录
func (s *RecordService) GetByCustomerID(ctx context.Context, customerID primitive.ObjectID) ([]*domain.Record, error) {
	return s.recordRepo.FindByCustomerID(ctx, customerID)
}

// GetCheckInList 获取入住信息列表
func (s *RecordService) GetCheckInList(ctx context.Context, name, roomNumber, nursingLevel string, startDate, endDate string, skip, limit int64) ([]map[string]interface{}, error) {
	filter := bson.M{}

	// 名字模糊查询
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	// 房间号查询逻辑：Customer -> Bed -> Room
	// 如果提供了房间号，需要先找到该房间下的所有床位ID，然后 Filter customer.bed_id IN [...]
	if roomNumber != "" {

		bedFilter := bson.M{"number": bson.M{"$regex": "^" + roomNumber}} // 假设床位号以房间号开头
		beds, _, err := s.bedRepo.FindList(ctx, bedFilter, 0, 0)
		if err != nil {
			return nil, err
		}
		var bedIDs []primitive.ObjectID
		for _, b := range beds {
			bedIDs = append(bedIDs, b.ID)
		}
		if len(bedIDs) > 0 {
			filter["bed_id"] = bson.M{"$in": bedIDs}
		} else {
			return []map[string]interface{}{}, nil
		}
	}

	if nursingLevel != "" {
		if id, err := primitive.ObjectIDFromHex(nursingLevel); err == nil {
			filter["care_level_id"] = id
		}
	}

	// 时间范围
	if startDate != "" || endDate != "" {
		dateFilter := bson.M{}
		if startDate != "" {
			if t, err := time.Parse("2006-01-02", startDate); err == nil {
				dateFilter["$gte"] = t
			}
		}
		if endDate != "" {
			if t, err := time.Parse("2006-01-02", endDate); err == nil {
				// 结束日期通常包含当天，加一天或处理时间
				t = t.Add(24 * time.Hour)
				dateFilter["$lt"] = t
			}
		}
		filter["check_in_date"] = dateFilter
	}

	// 查询客户
	customers, _, err := s.customerRepo.FindList(ctx, filter, skip, limit) // 暂未分页，全量返回
	if err != nil {
		return nil, err
	}

	// 组装结果
	var results []map[string]interface{}
	for _, c := range customers {
		item := map[string]interface{}{
			"id":            c.ID.Hex(),
			"name":          c.Name,
			"gender":        c.Gender,
			"age":           c.Age,
			"check_in_time": c.CheckInDate,
			"status":        c.Status,
		}

		// 填充床位号
		if !c.BedID.IsZero() {
			bed, _ := s.bedRepo.FindById(ctx, c.BedID)
			if bed != nil {
				item["bed_number"] = bed.Number
			} else {
				item["bed_number"] = "未知"
			}
		}

		// 填充护理级别
		// 注意: Customer 仅存储 CareLevelID，若需显示名称需关联查询
		//这里进行关联查询
		if c.CareLevelID.IsZero() {
			item["nursing_level"] = "未知"
		} else {
			careLevel, _ := s.careLevelRepo.FindById(ctx, c.CareLevelID)
			if careLevel != nil {
				item["nursing_level"] = careLevel.Name
			} else {
				item["nursing_level"] = "未知"
			}
		}

		results = append(results, item)
	}

	return results, nil
}

// GetCheckOutList 获取退住信息列表
func (s *RecordService) GetCheckOutList(ctx context.Context, name, bedId, reason, startDate, endDate string, skip, limit int64) ([]map[string]interface{}, error) {
	filter := bson.M{
		"type": "退住",
	}

	// 名字筛选 - 需先查客户
	if name != "" {
		customers, _, err := s.customerRepo.FindList(ctx, bson.M{"name": bson.M{"$regex": name, "$options": "i"}}, 0, 0)
		if err != nil {
			return nil, err
		}
		var customerIDs []primitive.ObjectID
		for _, c := range customers {
			customerIDs = append(customerIDs, c.ID)
		}
		if len(customerIDs) > 0 {
			filter["customer_id"] = bson.M{"$in": customerIDs}
		} else {
			return []map[string]interface{}{}, nil
		}
	}

	// 退住时间范围
	if startDate != "" || endDate != "" {
		dateFilter := bson.M{}
		if startDate != "" {
			if t, err := time.Parse("2006-01-02", startDate); err == nil {
				dateFilter["$gte"] = t
			}
		}
		if endDate != "" {
			if t, err := time.Parse("2006-01-02", endDate); err == nil {
				t = t.Add(24 * time.Hour)
				dateFilter["$lt"] = t
			}
		}
		// 记录的 EndTime 是退住时间?
		// CheckOut 逻辑: record.EndTime = time.Now(). record.StartTime = customer.CheckInDate.
		// 所以查询 "退住时间" 应该是 EndTime.
		filter["end_time"] = dateFilter
	}

	// 原因筛选(模糊)
	if reason != "" {
		filter["note"] = bson.M{"$regex": reason, "$options": "i"}
	}

	// 如果有bedId,先根据bedId查询出所有相关的customer_id
	// 因为Record表中没有bed_id字段,需要通过Customer表关联查询
	if bedId != "" {
		hex, err := primitive.ObjectIDFromHex(bedId)
		if err != nil {
			return nil, err
		}

		// 查询使用该床位的所有客户(包括历史客户)
		customers, _, err := s.customerRepo.FindList(ctx, bson.M{"bed_id": hex}, skip, limit)
		if err != nil {
			return nil, err
		}

		if len(customers) > 0 {
			var customerIDs []primitive.ObjectID
			for _, c := range customers {
				customerIDs = append(customerIDs, c.ID)
			}
			// 如果已经有customer_id筛选条件,需要取交集
			if existingFilter, ok := filter["customer_id"].(bson.M); ok {
				// 已有customer_id筛选,需要取交集
				if existingIDs, ok := existingFilter["$in"].([]primitive.ObjectID); ok {
					// 计算交集
					intersection := make([]primitive.ObjectID, 0)
					for _, id := range customerIDs {
						for _, existingID := range existingIDs {
							if id == existingID {
								intersection = append(intersection, id)
								break
							}
						}
					}
					customerIDs = intersection
				}
			}

			if len(customerIDs) > 0 {
				filter["customer_id"] = bson.M{"$in": customerIDs}
			} else {
				// 没有匹配的客户,返回空结果
				return []map[string]interface{}{}, nil
			}
		} else {
			// 该床位没有客户,返回空结果
			return []map[string]interface{}{}, nil
		}
	}
	records, _, err := s.recordRepo.FindList(ctx, filter, 0, 0)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, r := range records {
		customer, _ := s.customerRepo.FindById(ctx, r.CustomerID)
		customerName := "未知"
		careLevel := ""
		if customer != nil {
			customerName = customer.Name
			careLevel = customer.CareLevelID.Hex() // 同样只有ID
		}
		careName, err := s.careLevelRepo.FindById(ctx, customer.CareLevelID)
		if err != nil || careName == nil {
			return nil, err
		}

		days := int(r.EndTime.Sub(r.StartTime).Hours() / 24)

		item := map[string]interface{}{
			"id":              r.ID.Hex(),
			"customer_name":   customerName,
			"room_number":     "-", // 记录中未存，且客户已退住，难以获取历史床位
			"check_in_date":   r.StartTime,
			"check_out_date":  r.EndTime,
			"days":            days,
			"care_level":      careLevel,
			"care_level_name": careName.Name,
			"reason":          r.Note,
		}
		results = append(results, item)
	}
	return results, nil
}

// GetOutgoingList 获取外出登记信息列表
func (s *RecordService) GetOutgoingList(ctx context.Context, name, startDate, endDate, status string) ([]map[string]interface{}, error) {
	filter := bson.M{
		"type": "外出",
	}

	// 名字筛选
	if name != "" {
		customers, _, err := s.customerRepo.FindList(ctx, bson.M{"name": bson.M{"$regex": name, "$options": "i"}}, 0, 0)
		if err != nil {
			return nil, err
		}
		var customerIDs []primitive.ObjectID
		for _, c := range customers {
			customerIDs = append(customerIDs, c.ID)
		}
		if len(customerIDs) > 0 {
			filter["customer_id"] = bson.M{"$in": customerIDs}
		} else {
			return []map[string]interface{}{}, nil
		}
	}

	// 外出时间范围 (StartTime)
	if startDate != "" || endDate != "" {
		dateFilter := bson.M{}
		if startDate != "" {
			if t, err := time.Parse("2006-01-02", startDate); err == nil {
				dateFilter["$gte"] = t
			}
		}
		if endDate != "" {
			if t, err := time.Parse("2006-01-02", endDate); err == nil {
				t = t.Add(24 * time.Hour)
				dateFilter["$lt"] = t
			}
		}
		filter["start_time"] = dateFilter
	}

	records, _, err := s.recordRepo.FindList(ctx, filter, 0, 0)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, r := range records {
		customer, _ := s.customerRepo.FindById(ctx, r.CustomerID)
		customerName := "未知"
		phone := ""
		if customer != nil {
			customerName = customer.Name
			phone = customer.ContactPhone // 使用紧急联系人电话
		}

		isReturned := !r.EndTime.IsZero()
		currentStatus := "已外出"
		if isReturned {
			currentStatus = "已返回"
		}

		if status != "" && status != "全部" {
			if status == "已外出" && isReturned {
				continue
			}
			if status == "已返回" && !isReturned {
				continue
			}
		}

		item := map[string]interface{}{
			"id":            r.ID.Hex(),
			"customer_name": customerName,
			"contact_phone": phone,
			"outgoing_time": r.StartTime,
			//TODO登记时候显示实际返回时间
			"expected_return_time": "", // 暂无数据
			"actual_return_time":   r.EndTime,
			"status":               currentStatus,
			"destination":          r.Note, // 备注作为目的地
			// "customer_id" 用于前端操作 (登记返回等)?
			"customer_id": r.CustomerID.Hex(),
		}
		results = append(results, item)
	}
	return results, nil
}
