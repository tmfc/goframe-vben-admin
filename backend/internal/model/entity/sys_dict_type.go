// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDictType is the golang structure for table sys_dict_type.
type SysDictType struct {
	Id          int64       `json:"id"          orm:"id"          ` //
	TenantId    int64       `json:"tenantId"    orm:"tenant_id"   ` //
	TypeCode    string      `json:"typeCode"    orm:"type_code"   ` //
	TypeName    string      `json:"typeName"    orm:"type_name"   ` //
	Description string      `json:"description" orm:"description" ` //
	IsSystem    bool        `json:"isSystem"    orm:"is_system"   ` //
	Status      int         `json:"status"      orm:"status"      ` //
	Sort        int         `json:"sort"        orm:"sort"        ` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  ` //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  ` //
	CreatorId   int64       `json:"creatorId"   orm:"creator_id"  ` //
	ModifierId  int64       `json:"modifierId"  orm:"modifier_id" ` //
	DeptId      int64       `json:"deptId"      orm:"dept_id"     ` //
}
