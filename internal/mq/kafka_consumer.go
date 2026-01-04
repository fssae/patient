package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

// AlertDetails 定义了会话完成消息中的嵌套详细信息
type AlertDetails struct {
	SourceVideo    string `json:"source_video"`
	VideoURL       string `json:"video_url"`
	LocalPath      string `json:"local_path"`
	FramesAnalysed int    `json:"frames_analysed"`
	TimestampEnd   string `json:"timestamp_end"`
}

// AlertMessage 定义了来自 Python 算法的告警结构。
// 它处理即时事件告警和会话完成总结。
type AlertMessage struct {
	ID        string      `json:"id"`
	Timestamp interface{} `json:"timestamp"` // Unix 时间戳 (秒)

	// --- 类型 1: 事件告警 (存在 event_type) ---
	EventType string `json:"event_type,omitempty"` // fall(跌倒), pain_expression(疼痛表情), call_help(呼救)
	Level     string `json:"level,omitempty"`      // critical(严重), warning(警告), info(信息)
	Message   string `json:"message,omitempty"`
	VideoURL  string `json:"video_url,omitempty"`

	PatientID  string `json:"patient_id,omitempty"`
	BedID      string `json:"bed_id,omitempty"`
	IsResolved bool   `json:"isResolved"` // 初始状态必须为 false

	// --- 类型 2: 会话完成 (alert_type 为 "session_completed") ---
	AlertType  string        `json:"alert_type,omitempty"`
	Confidence float64       `json:"confidence,omitempty"`
	DeviceID   string        `json:"device_id,omitempty"`
	Location   string        `json:"location,omitempty"`
	Details    *AlertDetails `json:"details,omitempty"`
}

// AlertConsumer 处理从 Kafka 消费消息
type AlertConsumer struct {
	brokers []string
	topic   string
	groupID string
	handler func(msg *AlertMessage)
}

func NewAlertConsumer() *AlertConsumer {
	brokers := viper.GetStringSlice("kafka.brokers")
	if len(brokers) == 0 {
		brokers = []string{"82.156.64.69:9092"} // 备用 IP
	}
	// 默认配置可以根据需要调整
	return &AlertConsumer{
		brokers: brokers,
		topic:   "elderly_alerts",
		groupID: "elderly-alert-group",
	}
}

// SetMessageHandler 设置收到有效告警时的回调函数
func (c *AlertConsumer) SetMessageHandler(handler func(msg *AlertMessage)) {
	c.handler = handler
}

// Start 启动消费者循环
func (c *AlertConsumer) Start(ctx context.Context) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0 // 根据 Kafka 版本调整
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// 添加超时配置,防止连接阻塞
	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second
	config.Metadata.Timeout = 10 * time.Second

	// 消费者组会话超时配置
	sessionTimeout := viper.GetDuration("kafka.consumer.group_session_timeout")
	if sessionTimeout == 0 {
		sessionTimeout = 30 * time.Second // 默认 30 秒
	}
	config.Consumer.Group.Session.Timeout = sessionTimeout

	heartbeatInterval := viper.GetDuration("kafka.consumer.heartbeat_interval")
	if heartbeatInterval == 0 {
		heartbeatInterval = 3 * time.Second // 默认 3 秒
	}
	config.Consumer.Group.Heartbeat.Interval = heartbeatInterval

	// 设置重平衡超时
	rebalanceTimeout := viper.GetDuration("kafka.consumer.rebalance_timeout")
	if rebalanceTimeout == 0 {
		rebalanceTimeout = 60 * time.Second // 默认 60 秒
	}
	config.Consumer.Group.Rebalance.Timeout = rebalanceTimeout

	// 如果需要认证，匹配现有项目风格
	username := viper.GetString("kafka.username")
	password := viper.GetString("kafka.password")
	if username != "" && password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.Handshake = true
	}

	// 使用带超时的上下文创建客户端
	createCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// 在独立 goroutine 中创建客户端,避免阻塞主线程
	clientChan := make(chan sarama.ConsumerGroup, 1)
	errChan := make(chan error, 1)

	go func() {
		client, err := sarama.NewConsumerGroup(c.brokers, c.groupID, config)
		if err != nil {
			errChan <- err
			return
		}
		clientChan <- client
	}()

	var client sarama.ConsumerGroup
	select {
	case client = <-clientChan:
		log.Printf("✓ Kafka 消费者组客户端创建成功")
	case err := <-errChan:
		log.Printf("✗ 创建 Kafka 消费者组客户端失败: %v (将在后台重试)", err)
		// 启动后台重连逻辑
		go c.retryConnect(ctx)
		return
	case <-createCtx.Done():
		log.Printf("✗ 创建 Kafka 消费者组客户端超时 (将在后台重试)")
		go c.retryConnect(ctx)
		return
	}

	handler := &consumerGroupHandler{
		callback: c.handler,
	}

	go func() {
		defer client.Close()
		for {
			if err := client.Consume(ctx, []string{c.topic}, handler); err != nil {
				log.Printf("Kafka 消费者错误: %v (将在 5 秒后重试)", err)
				time.Sleep(time.Second * 5) // 重试延迟
			}
			// 检查上下文是否已取消，发出消费者停止信号
			if ctx.Err() != nil {
				log.Printf("Kafka 消费者已停止")
				return
			}
		}
	}()
	log.Printf("✓ Kafka 告警消费者已在主题 '%s' 上启动", c.topic)
}

// retryConnect 后台重试连接 Kafka
func (c *AlertConsumer) retryConnect(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	retryCount := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			retryCount++
			log.Printf("尝试重新连接 Kafka (第 %d 次)...", retryCount)

			// 尝试重新启动
			c.Start(ctx)
			return // Start 会处理后续逻辑
		}
	}
}

// consumerGroupHandler 实现 sarama.ConsumerGroupHandler 接口
type consumerGroupHandler struct {
	callback func(msg *AlertMessage)
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var alert AlertMessage
		if err := json.Unmarshal(msg.Value, &alert); err != nil {
			log.Printf("解码告警消息错误: %v. 原始消息: %s", err, string(msg.Value))
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
