package service

import (
	"classroom-analysis/internal/domain"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJWTService_GenerateAndValidatePatientToken(t *testing.T) {
	viper.Set("jwt.secret", "test-secret-key")
	viper.Set("jwt.expiration_hours", 24)
	viper.Set("jwt.refresh_threshold_hours", 1)

	service := NewJWTService()

	patientID := primitive.NewObjectID()
	now := time.Now()
	claims := &domain.PatientClaims{
		Id:          patientID,
		PatientId:   "P12345",
		Name:        "Test Patient",
		RefreshedAt: now.Unix(),
		TokenType:   "patient",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token, err := service.GeneratePatientToken(claims)
	if err != nil {
		t.Fatalf("Failed to generate patient token: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}

	validatedClaims, err := service.ValidatePatientToken(token)
	if err != nil {
		t.Fatalf("Failed to validate patient token: %v", err)
	}

	if validatedClaims.PatientId != claims.PatientId {
		t.Errorf("Expected PatientId %s, got %s", claims.PatientId, validatedClaims.PatientId)
	}

	if validatedClaims.TokenType != "patient" {
		t.Errorf("Expected TokenType 'patient', got %s", validatedClaims.TokenType)
	}
}

func TestJWTService_GenerateAndValidateUserToken(t *testing.T) {
	viper.Set("jwt.secret", "test-secret-key")
	viper.Set("jwt.expiration_hours", 24)

	service := NewJWTService()

	userID := primitive.NewObjectID()
	now := time.Now()
	claims := &domain.UserClaims{
		UserID:      userID,
		Phone:       "13800138000",
		Name:        "Test User",
		RefreshedAt: now.Unix(),
		TokenType:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token, err := service.GenerateUserToken(claims)
	if err != nil {
		t.Fatalf("Failed to generate user token: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}

	validatedClaims, err := service.ValidateUserToken(token)
	if err != nil {
		t.Fatalf("Failed to validate user token: %v", err)
	}

	if validatedClaims.Phone != claims.Phone {
		t.Errorf("Expected Phone %s, got %s", claims.Phone, validatedClaims.Phone)
	}

	if validatedClaims.TokenType != "user" {
		t.Errorf("Expected TokenType 'user', got %s", validatedClaims.TokenType)
	}
}

func TestJWTService_RefreshPatientToken(t *testing.T) {
	viper.Set("jwt.secret", "test-secret-key")
	viper.Set("jwt.expiration_hours", 24)

	service := NewJWTService()

	patientID := primitive.NewObjectID()
	now := time.Now()
	oldRefreshedAt := now.Add(-2 * time.Hour).Unix()

	claims := &domain.PatientClaims{
		Id:          patientID,
		PatientId:   "P12345",
		Name:        "Test Patient",
		RefreshedAt: oldRefreshedAt,
		TokenType:   "patient",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now.Add(-23 * time.Hour)),
			NotBefore: jwt.NewNumericDate(now.Add(-23 * time.Hour)),
		},
	}

	newToken, err := service.RefreshPatientToken(claims)
	if err != nil {
		t.Fatalf("Failed to refresh patient token: %v", err)
	}

	if newToken == "" {
		t.Fatal("Refreshed token is empty")
	}

	validatedClaims, err := service.ValidatePatientToken(newToken)
	if err != nil {
		t.Fatalf("Failed to validate refreshed token: %v", err)
	}

	if validatedClaims.RefreshedAt == oldRefreshedAt {
		t.Error("RefreshedAt timestamp was not updated")
	}

	if validatedClaims.RefreshedAt <= oldRefreshedAt {
		t.Errorf("Expected RefreshedAt to be greater than %d, got %d", oldRefreshedAt, validatedClaims.RefreshedAt)
	}
}

func TestJWTService_ShouldRefreshToken(t *testing.T) {
	viper.Set("jwt.expiration_hours", 24)
	viper.Set("jwt.refresh_threshold_hours", 1)

	service := NewJWTService()

	tests := []struct {
		name      string
		expiresAt *jwt.NumericDate
		expected  bool
	}{
		{
			name:      "Token expiring in 30 minutes - should refresh",
			expiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			expected:  true,
		},
		{
			name:      "Token expiring in 45 minutes - should refresh",
			expiresAt: jwt.NewNumericDate(time.Now().Add(45 * time.Minute)),
			expected:  true,
		},
		{
			name:      "Token expiring in 2 hours - should not refresh",
			expiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			expected:  false,
		},
		{
			name:      "Token expiring in 10 hours - should not refresh",
			expiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Hour)),
			expected:  false,
		},
		{
			name:      "Token already expired - should not refresh",
			expiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			expected:  false,
		},
		{
			name:      "Nil expiresAt - should not refresh",
			expiresAt: nil,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ShouldRefreshToken(tt.expiresAt)
			if result != tt.expected {
				t.Errorf("ShouldRefreshToken() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestJWTService_InvalidToken(t *testing.T) {
	viper.Set("jwt.secret", "test-secret-key")

	service := NewJWTService()

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "Empty token",
			token: "",
		},
		{
			name:  "Invalid format",
			token: "not.a.valid.jwt.token",
		},
		{
			name:  "Malformed token",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ValidatePatientToken(tt.token)
			if err == nil {
				t.Error("Expected error for invalid token, got nil")
			}

			_, err = service.ValidateUserToken(tt.token)
			if err == nil {
				t.Error("Expected error for invalid token, got nil")
			}
		})
	}
}
