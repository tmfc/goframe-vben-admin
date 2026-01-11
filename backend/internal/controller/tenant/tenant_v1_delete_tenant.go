package tenant

import (
	"context"

	"backend/api/tenant/v1"
	"backend/internal/service"
)

func (c *ControllerV1) DeleteTenant(ctx context.Context, req *v1.DeleteTenantReq) (res *v1.DeleteTenantRes, err error) {
	err = service.Tenant().Delete(ctx, req.ID)
	return
}