package service

import (
	"classroom-analysis/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWTService struct {
	secret                string
	expirationHours       int
	refreshThresholdHours int
}

func NewJWTService() *JWTService {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		secret = viper.GetString("general.jwt") // fallback to old config
	}
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}

	expirationHours := viper.GetInt("jwt.expiration_hours")
	if expirationHours == 0 {
		expirationHours = 24
	}

	refreshThresholdHours := viper.GetInt("jwt.refresh_threshold_hours")
	if refreshThresholdHours == 0 {
		refreshThresholdHours = 1
	}

	return &JWTService{
		secret:                secret,
		expirationHours:       expirationHours,
		refreshThresholdHours: refreshThresholdHours,
	}
}

// GeneratePatientToken 生成患者Token
func (s *JWTService) GeneratePatientToken(claims *domain.PatientClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// GenerateUserToken 生成用户Token
func (s *JWTService) GenerateUserToken(claims *domain.UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// ValidatePatientToken 验证患者Token
func (s *JWTService) ValidatePatientToken(tokenStr string) (*domain.PatientClaims, error) {
	claims := &domain.PatientClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token无效")
	}

	return claims, nil
}

// ValidateUserToken 验证用户Token
func (s *JWTService) ValidateUserToken(tokenStr string) (*domain.UserClaims, error) {
	claims := &domain.UserClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token无效")
	}

	return claims, nil
}

// RefreshPatientToken 刷新患者Token
func (s *JWTService) RefreshPatientToken(claims *domain.PatientClaims) (string, error) {
	now := time.Now()
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Duration(s.expirationHours) * time.Hour))
	claims.RefreshedAt = now.Unix()
	return s.GeneratePatientToken(claims)
}

// RefreshUserToken 刷新用户Token
func (s *JWTService) RefreshUserToken(claims *domain.UserClaims) (string, error) {
	now := time.Now()
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Duration(s.expirationHours) * time.Hour))
	claims.RefreshedAt = now.Unix()
	return s.GenerateUserToken(claims)
}

// ShouldRefreshToken 判断Token是否应该刷新
func (s *JWTService) ShouldRefreshToken(expiresAt *jwt.NumericDate) bool {
	if expiresAt == nil {
		return false
	}
	timeUntilExpiry := time.Until(expiresAt.Time)
	threshold := time.Duration(s.refreshThresholdHours) * time.Hour
	return timeUntilExpiry > 0 && timeUntilExpiry < threshold
}

// GetSecret 获取密钥
func (s *JWTService) GetSecret() string {
	return s.secret
}

// GetExpirationDuration 获取过期时间
func (s *JWTService) GetExpirationDuration() time.Duration {
	return time.Duration(s.expirationHours) * time.Hour
}
