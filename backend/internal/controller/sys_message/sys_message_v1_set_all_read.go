package sys_message

import (
	"context"

	"backend/api/sys_message/v1"
	"backend/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) SetAllRead(ctx context.Context, req *v1.SetAllReadReq) (res *v1.SetAllReadRes, err error) {
	token, err := service.ResolveAccessToken(ctx, "")
	if err != nil {
		return nil, err
	}
	claims, err := service.ParseAccessToken(token)
	if err != nil {
		return nil, err
	}
	userId := gconv.Int64(claims["id"])

	// For SetAllRead, we need to find all unread message IDs for this user
	// But actually, we can implement it more efficiently in logic layer
	// For now, let's keep controller simple and add SetAllRead to service
	
	// I'll add SetAllRead to service and logic
	err = service.SysMessage().SetAllRead(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &v1.SetAllReadRes{}, nil
}