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
	AssignUsersToRole(ctx context.Context, req *v1.AssignUsersToRoleReq) (res *v1.AssignUsersToRoleRes, err error)
	RemoveUsersFromRole(ctx context.Context, req *v1.RemoveUsersFromRoleReq) (res *v1.RemoveUsersFromRoleRes, err error)
	GetRoleUsers(ctx context.Context, req *v1.GetRoleUsersReq) (res *v1.GetRoleUsersRes, err error)
	AssignPermissionsToRole(ctx context.Context, req *v1.AssignPermissionsToRoleReq) (res *v1.AssignPermissionsToRoleRes, err error)
	RemovePermissionsFromRole(ctx context.Context, req *v1.RemovePermissionsFromRoleReq) (res *v1.RemovePermissionsFromRoleRes, err error)
	GetRolePermissions(ctx context.Context, req *v1.GetRolePermissionsReq) (res *v1.GetRolePermissionsRes, err error)
}
