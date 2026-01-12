package sys_message

import (
	"context"

	"backend/api/sys_message/v1"
)

type ISysMessageV1 interface {
	GetList(ctx context.Context, req *v1.GetListReq) (res *v1.GetListRes, err error)
	SetRead(ctx context.Context, req *v1.SetReadReq) (res *v1.SetReadRes, err error)
	SetAllRead(ctx context.Context, req *v1.SetAllReadReq) (res *v1.SetAllReadRes, err error)
}
