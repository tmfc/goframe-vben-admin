package service

import (
	"context"
	"testing"

	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/grand"
)

func TestTenant_CRUD(t *testing.T) {
	testutil.RequireDatabase(t)
	ctx := context.TODO()

	// Reset sequence to handle potential out-of-sync issues
	_, _ = g.DB().Exec(ctx, "SELECT setval('sys_tenant_id_seq', COALESCE((SELECT MAX(id) FROM sys_tenant), 1));")

	gtest.C(t, func(t *gtest.T) {
		// 1. Create
		in := model.TenantCreateIn{
			Name:   "Test CRUD Tenant",
			Code:   "test_crud_tenant_" + grand.S(5),
			Status: 1,
		}
		id, err := Tenant().Create(ctx, in)
		t.AssertNil(err)
		t.AssertGT(id, 0)

		// Cleanup
		defer dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, id).Delete()

		// 2. Get
		item, err := Tenant().Get(ctx, id)
		t.AssertNil(err)
		t.Assert(item.Id, id)
		t.Assert(item.Name, in.Name)
		t.Assert(item.Code, in.Code)

		// 3. Update
		updateIn := model.TenantUpdateIn{
			ID:   id,
			Name: "Updated Tenant Name",
			Code: in.Code,
		}
		err = Tenant().Update(ctx, updateIn)
		t.AssertNil(err)

		item, err = Tenant().Get(ctx, id)
		t.AssertNil(err)
		t.Assert(item.Name, "Updated Tenant Name")

		// 4. List
		listIn := model.TenantListIn{
			Name: "Updated",
		}
		listOut, err := Tenant().List(ctx, listIn)
		t.AssertNil(err)
		t.AssertGT(listOut.Total, 0)

		// 5. Delete
		err = Tenant().Delete(ctx, id)
		t.AssertNil(err)

		item, err = Tenant().Get(ctx, id)
		t.AssertNE(err, nil) // Should return not found error
	})
}
