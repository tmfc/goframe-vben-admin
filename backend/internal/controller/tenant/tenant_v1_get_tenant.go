package tenant

import (
	"context"

	"backend/api/tenant/v1"
	"backend/internal/service"
)

func (c *ControllerV1) GetTenant(ctx context.Context, req *v1.GetTenantReq) (res *v1.GetTenantRes, err error) {
	out, err := service.Tenant().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	res = &v1.GetTenantRes{
		TenantListItem: out,
	}
	return
}