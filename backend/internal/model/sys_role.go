package model

import (

	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
)

// SysRoleCreateIn is the input for creating a new role.
type SysRoleCreateIn struct {
	Name        string `json:"name" v:"required#名称不能为空"`
	Description string `json:"description"`
	ParentId    uint   `json:"parent_id"`
	Status      uint   `json:"status"`
	CreatorId   uint   `json:"creator_id"`
	ModifierId  uint   `json:"modifier_id"`
	DeptId      uint   `json:"dept_id"`
}

// SysRoleCreateOut is the output for creating a new role.
type SysRoleCreateOut struct {
	Id uint `json:"id"`
}

// SysRoleGetIn is the input for retrieving a role.
type SysRoleGetIn struct {
	Id uint `json:"id" v:"required#ID不能为空"`
}

// SysRoleGetOut is the output for retrieving a role.
type SysRoleGetOut struct {
	*entity.SysRole
}

// SysRoleUpdateIn is the input for updating a role.
type SysRoleUpdateIn struct {
	Id          uint        `json:"id" v:"required#ID不能为空"`
	Name        string      `json:"name" v:"required#名称不能为空"`
	Description string      `json:"description"`
	ParentId    uint        `json:"parent_id"`
	Status      uint        `json:"status"`
	ModifierId  uint        `json:"modifier_id"`
	DeptId      uint        `json:"dept_id"`
	UpdatedAt   *gtime.Time `json:"updated_at"`
}

// SysRoleDeleteIn is the input for deleting a role.
type SysRoleDeleteIn struct {
	Id uint `json:"id" v:"required#ID不能为空"`
}

// AssignRoleToUserIn is the input for assigning a role to a user.
type AssignRoleToUserIn struct {
	UserId    string `json:"userId" v:"required#用户ID不能为空"`
	RoleId    uint   `json:"roleId" v:"required#角色ID不能为空"`
	CreatedBy uint   `json:"createdBy"`
}

// AssignRoleToUserOut is the output for assigning a role to a user.
type AssignRoleToUserOut struct {
	Success bool `json:"success"`
}

// RemoveRoleFromUserIn is the input for removing a role from a user.
type RemoveRoleFromUserIn struct {
	UserId string `json:"userId" v:"required#用户ID不能为空"`
	RoleId uint   `json:"roleId" v:"required#角色ID不能为空"`
}

// RemoveRoleFromUserOut is the output for removing a role from a user.
type RemoveRoleFromUserOut struct {
	Success bool `json:"success"`
}

// GetUserRolesIn is the input for getting user's roles.
type GetUserRolesIn struct {
	UserId string `json:"userId" v:"required#用户ID不能为空"`
}

// GetUserRolesOut is the output for getting user's roles.
type GetUserRolesOut struct {
	Roles []*entity.SysRole `json:"roles"`
}

// AssignRolesToUserIn is the input for assigning multiple roles to a user.
type AssignRolesToUserIn struct {
	UserId    string `json:"userId" v:"required#用户ID不能为空"`
	RoleIds   []uint `json:"roleIds" v:"required#角色ID列表不能为空"`
	CreatedBy uint   `json:"createdBy"`
}

// AssignRolesToUserOut is the output for assigning multiple roles to a user.
type AssignRolesToUserOut struct {
	Success bool `json:"success"`
	Count   int  `json:"count"`
}

// GetUsersByRoleIn is the input for getting users by role.
type GetUsersByRoleIn struct {
	RoleId uint `json:"roleId" v:"required#角色ID不能为空"`
}

// GetUsersByRoleOut is the output for getting users by role.
type GetUsersByRoleOut struct {
	Users []*entity.SysUser `json:"users"`
}
