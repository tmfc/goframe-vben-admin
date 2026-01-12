package model

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysDictDataCreateIn is the input for creating a new dict data.
type SysDictDataCreateIn struct {
	DictTypeId  int64             `json:"dictTypeId" v:"required#字典类型不能为空"`
	Label       string            `json:"label" v:"required#标签不能为空"`
	LabelI18n   map[string]string `json:"labelI18n"`
	Value       string            `json:"value" v:"required#值不能为空"`
	Description string            `json:"description"`
	Color       string            `json:"color"`
	Icon        string            `json:"icon"`
	CssClass    string            `json:"cssClass"`
	Status      int               `json:"status"`
	Sort        int               `json:"sort"`
	IsDefault   bool              `json:"isDefault"`
	CreatorId   int64             `json:"creatorId"`
	ModifierId  int64             `json:"modifierId"`
	DeptId      int64             `json:"deptId"`
}

// SysDictDataCreateOut is the output for creating a new dict data.
type SysDictDataCreateOut struct {
	Id int64 `json:"id"`
}

// SysDictDataGetIn is the input for retrieving a dict data.
type SysDictDataGetIn struct {
	Id int64 `json:"id" v:"required#ID不能为空"`
}

// SysDictDataDetail is the output detail for dict data.
type SysDictDataDetail struct {
	Id          int64             `json:"id"`
	TenantId    int64             `json:"tenantId"`
	DictTypeId  int64             `json:"dictTypeId"`
	Label       string            `json:"label"`
	LabelI18n   map[string]string `json:"labelI18n,omitempty"`
	Value       string            `json:"value"`
	Description string            `json:"description"`
	Color       string            `json:"color"`
	Icon        string            `json:"icon"`
	CssClass    string            `json:"cssClass"`
	Status      int               `json:"status"`
	Sort        int               `json:"sort"`
	IsDefault   bool              `json:"isDefault"`
	CreatedAt   *gtime.Time       `json:"createdAt"`
	UpdatedAt   *gtime.Time       `json:"updatedAt"`
	CreatorId   int64             `json:"creatorId"`
	ModifierId  int64             `json:"modifierId"`
	DeptId      int64             `json:"deptId"`
}

// SysDictDataGetOut is the output for retrieving a dict data.
type SysDictDataGetOut struct {
	*SysDictDataDetail
}

// SysDictDataUpdateIn is the input for updating a dict data.
type SysDictDataUpdateIn struct {
	Id          int64             `json:"id" v:"required#ID不能为空"`
	DictTypeId  int64             `json:"dictTypeId" v:"required#字典类型不能为空"`
	Label       string            `json:"label" v:"required#标签不能为空"`
	LabelI18n   map[string]string `json:"labelI18n"`
	Value       string            `json:"value" v:"required#值不能为空"`
	Description string            `json:"description"`
	Color       string            `json:"color"`
	Icon        string            `json:"icon"`
	CssClass    string            `json:"cssClass"`
	Status      int               `json:"status"`
	Sort        int               `json:"sort"`
	IsDefault   bool              `json:"isDefault"`
	ModifierId  int64             `json:"modifierId"`
	UpdatedAt   *gtime.Time       `json:"updatedAt"`
}

// SysDictDataDeleteIn is the input for deleting a dict data.
type SysDictDataDeleteIn struct {
	Id int64 `json:"id" v:"required#ID不能为空"`
}

// SysDictDataListIn is the input for listing dict data.
type SysDictDataListIn struct {
	Page       int    `json:"page" d:"1"`
	PageSize   int    `json:"pageSize" d:"10"`
	DictTypeId int64  `json:"dictTypeId"`
	TypeCode   string `json:"typeCode"`
	Label      string `json:"label"`
	Value      string `json:"value"`
	Status     string `json:"status"`
}

// SysDictDataListOut is the output for listing dict data.
type SysDictDataListOut struct {
	Items []*SysDictDataDetail `json:"items"`
	Total int                  `json:"total"`
}
