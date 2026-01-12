// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysMessageRead is the golang structure of table sys_message_read for DAO operations like Where/Data.
type SysMessageRead struct {
	g.Meta    `orm:"table:sys_message_read, do:true"`
	Id        any         //
	MessageId any         //
	UserId    any         //
	ReadAt    *gtime.Time //
}
