package ws

import (
	"context"
	"net/http"

	"backend/api/ws/v1"
	"backend/internal/logic/ws"
	"backend/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now, can be restricted in production
	},
}

func (c *ControllerV1) Connect(ctx context.Context, req *v1.ConnectReq) (res *v1.ConnectRes, err error) {
	r := g.RequestFromCtx(ctx)

	// 1. Authenticate using token from query/req
	claims, err := service.NewJWTTokenService().ParseAccessToken(req.Token)
	if err != nil {
		return nil, err
	}

	userId := int64(claims["id"].(float64))

	// 2. Upgrade to WebSocket
	conn, err := upgrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		g.Log().Error(ctx, "WebSocket upgrade failed", err)
		return nil, err
	}

	// 3. Register client in Hub
	hub := ws.GetHub()
	hub.Register(userId, conn)

	// 4. Handle connection closure and heartbeat/read loop
	go func() {
		defer func() {
			hub.Unregister(userId, conn)
			conn.Close()
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// Client disconnected
				break
			}
			// We don't expect messages from client for now, just keep connection open
		}
	}()

	// GoFrame controllers returning nil for both res and err will stop further processing for the request
	// which is what we want for WebSocket after hijacking the connection via Upgrade.
	return nil, nil
}