package dao

import (
	"classroom-analysis/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PatientDAO struct {
	collection *mongo.Collection
}

func NewPatientDAO(db *mongo.Database) *PatientDAO {
	return &PatientDAO{
		collection: db.Collection("patients"),
	}
}

// Create 创建教师
func (dao *PatientDAO) Create(ctx context.Context, patient *domain.Patient) error {
	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()
	result, err := dao.collection.InsertOne(ctx, patient)
	if err != nil {
		return err
	}
	patient.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByPatientId 根据工号查找教师
func (dao *PatientDAO) FindByPatientAccount(ctx context.Context, patientId string) (*domain.Patient, error) {
	var patient domain.Patient
	err := dao.collection.FindOne(ctx, bson.M{"patientId": patientId}).Decode(&patient)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// FindById 根据ID查找教师
func (dao *PatientDAO) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Patient, error) {
	var patient domain.Patient
	err := dao.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&patient)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// Update 更新教师信息
func (dao *PatientDAO) Update(ctx context.Context, patient *domain.Patient) error {
	patient.UpdatedAt = time.Now()
	_, err := dao.collection.UpdateOne(
		ctx,
		bson.M{"_id": patient.Id},
		bson.M{"$set": patient},
	)
	return err
}
