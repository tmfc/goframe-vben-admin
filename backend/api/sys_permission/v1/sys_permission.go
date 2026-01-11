package v1

import (
	"backend/internal/model"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
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

// GetPermissionsByUserReq defines the request structure for retrieving permissions of a user.
type GetPermissionsByUserReq struct {
	g.Meta `path:"/sys-permission/by-user/{userId}" method:"get" summary:"Get permissions by user" tags:"System Permission"`
	UserID int64 `json:"userId" v:"required#用户ID不能为空"`
}

// GetPermissionsByUserRes defines the response structure for retrieving permissions of a user.
type GetPermissionsByUserRes struct {
	*model.UserPermissionsOut
}

// GetPermissionListReq defines the request structure for listing permissions.
type GetPermissionListReq struct {
	g.Meta   `path:"/sys-permission/list" method:"get" summary:"List permissions" tags:"System Permission"`
	Page     int    `json:"page" d:"1"`
	PageSize int    `json:"pageSize" d:"10"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

// GetPermissionListRes defines the response structure for listing permissions.
type GetPermissionListRes struct {
	Items []*entity.SysPermission `json:"items"`
	Total int                     `json:"total"`
}

// GetPermissionTreeReq defines the request structure for retrieving permission tree.
type GetPermissionTreeReq struct {
	g.Meta `path:"/sys-permission/tree" method:"get" summary:"Retrieve permission tree" tags:"System Permission"`
}

// GetPermissionTreeRes defines the response structure for retrieving permission tree.
type GetPermissionTreeRes struct {
	List []*model.SysPermissionTreeOut `json:"list"`
}
