// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysRole is the golang structure of table sys_role for DAO operations like Where/Data.
type SysRole struct {
	g.Meta      `orm:"table:sys_role, do:true"`
	Id          any         //
	Name        any         //
	Description any         //
	ParentId    any         //
	Status      any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	CreatorId   any         //
	ModifierId  any         //
	DeptId      any         //
}
