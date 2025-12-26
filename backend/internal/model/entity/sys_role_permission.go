// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysRolePermission is the golang structure for table sys_role_permission.
type SysRolePermission struct {
	RoleId       int64       `json:"roleId"       orm:"role_id"       ` //
	PermissionId int64       `json:"permissionId" orm:"permission_id" ` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` //
	Scope        string      `json:"scope"        orm:"scope"         ` //
}
