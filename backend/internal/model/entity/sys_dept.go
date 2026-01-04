// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDept is the golang structure for table sys_dept.
type SysDept struct {
	Id         string      `json:"id"         orm:"id"          ` //
	TenantId   string      `json:"tenantId"   orm:"tenant_id"   ` //
	ParentId   string      `json:"parentId"   orm:"parent_id"   ` //
	Name       string      `json:"name"       orm:"name"        ` //
	Order      int         `json:"order"      orm:"order"       ` //
	Status     int         `json:"status"     orm:"status"      ` //
	CreatorId  int64       `json:"creatorId"  orm:"creator_id"  ` //
	ModifierId int64       `json:"modifierId" orm:"modifier_id" ` //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` //
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` //
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  ` //
}
