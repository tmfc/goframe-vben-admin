// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysRolePermission is the golang structure of table sys_role_permission for DAO operations like Where/Data.
type SysRolePermission struct {
	g.Meta       `orm:"table:sys_role_permission, do:true"`
	RoleId       any         //
	PermissionId any         //
	CreatedAt    *gtime.Time //
	UpdatedAt    *gtime.Time //
	Scope        any         //
}
