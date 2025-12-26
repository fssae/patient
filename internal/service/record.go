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
func (s *RecordService) CheckIn(ctx context.Context, customerID primitive.ObjectID, bedID primitive.ObjectID, note, createdBy string) error {
	// 验证客户是否存在
	customer, err := s.customerRepo.FindById(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("客户不存在")
	}

	// 如果指定了床位，分配床位
	if !bedID.IsZero() {
		err = s.bedRepo.AssignToCustomer(ctx, bedID, customerID)
		if err != nil {
			return err
		}
		customer.BedID = bedID
	}

	// 更新客户状态
	customer.Status = "入住中"
	customer.CheckInDate = time.Now()
	err = s.customerRepo.Update(ctx, customer)
	if err != nil {
		return err
	}

	// 创建入住记录
	record := &domain.Record{
		CustomerID: customerID,
		Type:       "入住",
		StartTime:  time.Now(),
		Note:       note,
		CreatedBy:  createdBy,
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
	customer.Status = "已退住"
	customer.CheckOutDate = time.Now()
	err = s.customerRepo.Update(ctx, customer)
	if err != nil {
		return err
	}

	// 创建退住记录
	record := &domain.Record{
		CustomerID: customerID,
		Type:      "退住",
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
	customer.Status = "外出中"
	err = s.customerRepo.Update(ctx, customer)
	if err != nil {
		return err
	}

	// 创建外出记录
	record := &domain.Record{
		CustomerID: customerID,
		Type:      "外出",
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
	customer.Status = "入住中"
	err = s.customerRepo.Update(ctx, customer)
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

