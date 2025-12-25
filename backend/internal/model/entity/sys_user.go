// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysUser is the golang structure for table sys_user.
type SysUser struct {
	Id        string      `json:"id"        orm:"id"         description:""` //
	TenantId  string      `json:"tenantId"  orm:"tenant_id"  description:""` //
	Username  string      `json:"username"  orm:"username"   description:""` //
	Password  string      `json:"password"  orm:"password"   description:""` //
	RealName  string      `json:"realName"  orm:"real_name"  description:""` //
	Avatar    string      `json:"avatar"    orm:"avatar"     description:""` //
	HomePath  string      `json:"homePath"  orm:"home_path"  description:""` //
	Status    int         `json:"status"    orm:"status"     description:""` //
	Roles     string      `json:"roles"     orm:"roles"      description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:""` //
}
