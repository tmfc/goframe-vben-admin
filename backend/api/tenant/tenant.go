// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package tenant

import (
	"context"

	"backend/api/tenant/v1"
)

type ITenantV1 interface {
	CreateTenant(ctx context.Context, req *v1.CreateTenantReq) (res *v1.CreateTenantRes, err error)
	UpdateTenant(ctx context.Context, req *v1.UpdateTenantReq) (res *v1.UpdateTenantRes, err error)
	DeleteTenant(ctx context.Context, req *v1.DeleteTenantReq) (res *v1.DeleteTenantRes, err error)
	GetTenant(ctx context.Context, req *v1.GetTenantReq) (res *v1.GetTenantRes, err error)
	ListTenant(ctx context.Context, req *v1.ListTenantReq) (res *v1.ListTenantRes, err error)
}
