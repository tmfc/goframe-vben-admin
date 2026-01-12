package ws

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gorilla/websocket"
)

type Hub struct {
	// Map userId -> []*websocket.Conn (support multiple tabs for one user)
	clients map[int64][]*websocket.Conn
	mu      sync.RWMutex
}

var globalHub = &Hub{
	clients: make(map[int64][]*websocket.Conn),
}

func GetHub() *Hub {
	return globalHub
}

func (h *Hub) Register(userId int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userId] = append(h.clients[userId], conn)
	g.Log().Debugf(context.Background(), "User %d registered WebSocket connection. Total user connections: %d", userId, len(h.clients[userId]))
}

func (h *Hub) Unregister(userId int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.clients[userId]; ok {
		for i, c := range conns {
			if c == conn {
				h.clients[userId] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		if len(h.clients[userId]) == 0 {
			delete(h.clients, userId)
		}
	}
}

func (h *Hub) SendToUser(userId int64, message interface{}) {
	h.mu.RLock()
	conns, ok := h.clients[userId]
	h.mu.RUnlock()

	if !ok {
		return
	}

	for _, conn := range conns {
		err := conn.WriteJSON(message)
		if err != nil {
			g.Log().Error(context.Background(), "Failed to send WebSocket message to user", userId, err)
			// Connection might be closed, it will be cleaned up in the read loop or Heartbeat
		}
	}
}

func (h *Hub) Broadcast(message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for userId, conns := range h.clients {
		for _, conn := range conns {
			err := conn.WriteJSON(message)
			if err != nil {
				g.Log().Error(context.Background(), "Failed to broadcast WebSocket message to user", userId, err)
			}
		}
	}
}
