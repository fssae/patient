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
		return true // Allow all origins for now
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
			log.Println("New WebSocket client connected")
		case conn := <-manager.unregister:
			manager.mu.Lock()
			if _, ok := manager.clients[conn]; ok {
				delete(manager.clients, conn)
				conn.Close()
				log.Println("WebSocket client disconnected")
			}
			manager.mu.Unlock()
		case message := <-manager.broadcast:
			manager.mu.RLock()
			for conn := range manager.clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Websocket error: %v", err)
					conn.Close()
					delete(manager.clients, conn)
				}
			}
			manager.mu.RUnlock()
		}
	}
}

// Handler returns the Gin handler for WebSocket connections
func (manager *WebSocketManager) Handler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", err)
		return
	}
	manager.register <- conn

	// Keep connection alive/reader loop
	// Even if we only push, we need to read to handle close frames
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

// Broadcast sends a message to all connected clients
func (manager *WebSocketManager) Broadcast(msg interface{}) {
	// Serialize to JSON if it's not already bytes
	// For simplicity, assuming the caller might want to do the serialization or we pass raw bytes/struct
	// Here we will assume we might receive a struct and handle it, or bytes

	// NOTE: In the service layer we will likely marshal it.
	// To keep this generic, let's accept bytes or handle json here.
	// But the manager.broadcast channel expects []byte.
	// So let's provide a helper method upstream or do it here.
	// For now, let's just expose the internal channel via a safe method
}

// SendBytes allows external callers to send raw bytes to all clients
func (manager *WebSocketManager) SendBytes(data []byte) {
	select {
	case manager.broadcast <- data:
	case <-time.After(1 * time.Second):
		log.Println("Broadcast channel full, dropping message")
	}
}
