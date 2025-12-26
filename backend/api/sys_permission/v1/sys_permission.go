package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"backend/internal/model"
)

// CreatePermissionReq defines the request structure for creating a new permission.
type CreatePermissionReq struct {
	g.Meta `path:"/sys-permission" method:"post" summary:"Create a new permission" tags:"System Permission"`
	model.SysPermissionCreateIn
}

// CreatePermissionRes defines the response structure for creating a new permission.
type CreatePermissionRes struct {
	Id uint `json:"id"`
}

// GetPermissionReq defines the request structure for retrieving a permission.
type GetPermissionReq struct {
	g.Meta `path:"/sys-permission/{id}" method:"get" summary:"Retrieve a permission by ID" tags:"System Permission"`
	ID     uint `json:"id" v:"required#ID不能为空"`
}

// GetPermissionRes defines the response structure for retrieving a permission.
type GetPermissionRes struct {
	*model.SysPermissionGetOut
}

// UpdatePermissionReq defines the request structure for updating a permission.
type UpdatePermissionReq struct {
	g.Meta `path:"/sys-permission/{id}" method:"put" summary:"Update a permission by ID" tags:"System Permission"`
	ID     uint `json:"id" v:"required#ID不能为空"`
	model.SysPermissionUpdateIn
}

// UpdatePermissionRes defines the response structure for updating a permission.
type UpdatePermissionRes struct{}

// DeletePermissionReq defines the request structure for deleting a permission.
type DeletePermissionReq struct {
	g.Meta `path:"/sys-permission/{id}" method:"delete" summary:"Delete a permission by ID" tags:"System Permission"`
	ID     uint `json:"id" v:"required#ID不能为空"`
}

// DeletePermissionRes defines the response structure for deleting a permission.
type DeletePermissionRes struct{}
