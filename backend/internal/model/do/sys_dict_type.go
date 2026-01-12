// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDictType is the golang structure of table sys_dict_type for DAO operations like Where/Data.
type SysDictType struct {
	g.Meta      `orm:"table:sys_dict_type, do:true"`
	Id          any         //
	TenantId    any         //
	TypeCode    any         //
	TypeName    any         //
	Description any         //
	IsSystem    any         //
	Status      any         //
	Sort        any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	CreatorId   any         //
	ModifierId  any         //
	DeptId      any         //
}
