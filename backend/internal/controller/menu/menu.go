package menu

import (
	"context"

	"backend/api/menu"
	"backend/api/menu/v1"
	"backend/internal/service"
)

// ControllerV1 handles menu endpoints.
type ControllerV1 struct{}

// NewV1 creates the menu controller instance.
func NewV1() menu.IMenuV1 {
	return &ControllerV1{}
}

// MenuAll returns the full menu list for the current user.
func (c *ControllerV1) MenuAll(ctx context.Context, req *v1.MenuAllReq) (res *v1.MenuAllRes, err error) {
	out, err := service.Menu().All(ctx)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
