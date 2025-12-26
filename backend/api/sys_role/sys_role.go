package sys_role

import (
	"context"

	"backend/api/sys_role/v1"
)

type ISysRoleV1 interface {
	CreateRole(ctx context.Context, req *v1.CreateRoleReq) (res *v1.CreateRoleRes, err error)
	GetRole(ctx context.Context, req *v1.GetRoleReq) (res *v1.GetRoleRes, err error)
	UpdateRole(ctx context.Context, req *v1.UpdateRoleReq) (res *v1.UpdateRoleRes, err error)
	DeleteRole(ctx context.Context, req *v1.DeleteRoleReq) (res *v1.DeleteRoleRes, err error)
}
