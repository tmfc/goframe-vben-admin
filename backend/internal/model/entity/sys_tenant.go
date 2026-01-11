// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysTenant is the golang structure for table sys_tenant.
type SysTenant struct {
	Id             int64       `json:"id"             orm:"id"              ` //
	Name           string      `json:"name"           orm:"name"            ` //
	Status         int         `json:"status"         orm:"status"          ` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` //
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      ` //
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"      ` //
	Code           string      `json:"code"           orm:"code"            ` //
	ContactName    string      `json:"contactName"    orm:"contact_name"    ` //
	ContactMobile  string      `json:"contactMobile"  orm:"contact_mobile"  ` //
	StartAt        *gtime.Time `json:"startAt"        orm:"start_at"        ` //
	ExpireAt       *gtime.Time `json:"expireAt"       orm:"expire_at"       ` //
	PackageVersion string      `json:"packageVersion" orm:"package_version" ` //
	Domain         string      `json:"domain"         orm:"domain"          ` //
}
