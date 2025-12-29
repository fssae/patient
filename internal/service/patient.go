package service

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PatientService struct {
	patientRepo *repository.PatientRepository
}

func NewPatientService(
	patientRepo *repository.PatientRepository,
) *PatientService {
	return &PatientService{
		patientRepo: patientRepo,
	}
}

func (s *PatientService) Login(ctx context.Context, req *domain.PatientLoginRequest) (*domain.PatientLoginResponse, error) {
	patient, err := s.patientRepo.Login(ctx, req.PatientId, req.Password)
	if err != nil {
		return nil, errors.New("工号或密码错误")
	}
	// 生成JWT
	tokenString, err := createToken(patient.Id, patient.PatientId)
	return &domain.PatientLoginResponse{
		Token:   tokenString,
		Patient: patient,
	}, nil
}

// Register 教师注册
func (s *PatientService) Register(ctx context.Context, req *domain.PatientRegisterRequest) error {
	err := s.patientRepo.Register(ctx, req.PatientPhone, req.Password)

	if err != nil { // 处理注册失败的情况
		return err
	}

	//注册成功
	return nil
}
func createToken(Id primitive.ObjectID, patientId string) (tokenString string, err error) {
	secret := viper.GetString("general.jwt")
	claims := domain.PatientClaims{
		//设置参数
		RegisteredClaims: jwt.RegisteredClaims{
			//设置300天的过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 300)),
		},
		Id:        Id,
		PatientId: patientId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//加密
	tokenStr, err := token.SignedString([]byte(secret))
	return tokenStr, err
}
