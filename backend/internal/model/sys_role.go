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
