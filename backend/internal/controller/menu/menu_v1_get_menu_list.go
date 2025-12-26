package menu

import (
	"context"

	"backend/internal/model"
	"backend/internal/service"

	"backend/api/menu/v1"
)

func (c *ControllerV1) GetMenuList(ctx context.Context, req *v1.GetMenuListReq) (res *v1.GetMenuListRes, err error) {
	res = &v1.GetMenuListRes{}
	out, err := service.Menu().GetMenuList(ctx, model.SysMenuGetListIn{
		Name:   req.Name,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	res.List = out.List
	res.Total = out.Total
	return res, nil
}
