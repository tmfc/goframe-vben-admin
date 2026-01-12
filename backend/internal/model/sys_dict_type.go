package model

import (
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
)

// SysDictTypeCreateIn is the input for creating a new dict type.
type SysDictTypeCreateIn struct {
	TypeCode    string `json:"typeCode" v:"required#类型编码不能为空"`
	TypeName    string `json:"typeName" v:"required#类型名称不能为空"`
	Description string `json:"description"`
	IsSystem    bool   `json:"isSystem"`
	Status      int    `json:"status"`
	Sort        int    `json:"sort"`
	CreatorId   int64  `json:"creatorId"`
	ModifierId  int64  `json:"modifierId"`
	DeptId      int64  `json:"deptId"`
}

// SysDictTypeCreateOut is the output for creating a new dict type.
type SysDictTypeCreateOut struct {
	Id int64 `json:"id"`
}

// SysDictTypeGetIn is the input for retrieving a dict type.
type SysDictTypeGetIn struct {
	Id int64 `json:"id" v:"required#ID不能为空"`
}

// SysDictTypeGetOut is the output for retrieving a dict type.
type SysDictTypeGetOut struct {
	*entity.SysDictType
}

// SysDictTypeUpdateIn is the input for updating a dict type.
type SysDictTypeUpdateIn struct {
	Id          int64       `json:"id" v:"required#ID不能为空"`
	TypeCode    string      `json:"typeCode" v:"required#类型编码不能为空"`
	TypeName    string      `json:"typeName" v:"required#类型名称不能为空"`
	Description string      `json:"description"`
	IsSystem    bool        `json:"isSystem"`
	Status      int         `json:"status"`
	Sort        int         `json:"sort"`
	ModifierId  int64       `json:"modifierId"`
	UpdatedAt   *gtime.Time `json:"updatedAt"`
}

// SysDictTypeDeleteIn is the input for deleting a dict type.
type SysDictTypeDeleteIn struct {
	Id int64 `json:"id" v:"required#ID不能为空"`
}

// SysDictTypeListIn is the input for listing dict types.
type SysDictTypeListIn struct {
	Page     int    `json:"page" d:"1"`
	PageSize int    `json:"pageSize" d:"10"`
	TypeCode string `json:"typeCode"`
	TypeName string `json:"typeName"`
	Status   string `json:"status"`
}

// SysDictTypeListOut is the output for listing dict types.
type SysDictTypeListOut struct {
	Items []*entity.SysDictType `json:"items"`
	Total int                   `json:"total"`
}
