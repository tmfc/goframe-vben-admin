package tenant

import (
	"context"

	"backend/api/tenant/v1"
	"backend/internal/service"
)

func (c *ControllerV1) ListTenant(ctx context.Context, req *v1.ListTenantReq) (res *v1.ListTenantRes, err error) {
	out, err := service.Tenant().List(ctx, req.TenantListIn)
	if err != nil {
		return nil, err
	}
	res = &v1.ListTenantRes{
		TenantListOut: out,
	}
	return
}