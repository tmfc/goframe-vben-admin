// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysMenu is the golang structure of table sys_menu for DAO operations like Where/Data.
type SysMenu struct {
	g.Meta         `orm:"table:sys_menu, do:true"`
	Id             any         //
	TenantId       any         //
	ParentId       any         //
	Name           any         //
	Path           any         //
	Component      any         //
	Icon           any         //
	Order          any         //
	Type           any         //
	Visible        any         //
	Status         any         //
	PermissionCode any         //
	Meta           any         //
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	DeletedAt      *gtime.Time //
	CreatorId      any         //
	ModifierId     any         //
	DeptId         any         //
}
