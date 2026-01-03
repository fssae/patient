package ws

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 目前允许所有来源
	},
}

type WebSocketManager struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

func NewWebSocketManager() *WebSocketManager {
	manager := &WebSocketManager{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
	go manager.run()
	return manager
}

func (manager *WebSocketManager) run() {
	for {
		select {
		case conn := <-manager.register:
			manager.mu.Lock()
			manager.clients[conn] = true
			manager.mu.Unlock()
			log.Println("新的 WebSocket 客户端已连接")
		case conn := <-manager.unregister:
			manager.mu.Lock()
			if _, ok := manager.clients[conn]; ok {
				delete(manager.clients, conn)
				conn.Close()
				log.Println("WebSocket 客户端已断开")
			}
			manager.mu.Unlock()
		case message := <-manager.broadcast:
			manager.mu.Lock()
			for conn := range manager.clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Websocket 错误: %v", err)
					conn.Close()
					delete(manager.clients, conn)
				}
			}
			manager.mu.Unlock()
		}
	}
}

// Handler 返回用于 WebSocket 连接的 Gin 处理器
func (manager *WebSocketManager) Handler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("升级 websocket 失败: %v", err)
		return
	}
	manager.register <- conn

	// 保持连接活跃/读取循环
	// 即使我们只是推送数据，也需要读取以处理关闭帧
	go func() {
		defer func() {
			manager.unregister <- conn
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}

// Broadcast 向所有连接的客户端发送消息
func (manager *WebSocketManager) Broadcast(msg interface{}) {
	// 如果还不是字节数组，则序列化为 JSON
}

// SendBytes 允许外部调用者向所有客户端发送原始字节数据
func (manager *WebSocketManager) SendBytes(data []byte) {
	select {
	case manager.broadcast <- data:
	case <-time.After(1 * time.Second):
		log.Println("广播通道已满，丢弃消息")
	}
}
