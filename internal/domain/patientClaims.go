package domain

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// PatientClaims 患者JWT Claims
type PatientClaims struct {
    PatientId   string             `json:"patientId"`
    Id          primitive.ObjectID `json:"id"`
    Name        string             `json:"name"`
    RefreshedAt int64              `json:"refreshedAt"` // 最后刷新时间（Unix时间戳）
    TokenType   string             `json:"tokenType"`   // Token类型: "patient"
    jwt.RegisteredClaims
}

// NewPatientClaims 创建患者Claims
func NewPatientClaims(patient *Patient) *PatientClaims {
    now := time.Now()
    return &PatientClaims{
        Id:          patient.Id,
        RefreshedAt: now.Unix(),
        TokenType:   "patient",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
        },
    }
}
