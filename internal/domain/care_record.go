package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CareRecord 护理记录
type CareRecords struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerID   primitive.ObjectID `json:"customer_id,omitempty" bson:"customer_id,omitempty"`
	CustomerName string             `json:"customer_name" bson:"customer_name" binding:"required"` // 老人姓名
	Records      []RecordItems      `json:"records" bson:"records"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
type RecordItems struct {
	CareItem      string    `json:"care_item" bson:"care_item" binding:"required"`           // 护理项目
	CareTime      time.Time `json:"care_time" bson:"care_time" binding:"required"`           // 护理时间
	CarePersonnel string    `json:"care_personnel" bson:"care_personnel" binding:"required"` // 护理人员
	CareResult    string    `json:"care_result" bson:"care_result" binding:"required"`       // 护理结果
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}
