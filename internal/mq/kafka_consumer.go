package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

// AlertMessage defines the structure of the alert from Python algorithm
type AlertMessage struct {
	EventType string `json:"event_type"` // fall, help, emotion
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"` // critical, warning, info
	Message   string `json:"message"`
	VideoURL  string `json:"video_url"`
}

// AlertConsumer handles consuming messages from Kafka
type AlertConsumer struct {
	brokers []string
	topic   string
	groupID string
	handler func(msg *AlertMessage)
}

func NewAlertConsumer() *AlertConsumer {
	brokers := viper.GetStringSlice("kafka.brokers")
	if len(brokers) == 0 {
		brokers = []string{"82.156.64.69:9092"} // Fallback to provided IP
	}
	// Fallback/Default config can be adjusted
	return &AlertConsumer{
		brokers: brokers,
		topic:   "elderly_alerts",
		groupID: "elderly-alert-group",
	}
}

// SetMessageHandler sets the callback for when a valid alert is received
func (c *AlertConsumer) SetMessageHandler(handler func(msg *AlertMessage)) {
	c.handler = handler
}

// Start begins the consumer loop
func (c *AlertConsumer) Start(ctx context.Context) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0 // Adjust based on your Kafka version
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Auth configuration if needed, matching existing project style
	username := viper.GetString("kafka.username")
	password := viper.GetString("kafka.password")
	if username != "" && password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.Handshake = true
	}

	client, err := sarama.NewConsumerGroup(c.brokers, c.groupID, config)
	if err != nil {
		log.Printf("Error creating consumer group client: %v", err)
		return
	}

	handler := &consumerGroupHandler{
		callback: c.handler,
	}

	go func() {
		defer client.Close()
		for {
			if err := client.Consume(ctx, []string{c.topic}, handler); err != nil {
				log.Printf("Error from consumer: %v", err)
				time.Sleep(time.Second * 5) // Retry delay
			}
			// Check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()
	log.Printf("Alert Consumer started on topic %s", c.topic)
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	callback func(msg *AlertMessage)
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var alert AlertMessage
		if err := json.Unmarshal(msg.Value, &alert); err != nil {
			log.Printf("Error unmarshaling alert message: %v. Raw message: %s", err, string(msg.Value))
			sess.MarkMessage(msg, "")
			continue
		}

		if h.callback != nil {
			h.callback(&alert)
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
