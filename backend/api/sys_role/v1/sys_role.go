package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"backend/internal/model"
)

// CreateRoleReq defines the request structure for creating a new role.
type CreateRoleReq struct {
	g.Meta `path:"/sys-role" method:"post" summary:"Create a new role" tags:"System Role"`
	model.SysRoleCreateIn
}

// CreateRoleRes defines the response structure for creating a new role.
type CreateRoleRes struct {
	Id uint `json:"id"`
}

// GetRoleReq defines the request structure for retrieving a role.
type GetRoleReq struct {
	g.Meta `path:"/sys-role/{id}" method:"get" summary:"Retrieve a role by ID" tags:"System Role"`
	ID     uint `json:"id" v:"required#ID不能为空"`
}

// GetRoleRes defines the response structure for retrieving a role.
type GetRoleRes struct {
	*model.SysRoleGetOut
}

// UpdateRoleReq defines the request structure for updating a role.
type UpdateRoleReq struct {
	g.Meta `path:"/sys-role/{id}" method:"put" summary:"Update a role by ID" tags:"System Role"`
	ID     uint `json:"id" v:"required#ID不能为空"`
	model.SysRoleUpdateIn
}

// UpdateRoleRes defines the response structure for updating a role.
type UpdateRoleRes struct{}

// DeleteRoleReq defines the request structure for deleting a role.
type DeleteRoleReq struct {
	g.Meta `path:"/sys-role/{id}" method:"delete" summary:"Delete a role by ID" tags:"System Role"`
	ID     uint `json:"id" v:"required#ID不能为空"`
}

// DeleteRoleRes defines the response structure for deleting a role.
type DeleteRoleRes struct{}
