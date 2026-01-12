// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDictData is the golang structure for table sys_dict_data.
type SysDictData struct {
	Id          int64       `json:"id"         orm:"id"`           //
	TenantId    int64       `json:"tenantId"   orm:"tenant_id"`    //
	DictTypeId  int64       `json:"dictTypeId" orm:"dict_type_id"` //
	Label       string      `json:"label"      orm:"label"`        //
	LabelI18n   string      `json:"labelI18n"  orm:"label_i18n"`   //
	Value       string      `json:"value"      orm:"value"`        //
	Description string      `json:"description" orm:"description"` //
	Color       string      `json:"color"      orm:"color"`        //
	Icon        string      `json:"icon"       orm:"icon"`         //
	CssClass    string      `json:"cssClass"   orm:"css_class"`    //
	Status      int         `json:"status"     orm:"status"`       //
	Sort        int         `json:"sort"       orm:"sort"`         //
	IsDefault   bool        `json:"isDefault"  orm:"is_default"`   //
	CreatedAt   *gtime.Time `json:"createdAt"  orm:"created_at"`   //
	UpdatedAt   *gtime.Time `json:"updatedAt"  orm:"updated_at"`   //
	CreatorId   int64       `json:"creatorId"  orm:"creator_id"`   //
	ModifierId  int64       `json:"modifierId" orm:"modifier_id"`  //
	DeptId      int64       `json:"deptId"     orm:"dept_id"`      //
}
