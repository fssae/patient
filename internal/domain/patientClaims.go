package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PatientClaims 患者JWT Claims
type PatientClaims struct {
	PatientId string             `json:"patientId"`
	Id        primitive.ObjectID `json:"id"`
	Name      string             `json:"name"`
	jwt.RegisteredClaims
}

// NewPatientClaims 创建患者Claims
func NewPatientClaims(patient *Patient) *PatientClaims {
	return &PatientClaims{
		Id: patient.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
}
