package tenant

import (
	"context"

	"backend/api/tenant/v1"
	"backend/internal/service"
)

func (c *ControllerV1) CreateTenant(ctx context.Context, req *v1.CreateTenantReq) (res *v1.CreateTenantRes, err error) {
	id, err := service.Tenant().Create(ctx, req.TenantCreateIn)
	if err != nil {
		return nil, err
	}
	res = &v1.CreateTenantRes{
		Id: id,
	}
	return
}