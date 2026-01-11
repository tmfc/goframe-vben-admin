package tenant

import (
	"context"

	"backend/api/tenant/v1"
	"backend/internal/service"
)

func (c *ControllerV1) UpdateTenant(ctx context.Context, req *v1.UpdateTenantReq) (res *v1.UpdateTenantRes, err error) {
	req.TenantUpdateIn.ID = req.ID
	err = service.Tenant().Update(ctx, req.TenantUpdateIn)
	return
}