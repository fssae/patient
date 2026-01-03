package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// AnalysisLog 代表存储在 MongoDB 中的告警/分析数据
type AnalysisLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type      string             `bson:"type" json:"type"`           // "event"(事件) 或 "session"(会话)
	Timestamp interface{}        `bson:"timestamp" json:"timestamp"` // 支持 Long 或 String

	// 事件特定字段
	EventType string `bson:"event_type,omitempty" json:"event_type,omitempty"`
	Level     string `bson:"level,omitempty" json:"level,omitempty"`
	Message   string `bson:"message,omitempty" json:"message,omitempty"`
	VideoURL  string `bson:"video_url,omitempty" json:"video_url,omitempty"`

	PatientID  string `bson:"patient_id,omitempty" json:"patient_id,omitempty"`
	BedID      string `bson:"bed_id,omitempty" json:"bed_id,omitempty"`
	IsResolved bool   `bson:"is_resolved" json:"is_resolved"`

	// 会话特定字段
	AlertType  string        `bson:"alert_type,omitempty" json:"alert_type,omitempty"`
	Confidence float64       `bson:"confidence,omitempty" json:"confidence,omitempty"`
	DeviceID   string        `bson:"device_id,omitempty" json:"device_id,omitempty"`
	Location   string        `bson:"location,omitempty" json:"location,omitempty"`
	Details    *AlertDetails `bson:"details,omitempty" json:"details,omitempty"`

	CreatedAt int64 `bson:"created_at" json:"created_at"` // 内部存储时间
}

type AlertDetails struct {
	SourceVideo    string `bson:"source_video" json:"source_video"`
	VideoURL       string `bson:"video_url" json:"video_url"`
	LocalPath      string `bson:"local_path" json:"local_path"`
	FramesAnalysed int    `bson:"frames_analysed" json:"frames_analysed"`
	TimestampEnd   string `bson:"timestamp_end" json:"timestamp_end"`
}
