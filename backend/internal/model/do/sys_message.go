// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysMessage is the golang structure of table sys_message for DAO operations like Where/Data.
type SysMessage struct {
	g.Meta     `orm:"table:sys_message, do:true"`
	Id         any         //
	TenantId   any         //
	Title      any         //
	Content    any         //
	MsgType    any         //
	SenderId   any         //
	TargetType any         //
	Recipients any         //
	Path       any         //
	Params     any         //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
	JumpType   any         // 0: None, 1: Route, 2: External
}
