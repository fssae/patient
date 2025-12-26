package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserClaims 用户JWT Claims
type UserClaims struct {
	UserID primitive.ObjectID `json:"user_id"`
	Phone  string             `json:"phone"`
	Name   string             `json:"name"`
	jwt.RegisteredClaims
}

// NewUserClaims 创建用户Claims
func NewUserClaims(user *User) *UserClaims {
	return &UserClaims{
		UserID: user.ID,
		Phone:  user.Phone,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 300)), // 300天
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
}

