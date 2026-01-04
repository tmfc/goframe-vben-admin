package sys_permission

import (
	"context"

	"backend/api/sys_permission/v1"
)

type ISysPermissionV1 interface {
	CreatePermission(ctx context.Context, req *v1.CreatePermissionReq) (res *v1.CreatePermissionRes, err error)
	GetPermission(ctx context.Context, req *v1.GetPermissionReq) (res *v1.GetPermissionRes, err error)
	UpdatePermission(ctx context.Context, req *v1.UpdatePermissionReq) (res *v1.UpdatePermissionRes, err error)
	DeletePermission(ctx context.Context, req *v1.DeletePermissionReq) (res *v1.DeletePermissionRes, err error)
	GetPermissionsByUser(ctx context.Context, req *v1.GetPermissionsByUserReq) (res *v1.GetPermissionsByUserRes, err error)
}
