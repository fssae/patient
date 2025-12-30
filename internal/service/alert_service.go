package service

import (
	"classroom-analysis/internal/mq"
	"classroom-analysis/internal/ws"
	"context"
	"encoding/json"
	"log"
)

type AlertService struct {
	consumer *mq.AlertConsumer
	wsMgr    *ws.WebSocketManager
}

func NewAlertService(consumer *mq.AlertConsumer, wsMgr *ws.WebSocketManager) *AlertService {
	return &AlertService{
		consumer: consumer,
		wsMgr:    wsMgr,
	}
}

func (s *AlertService) Start() {
	// Set the callback for the consumer
	s.consumer.SetMessageHandler(func(msg *mq.AlertMessage) {
		log.Printf("Received alert: %+v", msg)

		// Serialize message to JSON
		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling alert message: %v", err)
			return
		}

		// Broadcast to all WebSocket clients
		s.wsMgr.SendBytes(data)
	})

	// Start the consumer in a background context (or passing a context from App)
	// For now, we use Background, effectively running until app exit.
	go s.consumer.Start(context.Background())
}
