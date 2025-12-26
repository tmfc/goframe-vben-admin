package model

import (

	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
)

// SysPermissionCreateIn is the input for creating a new permission.
type SysPermissionCreateIn struct {
	Name        string `json:"name" v:"required#名称不能为空"`
	Description string `json:"description"`
	ParentId    uint   `json:"parent_id"`
	Status      uint   `json:"status"`
	CreatorId   uint   `json:"creator_id"`
	ModifierId  uint   `json:"modifier_id"`
	DeptId      uint   `json:"dept_id"`
}

// SysPermissionCreateOut is the output for creating a new permission.
type SysPermissionCreateOut struct {
	Id uint `json:"id"`
}

// SysPermissionGetIn is the input for retrieving a permission.
type SysPermissionGetIn struct {
	Id uint `json:"id" v:"required#ID不能为空"`
}

// SysPermissionGetOut is the output for retrieving a permission.
type SysPermissionGetOut struct {
	*entity.SysPermission
}

// SysPermissionUpdateIn is the input for updating a permission.
type SysPermissionUpdateIn struct {
	Id          uint        `json:"id" v:"required#ID不能为空"`
	Name        string      `json:"name" v:"required#名称不能为空"`
	Description string      `json:"description"`
	ParentId    uint        `json:"parent_id"`
	Status      uint        `json:"status"`
	ModifierId  uint        `json:"modifier_id"`
	DeptId      uint        `json:"dept_id"`
	UpdatedAt   *gtime.Time `json:"updated_at"`
}

// SysPermissionDeleteIn is the input for deleting a permission.
type SysPermissionDeleteIn struct {
	Id uint `json:"id" v:"required#ID不能为空"`
}
