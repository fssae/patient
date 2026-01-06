package domain

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// UserClaims 用户JWT Claims
type UserClaims struct {
    UserID      primitive.ObjectID `json:"user_id"`
    Phone       string             `json:"phone"`
    Name        string             `json:"name"`
    RefreshedAt int64              `json:"refreshedAt"` // 最后刷新时间（Unix时间戳）
    TokenType   string             `json:"tokenType"`   // Token类型: "user"
    jwt.RegisteredClaims
}

// NewUserClaims 创建用户Claims
func NewUserClaims(user *User) *UserClaims {
    now := time.Now()
    return &UserClaims{
        UserID:      user.ID,
        Phone:       user.Phone,
        Name:        user.Name,
        RefreshedAt: now.Unix(),
        TokenType:   "user",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour * 300)), // 300天
            IssuedAt:  jwt.NewNumericDate(now),
            NotBefore: jwt.NewNumericDate(now),
        },
    }
}

