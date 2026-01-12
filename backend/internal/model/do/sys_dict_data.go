// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDictData is the golang structure of table sys_dict_data for DAO operations like Where/Data.
type SysDictData struct {
	g.Meta      `orm:"table:sys_dict_data, do:true"`
	Id          any         //
	TenantId    any         //
	DictTypeId  any         //
	Label       any         //
	LabelI18n   any         //
	Value       any         //
	Description any         //
	Color       any         //
	Icon        any         //
	CssClass    any         //
	Status      any         //
	Sort        any         //
	IsDefault   any         //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	CreatorId   any         //
	ModifierId  any         //
	DeptId      any         //
}
