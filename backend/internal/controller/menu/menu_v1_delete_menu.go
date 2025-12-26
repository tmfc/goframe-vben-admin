package menu

import (
	"context"

	"backend/internal/service"

	"backend/api/menu/v1"
)

func (c *ControllerV1) DeleteMenu(ctx context.Context, req *v1.DeleteMenuReq) (res *v1.DeleteMenuRes, err error) {
	err = service.Menu().DeleteMenu(ctx, req.ID)
	return nil, err
}
