package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Patient struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PatientId string             `json:"patientId" bson:"patientId"` // 患者ID
	Name      string             `json:"name" bson:"name"`           // 姓名
	Password  string             `json:"password" bson:"password"`   // 密码
	Email     string             `json:"email" bson:"email"`         // 邮箱
	Phone     string             `json:"phone" bson:"phone"`         // 电话
	Avatar    string             `json:"avatar" bson:"avatar"`       // 头像
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type PatientLoginRequest struct {
	PatientId string `json:"patientId" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type PatientRegisterRequest struct {
	PatientPhone string `json:"patientphone" binding:"required"`
	Password     string `json:"password" binding:"required"`
}

type PatientLoginResponse struct {
	Token   string   `json:"token"`
	Patient *Patient `json:"patient"`
}

type PatientSettings struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PatientId    primitive.ObjectID `json:"patientId" bson:"patientId"`
	Email        string             `json:"email" bson:"email"`
	Notification bool               `json:"notification" bson:"notification"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}
