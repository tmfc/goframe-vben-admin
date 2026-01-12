package ws

import (
	"context"

	"backend/api/ws/v1"
)

type IWsV1 interface {
	Connect(ctx context.Context, req *v1.ConnectReq) (res *v1.ConnectRes, err error)
}
