package service

import (
	"classroom-analysis/internal/domain"
	"classroom-analysis/internal/mq"
	"classroom-analysis/internal/ws"
	"context"
	"encoding/json"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlertService struct {
	consumer    *mq.AlertConsumer
	wsMgr       *ws.WebSocketManager
	analysisSvc AnalysisService
}

func NewAlertService(consumer *mq.AlertConsumer, wsMgr *ws.WebSocketManager, analysisSvc AnalysisService) *AlertService {
	return &AlertService{
		consumer:    consumer,
		wsMgr:       wsMgr,
		analysisSvc: analysisSvc,
	}
}

func (s *AlertService) Start() {
	// 为消费者设置回调函数
	s.consumer.SetMessageHandler(func(msg *mq.AlertMessage) {
		log.Printf("收到告警: %+v", msg)

		// 1. 根据前端规范初始化字段
		analysisLog := s.mapToAnalysisLog(msg)
		if analysisLog.ID.IsZero() {
			analysisLog.ID = primitive.NewObjectID()
		}
		analysisLog.IsResolved = false

		// 使用生成的 ID 和状态更新消息
		msg.ID = analysisLog.ID.Hex()
		msg.IsResolved = false

		// 2. 保存到 MongoDB (异步执行即可，因为我们已经有了 ID)
		go func() {
			err := s.analysisSvc.RecordAnalysis(context.Background(), analysisLog)
			if err != nil {
				log.Printf("记录分析日志到 MongoDB 时出错: %v", err)
			}
		}()

		// 3. 将消息序列化为 JSON
		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("序列化告警消息时出错: %v", err)
			return
		}

		// 4. 广播给所有 WebSocket 客户端
		s.wsMgr.SendBytes(data)
	})

	// 在后台上下文中启动消费者
	go s.consumer.Start(context.Background())
}

func (s *AlertService) mapToAnalysisLog(msg *mq.AlertMessage) *domain.AnalysisLog {
	logType := "event"
	if msg.AlertType == "session_completed" {
		logType = "session"
	}

	analysisLog := &domain.AnalysisLog{
		Type:       logType,
		Timestamp:  msg.Timestamp,
		EventType:  msg.EventType,
		Level:      msg.Level,
		Message:    msg.Message,
		VideoURL:   msg.VideoURL,
		AlertType:  msg.AlertType,
		Confidence: msg.Confidence,
		DeviceID:   msg.DeviceID,
		Location:   msg.Location,
	}

	if msg.Details != nil {
		analysisLog.Details = &domain.AlertDetails{
			SourceVideo:    msg.Details.SourceVideo,
			VideoURL:       msg.Details.VideoURL,
			LocalPath:      msg.Details.LocalPath,
			FramesAnalysed: msg.Details.FramesAnalysed,
			TimestampEnd:   msg.Details.TimestampEnd,
		}
	}

	analysisLog.PatientID = msg.PatientID
	analysisLog.BedID = msg.BedID

	return analysisLog
}
