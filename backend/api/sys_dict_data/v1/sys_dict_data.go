package v1

import (
	"backend/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// CreateDictDataReq defines the request structure for creating a new dict data.
type CreateDictDataReq struct {
	g.Meta `path:"/sys-dict-data" method:"post" summary:"Create a new dict data" tags:"System Dict Data"`
	model.SysDictDataCreateIn
}

// CreateDictDataRes defines the response structure for creating a new dict data.
type CreateDictDataRes struct {
	Id int64 `json:"id"`
}

// GetDictDataReq defines the request structure for retrieving a dict data.
type GetDictDataReq struct {
	g.Meta `path:"/sys-dict-data/{id}" method:"get" summary:"Retrieve a dict data by ID" tags:"System Dict Data"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
}

// GetDictDataRes defines the response structure for retrieving a dict data.
type GetDictDataRes struct {
	*model.SysDictDataGetOut
}

// UpdateDictDataReq defines the request structure for updating a dict data.
type UpdateDictDataReq struct {
	g.Meta `path:"/sys-dict-data/{id}" method:"put" summary:"Update a dict data by ID" tags:"System Dict Data"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
	model.SysDictDataUpdateIn
}

// UpdateDictDataRes defines the response structure for updating a dict data.
type UpdateDictDataRes struct{}

// DeleteDictDataReq defines the request structure for deleting a dict data.
type DeleteDictDataReq struct {
	g.Meta `path:"/sys-dict-data/{id}" method:"delete" summary:"Delete a dict data by ID" tags:"System Dict Data"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
}

// DeleteDictDataRes defines the response structure for deleting a dict data.
type DeleteDictDataRes struct{}

// GetDictDataListReq defines the request structure for listing dict data.
type GetDictDataListReq struct {
	g.Meta     `path:"/sys-dict-data/list" method:"get" summary:"List dict data" tags:"System Dict Data"`
	Page       int    `json:"page" d:"1"`
	PageSize   int    `json:"pageSize" d:"10"`
	DictTypeId int64  `json:"dictTypeId"`
	TypeCode   string `json:"typeCode"`
	Label      string `json:"label"`
	Value      string `json:"value"`
	Status     string `json:"status"`
}

// GetDictDataListRes defines the response structure for listing dict data.
type GetDictDataListRes struct {
	Items []*model.SysDictDataDetail `json:"items"`
	Total int                        `json:"total"`
}
