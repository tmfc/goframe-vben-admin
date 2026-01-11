// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysUser is the golang structure for table sys_user.
type SysUser struct {
	Id         int64       `json:"id"         orm:"id"          ` //
	TenantId   int64       `json:"tenantId"   orm:"tenant_id"   ` //
	Username   string      `json:"username"   orm:"username"    ` //
	Password   string      `json:"password"   orm:"password"    ` //
	RealName   string      `json:"realName"   orm:"real_name"   ` //
	Avatar     string      `json:"avatar"     orm:"avatar"      ` //
	HomePath   string      `json:"homePath"   orm:"home_path"   ` //
	Status     int         `json:"status"     orm:"status"      ` //
	Roles      string      `json:"roles"      orm:"roles"       ` //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` //
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` //
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  ` //
	CreatorId  int64       `json:"creatorId"  orm:"creator_id"  ` //
	ModifierId int64       `json:"modifierId" orm:"modifier_id" ` //
	DeptId     int64       `json:"deptId"     orm:"dept_id"     ` //
	ExtInfo    string      `json:"extInfo"    orm:"ext_info"    ` //
}
