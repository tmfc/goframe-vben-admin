// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package menu

import (
	"context"

	"backend/api/menu/v1"
)

type IMenuV1 interface {
	CreateMenu(ctx context.Context, req *v1.CreateMenuReq) (res *v1.CreateMenuRes, err error)
	GetMenu(ctx context.Context, req *v1.GetMenuReq) (res *v1.GetMenuRes, err error)
	UpdateMenu(ctx context.Context, req *v1.UpdateMenuReq) (res *v1.UpdateMenuRes, err error)
	DeleteMenu(ctx context.Context, req *v1.DeleteMenuReq) (res *v1.DeleteMenuRes, err error)
	GetMenuList(ctx context.Context, req *v1.GetMenuListReq) (res *v1.GetMenuListRes, err error)
}
