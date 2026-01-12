// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysMessage is the golang structure for table sys_message.
type SysMessage struct {
	Id         int64       `json:"id"         orm:"id"          ` //
	TenantId   int64       `json:"tenantId"   orm:"tenant_id"   ` //
	Title      string      `json:"title"      orm:"title"       ` //
	Content    string      `json:"content"    orm:"content"     ` //
	MsgType    int         `json:"msgType"    orm:"msg_type"    ` //
	SenderId   int64       `json:"senderId"   orm:"sender_id"   ` //
	TargetType int         `json:"targetType" orm:"target_type" ` //
	Recipients string      `json:"recipients" orm:"recipients"  ` //
	Path       string      `json:"path"       orm:"path"        ` //
	Params     string      `json:"params"     orm:"params"      ` //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` //
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` //
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  ` //
	JumpType   int         `json:"jumpType"   orm:"jump_type"   ` // 0: None, 1: Route, 2: External
}
