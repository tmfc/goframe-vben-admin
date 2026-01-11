// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysUserRole is the golang structure for table sys_user_role.
type SysUserRole struct {
	Id        int64       `json:"id"        orm:"id"         ` //
	TenantId  int64       `json:"tenantId"  orm:"tenant_id"  ` //
	UserId    int64       `json:"userId"    orm:"user_id"    ` //
	RoleId    int64       `json:"roleId"    orm:"role_id"    ` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	CreatedBy int64       `json:"createdBy" orm:"created_by" ` //
}
