package menu

import (
	"context"

	"backend/internal/service"

	"backend/api/menu/v1"
)

func (c *ControllerV1) UpdateMenu(ctx context.Context, req *v1.UpdateMenuReq) (res *v1.UpdateMenuRes, err error) {
	err = service.Menu().UpdateMenu(ctx, req.SysMenuUpdateIn)
	return nil, err
}
