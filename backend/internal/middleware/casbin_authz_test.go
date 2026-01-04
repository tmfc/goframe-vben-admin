package middleware

import (
	"context"
	"testing"

	"backend/internal/dao"
	"backend/internal/service"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestCasbinAuthz_AllowsRolePermission(t *testing.T) {
	testutil.RequireDatabase(t)
	ctx := context.TODO()

	roleName := "role_allow"
	path := "/api/allowed"
	method := "get"
	tenantID := "00000000-0000-0000-0000-000000000001"
	userID := "550e8400-e29b-41d4-a716-446655440201"

	t.Cleanup(func() {
		dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Id, tenantID).Delete()
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()
		dao.CasbinRule.Ctx(ctx).Unscoped().Where(dao.CasbinRule.Columns().V0, roleName).Delete()
		resetAuthCache()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctx, tenantID, "Authz Allow Tenant")
		t.AssertNil(err)

		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: tenantID,
			dao.SysUser.Columns().Username: "authz-allow",
			dao.SysUser.Columns().Password: "password",
			dao.SysUser.Columns().Roles:    `["role_allow"]`,
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		err = insertPolicy(ctx, roleName, service.NormalizeDomain(tenantID), path, method)
		t.AssertNil(err)

		enforcer, err := service.Casbin(ctx)
		t.AssertNil(err)
		t.AssertNil(enforcer.LoadPolicy())

		err = authorizeCasbin(ctx, authRequest{
			userID:   userID,
			tenantID: tenantID,
			path:     path,
			method:   method,
		})
		t.AssertNil(err)
	})
}

func TestCasbinAuthz_MultiRoleAllows(t *testing.T) {
	testutil.RequireDatabase(t)
	ctx := context.TODO()

	roleAllow := "role_multi_allow"
	path := "/api/multi"
	method := "post"
	tenantID := "00000000-0000-0000-0000-000000000002"
	userID := "550e8400-e29b-41d4-a716-446655440202"

	t.Cleanup(func() {
		dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Id, tenantID).Delete()
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()
		dao.CasbinRule.Ctx(ctx).Unscoped().Where(dao.CasbinRule.Columns().V0, roleAllow).Delete()
		resetAuthCache()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctx, tenantID, "Authz Multi Tenant")
		t.AssertNil(err)

		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: tenantID,
			dao.SysUser.Columns().Username: "authz-multi",
			dao.SysUser.Columns().Password: "password",
			dao.SysUser.Columns().Roles:    `["role_none","role_multi_allow"]`,
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		err = insertPolicy(ctx, roleAllow, service.NormalizeDomain(tenantID), path, method)
		t.AssertNil(err)

		enforcer, err := service.Casbin(ctx)
		t.AssertNil(err)
		t.AssertNil(enforcer.LoadPolicy())

		err = authorizeCasbin(ctx, authRequest{
			userID:   userID,
			tenantID: tenantID,
			path:     path,
			method:   method,
		})
		t.AssertNil(err)
	})
}

func TestCasbinAuthz_DeniesWithoutPermission(t *testing.T) {
	testutil.RequireDatabase(t)
	ctx := context.TODO()

	userID := "550e8400-e29b-41d4-a716-446655440203"
	tenantID := "00000000-0000-0000-0000-000000000003"

	t.Cleanup(func() {
		dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Id, tenantID).Delete()
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()
		resetAuthCache()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctx, tenantID, "Authz Deny Tenant")
		t.AssertNil(err)

		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: tenantID,
			dao.SysUser.Columns().Username: "authz-deny",
			dao.SysUser.Columns().Password: "password",
			dao.SysUser.Columns().Roles:    `["role_deny"]`,
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		err = authorizeCasbin(ctx, authRequest{
			userID:   userID,
			tenantID: tenantID,
			path:     "/api/denied",
			method:   "get",
		})
		t.AssertNE(err, nil)
	})
}

func TestCasbinAuthz_TenantIsolation(t *testing.T) {
	testutil.RequireDatabase(t)
	ctx := context.TODO()

	roleName := "role_tenant"
	path := "/api/tenant"
	method := "get"
	tenantA := "00000000-0000-0000-0000-000000000004"
	tenantB := "00000000-0000-0000-0000-000000000005"
	userID := "550e8400-e29b-41d4-a716-446655440204"

	t.Cleanup(func() {
		dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Id, tenantA).Delete()
		dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Id, tenantB).Delete()
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()
		dao.CasbinRule.Ctx(ctx).Unscoped().Where(dao.CasbinRule.Columns().V0, roleName).Delete()
		resetAuthCache()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctx, tenantA, "Authz Tenant A")
		t.AssertNil(err)
		err = ensureTenant(ctx, tenantB, "Authz Tenant B")
		t.AssertNil(err)

		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: tenantA,
			dao.SysUser.Columns().Username: "authz-tenant",
			dao.SysUser.Columns().Password: "password",
			dao.SysUser.Columns().Roles:    `["role_tenant"]`,
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		err = insertPolicy(ctx, roleName, service.NormalizeDomain(tenantA), path, method)
		t.AssertNil(err)

		enforcer, err := service.Casbin(ctx)
		t.AssertNil(err)
		t.AssertNil(enforcer.LoadPolicy())

		err = authorizeCasbin(ctx, authRequest{
			userID:   userID,
			tenantID: tenantB,
			path:     path,
			method:   method,
		})
		t.AssertNE(err, nil)

		err = authorizeCasbin(ctx, authRequest{
			userID:   userID,
			tenantID: tenantA,
			path:     path,
			method:   method,
		})
		t.AssertNil(err)
	})
}

func TestCasbinAuthz_UsesCache(t *testing.T) {
	testutil.RequireDatabase(t)
	ctx := context.TODO()

	roleName := "role_cache"
	path := "/api/cache"
	method := "get"
	tenantID := "00000000-0000-0000-0000-000000000006"
	userID := "550e8400-e29b-41d4-a716-446655440205"

	t.Cleanup(func() {
		dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Id, tenantID).Delete()
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()
		dao.CasbinRule.Ctx(ctx).Unscoped().Where(dao.CasbinRule.Columns().V0, roleName).Delete()
		resetAuthCache()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctx, tenantID, "Authz Cache Tenant")
		t.AssertNil(err)

		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: tenantID,
			dao.SysUser.Columns().Username: "authz-cache",
			dao.SysUser.Columns().Password: "password",
			dao.SysUser.Columns().Roles:    `["role_cache"]`,
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		err = insertPolicy(ctx, roleName, service.NormalizeDomain(tenantID), path, method)
		t.AssertNil(err)

		enforcer, err := service.Casbin(ctx)
		t.AssertNil(err)
		t.AssertNil(enforcer.LoadPolicy())

		err = authorizeCasbin(ctx, authRequest{
			userID:   userID,
			tenantID: tenantID,
			path:     path,
			method:   method,
		})
		t.AssertNil(err)

		_, err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Id, userID).Delete()
		t.AssertNil(err)

		err = authorizeCasbin(ctx, authRequest{
			userID:   userID,
			tenantID: tenantID,
			path:     path,
			method:   method,
		})
		t.AssertNil(err)
	})
}

func insertPolicy(ctx context.Context, roleName, domain, path, method string) error {
	_, err := dao.CasbinRule.Ctx(ctx).Data(g.Map{
		dao.CasbinRule.Columns().Ptype: "p",
		dao.CasbinRule.Columns().V0:    roleName,
		dao.CasbinRule.Columns().V1:    domain,
		dao.CasbinRule.Columns().V2:    path,
		dao.CasbinRule.Columns().V3:    method,
	}).Insert()
	return err
}

func ensureTenant(ctx context.Context, tenantID, name string) error {
	count, err := dao.SysTenant.Ctx(ctx).
		Where(dao.SysTenant.Columns().Id, tenantID).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err = dao.SysTenant.Ctx(ctx).Data(g.Map{
		dao.SysTenant.Columns().Id:     tenantID,
		dao.SysTenant.Columns().Name:   name,
		dao.SysTenant.Columns().Status: 1,
	}).Insert()
	return err
}
