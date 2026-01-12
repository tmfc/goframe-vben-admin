// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysMessageRead is the golang structure for table sys_message_read.
type SysMessageRead struct {
	Id        int64       `json:"id"        orm:"id"         ` //
	MessageId int64       `json:"messageId" orm:"message_id" ` //
	UserId    int64       `json:"userId"    orm:"user_id"    ` //
	ReadAt    *gtime.Time `json:"readAt"    orm:"read_at"    ` //
}
