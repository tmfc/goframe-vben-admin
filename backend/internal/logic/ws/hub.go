package ws

import (
	"context"
	"sync"

	"backend/internal/logic/sys_message/mq"
	"backend/internal/model"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[int64][]*websocket.Conn
	mu      sync.RWMutex
	once    sync.Once
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

// StartMessageConsumer starts a Redis subscription to consume messages and push to WebSocket
func (h *Hub) StartMessageConsumer(ctx context.Context) {
	h.once.Do(func() {
		go mq.Subscribe(ctx, "sys_message", func(ctx context.Context, msg string) error {
			j, err := gjson.DecodeToJson(msg)
			if err != nil {
				return err
			}

			// We use MessageCreateInput as the wire format from MQ
			var in model.MessageCreateInput
			if err := j.Scan(&in); err != nil {
				return err
			}

			if in.ReceiverID != nil {
				h.SendToUser(*in.ReceiverID, in)
			} else {
				h.Broadcast(in)
			}
			return nil
		})
	})
}