package sys_message

import (
	"context"

	"backend/api/sys_message/v1"
	"backend/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) SetRead(ctx context.Context, req *v1.SetReadReq) (res *v1.SetReadRes, err error) {
	token, err := service.ResolveAccessToken(ctx, "")
	if err != nil {
		return nil, err
	}
	claims, err := service.ParseAccessToken(token)
	if err != nil {
		return nil, err
	}
	userId := gconv.Int64(claims["id"])

	err = service.SysMessage().SetRead(ctx, userId, req.IDs)
	if err != nil {
		return nil, err
	}

	return &v1.SetReadRes{}, nil
}