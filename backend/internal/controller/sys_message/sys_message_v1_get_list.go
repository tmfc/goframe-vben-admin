package sys_message

import (
	"context"

	"backend/api/sys_message/v1"
	"backend/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

func (c *ControllerV1) GetList(ctx context.Context, req *v1.GetListReq) (res *v1.GetListRes, err error) {
	token, err := service.ResolveAccessToken(ctx, "")
	if err != nil {
		return nil, err
	}
	claims, err := service.ParseAccessToken(token)
	if err != nil {
		return nil, err
	}
	userId := gconv.Int64(claims["id"])

	list, total, err := service.SysMessage().GetList(ctx, userId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	return &v1.GetListRes{
		List:  list,
		Total: total,
	}, nil
}