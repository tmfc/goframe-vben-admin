package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ConnectReq struct {
	g.Meta `path:"/ws" method:"get" summary:"WebSocket connection" tags:"WebSocket"`
	Token  string `json:"token" v:"required#Token is required"`
}

type ConnectRes struct{}
