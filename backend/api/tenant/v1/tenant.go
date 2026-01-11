package v1

import (
	"backend/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// CreateTenantReq defines the request structure for creating a new tenant.
type CreateTenantReq struct {
	g.Meta `path:"/sys-tenant" method:"post" summary:"Create a new tenant" tags:"System Tenant"`
	model.TenantCreateIn
}

// CreateTenantRes defines the response structure for creating a new tenant.
type CreateTenantRes struct {
	Id int64 `json:"id"`
}

// UpdateTenantReq defines the request structure for updating a tenant.
type UpdateTenantReq struct {
	g.Meta `path:"/sys-tenant/{id}" method:"put" summary:"Update a tenant" tags:"System Tenant"`
	ID     int64 `json:"id" v:"required#ID不能为空" in:"path"`
	model.TenantUpdateIn
}

// UpdateTenantRes defines the response structure for updating a tenant.
type UpdateTenantRes struct{}

// DeleteTenantReq defines the request structure for deleting a tenant.
type DeleteTenantReq struct {
	g.Meta `path:"/sys-tenant/{id}" method:"delete" summary:"Delete a tenant" tags:"System Tenant"`
	ID     int64 `json:"id" v:"required#ID不能为空" in:"path"`
}

// DeleteTenantRes defines the response structure for deleting a tenant.
type DeleteTenantRes struct{}

// GetTenantReq defines the request structure for retrieving a tenant.
type GetTenantReq struct {
	g.Meta `path:"/sys-tenant/{id}" method:"get" summary:"Get tenant" tags:"System Tenant"`
	ID     int64 `json:"id" v:"required#ID不能为空" in:"path"`
}

// GetTenantRes defines the response structure for retrieving a tenant.
type GetTenantRes struct {
	*model.TenantListItem
}

// ListTenantReq defines the request structure for listing tenants.
type ListTenantReq struct {
	g.Meta `path:"/sys-tenant" method:"get" summary:"List tenants" tags:"System Tenant"`
	model.TenantListIn
}

// ListTenantRes defines the response structure for listing tenants.
type ListTenantRes struct {
	*model.TenantListOut
}
