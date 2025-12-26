// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package menu

import (
	"context"

	"backend/api/menu/v1"
)

// IMenuV1 defines the menu controller interface.
type IMenuV1 interface {
	All(ctx context.Context, req *v1.MenuAllReq) (res v1.MenuAllRes, err error)
}
