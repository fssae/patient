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
        return nil, errors.New("用户名或密码错误")
    }
    // 生成JWT
    tokenString, err := createToken(patient.Id, patient.PatientId)
    return &domain.PatientLoginResponse{
        Token:   tokenString,
        Patient: patient,
    }, nil
}

// Register 患者注册
func (s *PatientService) Register(ctx context.Context, req *domain.PatientRegisterRequest) error {
    err := s.patientRepo.Register(ctx, req.PatientPhone, req.Password)

    if err != nil { // 处理注册失败的情况
        return err
    }

    //注册成功
    return nil
}
func createToken(Id primitive.ObjectID, patientId string) (tokenString string, err error) {
    secret := viper.GetString("jwt.secret")
    if secret == "" {
        secret = viper.GetString("general.jwt")
    }
    
    expirationHours := viper.GetInt("jwt.expiration_hours")
    if expirationHours == 0 {
        expirationHours = 24
    }
    
    now := time.Now()
    claims := domain.PatientClaims{
        Id:          Id,
        PatientId:   patientId,
        RefreshedAt: now.Unix(),
        TokenType:   "patient",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * time.Duration(expirationHours))),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err := token.SignedString([]byte(secret))
    return tokenStr, err
}
