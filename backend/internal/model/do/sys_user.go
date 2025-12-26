// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysUser is the golang structure of table sys_user for DAO operations like Where/Data.
type SysUser struct {
	g.Meta     `orm:"table:sys_user, do:true"`
	Id         any         //
	TenantId   any         //
	Username   any         //
	Password   any         //
	RealName   any         //
	Avatar     any         //
	HomePath   any         //
	Status     any         //
	Roles      any         //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
	CreatorId  any         //
	ModifierId any         //
	DeptId     any         //
}
