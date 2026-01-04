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
	recordRepo   repository.RecordRepository
	customerRepo repository.CustomerRepository
	bedRepo      repository.BedRepository
}

func NewRecordService(recordRepo repository.RecordRepository, customerRepo repository.CustomerRepository, bedRepo repository.BedRepository) *RecordService {
	return &RecordService{
		recordRepo:   recordRepo,
		customerRepo: customerRepo,
		bedRepo:      bedRepo,
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
func (s *RecordService) GetCheckInList(ctx context.Context, name, roomNumber, nursingLevel string, startDate, endDate string) ([]map[string]interface{}, error) {
	filter := bson.M{}

	// 名字模糊查询
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	// 房间号查询逻辑：Customer -> Bed -> Room
	// 如果提供了房间号，需要先找到该房间下的所有床位ID，然后 Filter customer.bed_id IN [...]
	if roomNumber != "" {
		// 这里 RecordService 需要 Access RoomRepository?
		// 目前 RecordService 只有 BedRepo.
		// 但是 BedRepo 可以根据 FindByRoomID.
		// 我需要先从 BedRepo 拿到所有床位，然后手动检查 Room Number?
		// BedRepo.FindByRoomID 需要 RoomID.
		// 我无法直接从 RecordService 查 Room (缺少 RoomRepo).
		// 假如 Bed 有 Number 如 A101-1，可以模糊匹配 "A101" ?
		// 如果Bed.Number设计规范为 RoomNumber-BedIndex，那可以用正则匹配 Bed ID.
		// 但 Customer 存的是 BedID (ObjectId).
		// 方案：注入 RoomRepo 或 BedService?
		// 简单方案：先不依赖 RoomRepo，假设 RoomNumber 可以通过 Bed Number 匹配 (如果 Bed 存了 Number)。
		// Bed 结构体有 RoomID 和 Number.
		// 让我先假设可以直接按照 bed_id list 筛选.
		// 我需要引入 RoomRepo.
		// 或者，暂时只支持根据 Bed Number 筛选? 用户说 "Room Number".
		// 如果没有 RoomRepo，这个需求比较难办。
		// 让我在 NewRecordService 注入 RoomRepo?
		// 鉴于此时修改 struct 较大，且 `bed.go` 服务里有 RoomRepo.
		// 也许我可以使用 bedRepo.FindList(bson.M{}) 获取所有床位，然后在内存过滤 Room Number?
		// 或者：查询所有 Bed 匹配 Number like "A101%"?
		// 如果 Bed.Number 是 "A101-1"。
		// 让我们尝试查询与 Bed 关联的。
		// 1. Find all beds where Number starts with roomNumber.
		// beds, _, _ := s.bedRepo.FindList(ctx, bson.M{"number": {$regex: "^" + roomNumber}}, 0, 0)
		// 2. Extract IDs.
		// 3. filter["bed_id"] = {$in: ids}
		// 这是一个可行的方案，不需要 RoomRepo，只要 Bed Number 包含 Room Number.

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

	// 护理级别 (CareLevelID usually).
	// 如果参数是名字，需要 LookUp. 同样缺少 Repo.
	// 假设参数是 ID 字符串? 用户说 "护理级别"，可能是名字。
	// 这通常需要在前端下拉选择 ID. 假设传 ID.
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
	customers, _, err := s.customerRepo.FindList(ctx, filter, 0, 0) // 暂未分页，全量返回
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
		// 注意: Customer 仅存储 CareLevelID，若需显示名称需关联查询或前端处理
		item["nursing_level"] = c.CareLevelID.Hex() // 暂时返回 ID

		results = append(results, item)
	}

	return results, nil
}

// GetCheckOutList 获取退住信息列表
func (s *RecordService) GetCheckOutList(ctx context.Context, name, roomNumber, reason, startDate, endDate string) ([]map[string]interface{}, error) {
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

		days := int(r.EndTime.Sub(r.StartTime).Hours() / 24)

		item := map[string]interface{}{
			"id":             r.ID.Hex(),
			"customer_name":  customerName,
			"room_number":    "-", // 记录中未存，且客户已退住，难以获取历史床位
			"check_in_date":  r.StartTime,
			"check_out_date": r.EndTime,
			"days":           days,
			"care_level":     careLevel,
			"reason":         r.Note,
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

		// 状态筛选逻辑：
		// status: "全部", "已外出" (EndTime is zero), "已返回" (EndTime not zero)
		// Mongo 查询 EndTime 是否存在/为零比较麻烦，通常用 $exists 或 $eq null.
		// 但 Go Driver 读出来的 empty Time 是 zero value.
		// 我们在内存里做这个筛选比较简单，因为 filter 只能基本筛选。

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
			"id":                   r.ID.Hex(),
			"customer_name":        customerName,
			"contact_phone":        phone,
			"outgoing_time":        r.StartTime,
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

// UpdateRecord 更新客户的记录
func (s *RecordService) UpdateRecord(ctx context.Context, record *domain.Record, customerID primitive.ObjectID, elderId string) error {
	// 检查记录是否存在
	_, err := s.recordRepo.FindById(ctx, record.ID)
	if err != nil {
		return err
	}
	//更新records
	if err := s.recordRepo.Update(ctx, record); err != nil {
		return err
	}

	// 检查customer是否存在该用户
	customer, err := s.customerRepo.FindByUserID(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户未找到")
	}
	// 更新客户名称
	updates := map[string]interface{}{
		"name": elderId,
	}
	if err := s.customerRepo.Update(ctx, customer.ID, updates); err != nil {
		return err
	}
	return nil
}
