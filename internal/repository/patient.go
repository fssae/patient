package repository

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository/dao"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type PatientRepository struct {
	patientDAO *dao.PatientDAO
}

func NewPatientRepository(patientDAO *dao.PatientDAO) *PatientRepository {
	return &PatientRepository{
		patientDAO: patientDAO,
	}
}

// Create 创建教师
func (r *PatientRepository) Create(ctx context.Context, patient *domain.Patient) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(patient.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	patient.Password = string(hashedPassword)
	return r.patientDAO.Create(ctx, patient)
}

// Login 教师登录
func (r *PatientRepository) Login(ctx context.Context, patientId, password string) (*domain.Patient, error) {
	patient, err := r.patientDAO.FindByPatientAccount(ctx, patientId)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(patient.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return patient, nil
}

// Register 教师注册
func (r *PatientRepository) Register(ctx context.Context, patientphone, password string) error {
	//判断教师是否存在，用手机号判断
	patient, err := r.patientDAO.FindByPatientAccount(ctx, patientphone)

	//教师存在不可以在创建
	if err == nil {
		return errors.New("教师已存在")
	}

	//创建教师
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	patient = &domain.Patient{
		Phone:    patientphone,
		Password: string(hashedPassword),
	}
	return r.patientDAO.Create(ctx, patient)
}

// FindById 根据ID查找教师
func (r *PatientRepository) FindById(ctx context.Context, id primitive.ObjectID) (*domain.Patient, error) {
	return r.patientDAO.FindById(ctx, id)
}

// Update 更新教师信息
func (r *PatientRepository) Update(ctx context.Context, patient *domain.Patient) error {
	return r.patientDAO.Update(ctx, patient)
}
