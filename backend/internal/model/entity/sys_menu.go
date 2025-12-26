// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysMenu is the golang structure for table sys_menu.
type SysMenu struct {
	Id             string      `json:"id"             orm:"id"              ` //
	TenantId       string      `json:"tenantId"       orm:"tenant_id"       ` //
	ParentId       string      `json:"parentId"       orm:"parent_id"       ` //
	Name           string      `json:"name"           orm:"name"            ` //
	Path           string      `json:"path"           orm:"path"            ` //
	Component      string      `json:"component"      orm:"component"       ` //
	Icon           string      `json:"icon"           orm:"icon"            ` //
	Order          int         `json:"order"          orm:"order"           ` //
	Type           string      `json:"type"           orm:"type"            ` //
	Visible        int         `json:"visible"        orm:"visible"         ` //
	Status         int         `json:"status"         orm:"status"          ` //
	PermissionCode string      `json:"permissionCode" orm:"permission_code" ` //
	Meta           string      `json:"meta"           orm:"meta"            ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      ` //
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"      ` //
	CreatorId      int64       `json:"creatorId"      orm:"creator_id"      ` //
	ModifierId     int64       `json:"modifierId"     orm:"modifier_id"     ` //
	DeptId         int64       `json:"deptId"         orm:"dept_id"         ` //
}
