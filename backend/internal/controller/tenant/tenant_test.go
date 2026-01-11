package tenant

import (
	"context"
	"testing"

	"backend/api/tenant/v1"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/grand"
)

func TestTenantController_CRUD(t *testing.T) {
	testutil.RequireDatabase(t)
	ctx := context.TODO()
	ctrl := &ControllerV1{}

	// Reset sequence
	_, _ = g.DB().Exec(ctx, "SELECT setval('sys_tenant_id_seq', COALESCE((SELECT MAX(id) FROM sys_tenant), 1));")

	gtest.C(t, func(t *gtest.T) {
		// 1. Create
		in := v1.CreateTenantReq{
			TenantCreateIn: model.TenantCreateIn{
				Name: "API Test Tenant",
				Code: "api_test_" + grand.S(5),
				Status: 1,
			},
		}
		res, err := ctrl.CreateTenant(ctx, &in)
		t.AssertNil(err)
		t.AssertGT(res.Id, 0)

		id := res.Id
		defer dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, id).Delete()

		// 2. Get
		getRes, err := ctrl.GetTenant(ctx, &v1.GetTenantReq{ID: id})
		t.AssertNil(err)
		t.Assert(getRes.Id, id)
		t.Assert(getRes.Name, in.Name)

		// 3. Update
		updateReq := &v1.UpdateTenantReq{
			ID: id,
			TenantUpdateIn: model.TenantUpdateIn{
				Name: "Updated API Tenant",
				Code: in.Code,
			},
		}
		_, err = ctrl.UpdateTenant(ctx, updateReq)
		t.AssertNil(err)

		getRes, err = ctrl.GetTenant(ctx, &v1.GetTenantReq{ID: id})
		t.AssertNil(err)
		t.Assert(getRes.Name, "Updated API Tenant")

		// 4. List
		listRes, err := ctrl.ListTenant(ctx, &v1.ListTenantReq{
			TenantListIn: model.TenantListIn{
				Name: "Updated",
			},
		})
		t.AssertNil(err)
		t.AssertGT(listRes.Total, 0)

		// 5. Delete
		_, err = ctrl.DeleteTenant(ctx, &v1.DeleteTenantReq{ID: id})
		t.AssertNil(err)

		_, err = ctrl.GetTenant(ctx, &v1.GetTenantReq{ID: id})
		t.AssertNE(err, nil)
	})
}
