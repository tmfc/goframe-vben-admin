package v1

import (
	"backend/internal/model"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

// CreateDictTypeReq defines the request structure for creating a new dict type.
type CreateDictTypeReq struct {
	g.Meta `path:"/sys-dict-type" method:"post" summary:"Create a new dict type" tags:"System Dict Type"`
	model.SysDictTypeCreateIn
}

// CreateDictTypeRes defines the response structure for creating a new dict type.
type CreateDictTypeRes struct {
	Id int64 `json:"id"`
}

// GetDictTypeReq defines the request structure for retrieving a dict type.
type GetDictTypeReq struct {
	g.Meta `path:"/sys-dict-type/{id}" method:"get" summary:"Retrieve a dict type by ID" tags:"System Dict Type"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
}

// GetDictTypeRes defines the response structure for retrieving a dict type.
type GetDictTypeRes struct {
	*model.SysDictTypeGetOut
}

// UpdateDictTypeReq defines the request structure for updating a dict type.
type UpdateDictTypeReq struct {
	g.Meta `path:"/sys-dict-type/{id}" method:"put" summary:"Update a dict type by ID" tags:"System Dict Type"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
	model.SysDictTypeUpdateIn
}

// UpdateDictTypeRes defines the response structure for updating a dict type.
type UpdateDictTypeRes struct{}

// DeleteDictTypeReq defines the request structure for deleting a dict type.
type DeleteDictTypeReq struct {
	g.Meta `path:"/sys-dict-type/{id}" method:"delete" summary:"Delete a dict type by ID" tags:"System Dict Type"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
}

// DeleteDictTypeRes defines the response structure for deleting a dict type.
type DeleteDictTypeRes struct{}

// GetDictTypeListReq defines the request structure for listing dict types.
type GetDictTypeListReq struct {
	g.Meta   `path:"/sys-dict-type/list" method:"get" summary:"List dict types" tags:"System Dict Type"`
	Page     int    `json:"page" d:"1"`
	PageSize int    `json:"pageSize" d:"10"`
	TypeCode string `json:"typeCode"`
	TypeName string `json:"typeName"`
	Status   string `json:"status"`
}

// GetDictTypeListRes defines the response structure for listing dict types.
type GetDictTypeListRes struct {
	Items []*entity.SysDictType `json:"items"`
	Total int                   `json:"total"`
}
