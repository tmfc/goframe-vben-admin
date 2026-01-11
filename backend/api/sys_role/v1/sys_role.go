package v1

import (
	"backend/internal/model"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

// CreateRoleReq defines the request structure for creating a new role.
type CreateRoleReq struct {
	g.Meta `path:"/sys-role" method:"post" summary:"Create a new role" tags:"System Role"`
	model.SysRoleCreateIn
}

// CreateRoleRes defines the response structure for creating a new role.
type CreateRoleRes struct {
	Id uint `json:"id"`
}

// GetRoleReq defines the request structure for retrieving a role.
type GetRoleReq struct {
	g.Meta `path:"/sys-role/{id}" method:"get" summary:"Retrieve a role by ID" tags:"System Role"`
	ID     uint `json:"id" v:"required#ID不能为空"`
}

// GetRoleRes defines the response structure for retrieving a role.
type GetRoleRes struct {
	*model.SysRoleGetOut
}

// GetRoleListReq defines the request structure for listing roles.
type GetRoleListReq struct {
	g.Meta   `path:"/sys-role" method:"get" summary:"List roles" tags:"System Role"`
	Page     int    `json:"page" d:"1"`
	PageSize int    `json:"pageSize" d:"10"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

// GetRoleListRes defines the response structure for listing roles.
type GetRoleListRes struct {
	Items []*entity.SysRole `json:"items"`
	Total int               `json:"total"`
}

// UpdateRoleReq defines the request structure for updating a role.
type UpdateRoleReq struct {
	g.Meta `path:"/sys-role/{id}" method:"put" summary:"Update a role by ID" tags:"System Role"`
	ID     uint `json:"id" v:"required#ID不能为空"`
	model.SysRoleUpdateIn
}

// UpdateRoleRes defines the response structure for updating a role.
type UpdateRoleRes struct{}

// DeleteRoleReq defines the request structure for deleting a role.
type DeleteRoleReq struct {
	g.Meta `path:"/sys-role/{id}" method:"delete" summary:"Delete a role by ID" tags:"System Role"`
	ID     uint `json:"id" v:"required#ID不能为空"`
}

// DeleteRoleRes defines the response structure for deleting a role.
type DeleteRoleRes struct{}

// AssignUsersToRoleReq defines the request structure for assigning users to a role.
type AssignUsersToRoleReq struct {
	g.Meta    `path:"/sys-role/{id}/assign-users" method:"post" summary:"Assign users to role" tags:"System Role"`
	ID        uint     `json:"id" v:"required#ID不能为空"`
	UserIds   []int64 `json:"userIds" v:"required#用户ID列表不能为空"`
	CreatedBy uint     `json:"createdBy"`
}

// AssignUsersToRoleRes defines the response structure for assigning users to a role.
type AssignUsersToRoleRes struct {
	Success bool `json:"success"`
	Count   int  `json:"count"`
}

// RemoveUsersFromRoleReq defines the request structure for removing users from a role.
type RemoveUsersFromRoleReq struct {
	g.Meta  `path:"/sys-role/{id}/remove-users" method:"post" summary:"Remove users from role" tags:"System Role"`
	ID      uint     `json:"id" v:"required#ID不能为空"`
	UserIds []int64 `json:"userIds" v:"required#用户ID列表不能为空"`
}

// RemoveUsersFromRoleRes defines the response structure for removing users from a role.
type RemoveUsersFromRoleRes struct {
	Success bool `json:"success"`
	Count   int  `json:"count"`
}

// GetRoleUsersReq defines the request structure for retrieving users of a role.
type GetRoleUsersReq struct {
	g.Meta `path:"/sys-role/{id}/users" method:"get" summary:"Get users by role" tags:"System Role"`
	ID     uint `json:"id" v:"required#ID不能为空"`
}

// GetRoleUsersRes defines the response structure for retrieving users of a role.
type GetRoleUsersRes struct {
	*model.GetUsersByRoleOut
}

// AssignPermissionsToRoleReq defines the request structure for assigning permissions to a role.
type AssignPermissionsToRoleReq struct {
	g.Meta        `path:"/sys-role/{id}/assign-permissions" method:"post" summary:"Assign permissions to role" tags:"System Role"`
	ID            uint   `json:"id" v:"required#ID不能为空"`
	PermissionIds []uint `json:"permissionIds" v:"required#权限ID列表不能为空"`
}

// AssignPermissionsToRoleRes defines the response structure for assigning permissions to a role.
type AssignPermissionsToRoleRes struct {
	Success bool `json:"success"`
}

// RemovePermissionsFromRoleReq defines the request structure for removing permissions from a role.
type RemovePermissionsFromRoleReq struct {
	g.Meta        `path:"/sys-role/{id}/remove-permissions" method:"post" summary:"Remove permissions from role" tags:"System Role"`
	ID            uint   `json:"id" v:"required#ID不能为空"`
	PermissionIds []uint `json:"permissionIds" v:"required#权限ID列表不能为空"`
}

// RemovePermissionsFromRoleRes defines the response structure for removing permissions from a role.
type RemovePermissionsFromRoleRes struct {
	Success bool `json:"success"`
	Count   int  `json:"count"`
}

// GetRolePermissionsReq defines the request structure for retrieving permissions of a role.
type GetRolePermissionsReq struct {
	g.Meta `path:"/sys-role/{id}/permissions" method:"get" summary:"Get permissions by role" tags:"System Role"`
	ID     uint `json:"id" v:"required#ID不能为空"`
}

// GetRolePermissionsRes defines the response structure for retrieving permissions of a role.
type GetRolePermissionsRes struct {
	*model.SysRolePermissionOut
}
