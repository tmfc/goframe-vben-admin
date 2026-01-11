// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysPermission is the golang structure for table sys_permission.
type SysPermission struct {
	Id          int64       `json:"id"          orm:"id"          ` //
	Name        string      `json:"name"        orm:"name"        ` //
	Description string      `json:"description" orm:"description" ` //
	ParentId    int64       `json:"parentId"    orm:"parent_id"   ` //
	Status      int         `json:"status"      orm:"status"      ` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  ` //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  ` //
	CreatorId   int64       `json:"creatorId"   orm:"creator_id"  ` //
	ModifierId  int64       `json:"modifierId"  orm:"modifier_id" ` //
	DeptId      int64       `json:"deptId"      orm:"dept_id"     ` //
	TenantId    int64       `json:"tenantId"    orm:"tenant_id"   ` //
}
