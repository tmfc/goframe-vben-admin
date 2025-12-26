package menu

import (
	"context"

	"backend/internal/service"

	"backend/api/menu/v1"
)

func (c *ControllerV1) CreateMenu(ctx context.Context, req *v1.CreateMenuReq) (res *v1.CreateMenuRes, err error) {
	res = &v1.CreateMenuRes{}
	id, err := service.Menu().CreateMenu(ctx, req.SysMenuCreateIn)
	if err != nil {
		return nil, err
	}
	res.Id = id
	return res, nil
}
