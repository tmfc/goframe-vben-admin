package menu

import (
	"context"

	"backend/internal/service"

	v1 "backend/api/menu/v1"
)

// GenerateButtons handles batch generation of default buttons under a menu.
func (c *ControllerV1) GenerateButtons(ctx context.Context, req *v1.GenerateButtonsReq) (res *v1.GenerateButtonsRes, err error) {
	res = &v1.GenerateButtonsRes{}
	created, skipped, err := service.Menu().GenerateButtons(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	res.Created = created
	res.Skipped = skipped
	return res, nil
}
