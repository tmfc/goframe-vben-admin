package menu

import (
	"context"

	"backend/internal/service"

	"backend/api/menu/v1"
)

func (c *ControllerV1) GetMenu(ctx context.Context, req *v1.GetMenuReq) (res *v1.GetMenuRes, err error) {
	res = &v1.GetMenuRes{}
	res.SysMenuGetOut, err = service.Menu().GetMenu(ctx, req.ID)
	return res, err
}
