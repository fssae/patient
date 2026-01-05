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
	// 使用唯一的消费者组ID，避免与其他消费者冲突
	// 每次启动使用固定ID，确保消费组状态一致
	return &AlertConsumer{
		brokers: brokers,
		topic:   "elderly_alerts",
		groupID: "elderly-alert-backend-v2", // 使用新的消费者组ID
	}
}

// SetMessageHandler 设置收到有效告警时的回调函数
func (c *AlertConsumer) SetMessageHandler(handler func(msg *AlertMessage)) {
	c.handler = handler
}

// Start 启动消费者循环
func (c *AlertConsumer) Start(ctx context.Context) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0                      // 与 ioc/kafka.go 保持一致
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // 从最新消息开始，避免重复消费

	// 消费者组配置 - 确保正确的分区分配
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	config.Consumer.Group.Session.Timeout = 30 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	// 静态成员ID - 减少重平衡
	// 使用固定的实例ID，重启后重新加入不会触发不必要的重平衡
	config.Consumer.Group.InstanceId = "elderly-alert-backend-instance-1"

	// 如果需要认证，匹配现有项目风格
	username := viper.GetString("kafka.username")
	password := viper.GetString("kafka.password")

	log.Printf("========== Kafka 告警消费者配置 ==========")
	log.Printf("Brokers: %v", c.brokers)
	log.Printf("Topic: %s", c.topic)
	log.Printf("GroupID: %s", c.groupID)
	log.Printf("Username: %s", username)
	log.Printf("Password 长度: %d", len(password))

	if username != "" && password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.Handshake = true
		config.Net.TLS.Enable = false // 确保关闭 TLS，使用 SASL_PLAINTEXT
		log.Printf("SASL 认证已启用: User=%s, Mechanism=%s", username, sarama.SASLTypePlaintext)
	} else {
		log.Printf("警告: SASL 认证未启用 (username 或 password 为空)")
	}

	log.Printf("尝试连接到 Kafka brokers: %v", c.brokers)
	client, err := sarama.NewConsumerGroup(c.brokers, c.groupID, config)
	if err != nil {
		log.Printf("创建消费者组客户端错误: %v", err)
		log.Printf("请检查: 1) Broker 地址是否正确 2) 端口是否开放 3) SASL 认证信息是否正确")
		return
	}
	log.Printf("成功创建消费者组客户端")

	handler := &consumerGroupHandler{
		callback: c.handler,
	}

	go func() {
		defer client.Close()
		for {
			if err := client.Consume(ctx, []string{c.topic}, handler); err != nil {
				log.Printf("消费者错误: %v", err)
				time.Sleep(time.Second * 5) // 重试延迟
			}
			// 检查上下文是否已取消，发出消费者停止信号
			if ctx.Err() != nil {
				return
			}
		}
	}()
	log.Printf("告警消费者已在主题 %s 上启动", c.topic)
}

// consumerGroupHandler 实现 sarama.ConsumerGroupHandler 接口
type consumerGroupHandler struct {
	callback func(msg *AlertMessage)
}

func (h *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("消费者组 Setup - MemberID: %s, GenerationID: %d, Claims: %v",
		session.MemberID(), session.GenerationID(), session.Claims())
	return nil
}

func (h *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Printf("消费者组 Cleanup - MemberID: %s", session.MemberID())
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Printf("开始消费分区 - Topic: %s, Partition: %d, InitialOffset: %d",
		claim.Topic(), claim.Partition(), claim.InitialOffset())

	for msg := range claim.Messages() {
		log.Printf("收到消息 - Topic: %s, Partition: %d, Offset: %d, 内容长度: %d bytes",
			msg.Topic, msg.Partition, msg.Offset, len(msg.Value))

		var alert AlertMessage
		if err := json.Unmarshal(msg.Value, &alert); err != nil {
			log.Printf("解码告警消息错误: %v. 原始消息: %s", err, string(msg.Value))
			sess.MarkMessage(msg, "")
			continue
		}

		log.Printf("成功解析告警 - EventType: %s, AlertType: %s", alert.EventType, alert.AlertType)

		if h.callback != nil {
			h.callback(&alert)
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
