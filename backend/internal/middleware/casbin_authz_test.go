package middleware

import (
	"context"
	"testing"

	"backend/internal/consts"
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
	tenantID := "101"
	userID := "4201"
	ctxTenant := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)
	ctxNoTenant := dao.WithoutTenant(ctx)

	t.Cleanup(func() {
		dao.SysTenant.Ctx(ctxNoTenant).Unscoped().Where(dao.SysTenant.Columns().Id, tenantID).Delete()
		dao.SysUser.Ctx(ctxNoTenant).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()
		dao.CasbinRule.Ctx(ctxNoTenant).Unscoped().Where(dao.CasbinRule.Columns().V0, roleName).Delete()
		resetAuthCache()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctxNoTenant, tenantID, "Authz Allow Tenant")
		t.AssertNil(err)

		// ensure no residue user with same id
		dao.SysUser.Ctx(ctxNoTenant).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()

		_, err = dao.SysUser.Ctx(ctxTenant).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: tenantID,
			dao.SysUser.Columns().Username: "authz-allow",
			dao.SysUser.Columns().Password: "password",
			dao.SysUser.Columns().Roles:    `["role_allow"]`,
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		err = insertPolicy(ctxTenant, roleName, service.NormalizeDomain(tenantID), path, method)
		t.AssertNil(err)

		enforcer, err := service.Casbin(ctxTenant)
		t.AssertNil(err)
		t.AssertNil(enforcer.LoadPolicy())

		err = authorizeCasbin(ctxTenant, authRequest{
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
	tenantID := "102"
	userID := "4202"
	ctxTenant := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)
	ctxNoTenant := dao.WithoutTenant(ctx)

	t.Cleanup(func() {
		dao.SysTenant.Ctx(ctxNoTenant).Unscoped().Where(dao.SysTenant.Columns().Id, tenantID).Delete()
		dao.SysUser.Ctx(ctxNoTenant).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()
		dao.CasbinRule.Ctx(ctxNoTenant).Unscoped().Where(dao.CasbinRule.Columns().V0, roleAllow).Delete()
		resetAuthCache()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctxNoTenant, tenantID, "Authz Multi Tenant")
		t.AssertNil(err)

		// ensure no residue user with same id
		dao.SysUser.Ctx(ctxNoTenant).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()

		_, err = dao.SysUser.Ctx(ctxTenant).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: tenantID,
			dao.SysUser.Columns().Username: "authz-multi",
			dao.SysUser.Columns().Password: "password",
			dao.SysUser.Columns().Roles:    `["role_none","role_multi_allow"]`,
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		err = insertPolicy(ctxTenant, roleAllow, service.NormalizeDomain(tenantID), path, method)
		t.AssertNil(err)

		enforcer, err := service.Casbin(ctxTenant)
		t.AssertNil(err)
		t.AssertNil(enforcer.LoadPolicy())

		err = authorizeCasbin(ctxTenant, authRequest{
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

	userID := "4203"
	tenantID := "103"
	ctxTenant := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)
	ctxNoTenant := dao.WithoutTenant(ctx)

	t.Cleanup(func() {
		dao.SysTenant.Ctx(ctxNoTenant).Unscoped().Where(dao.SysTenant.Columns().Id, tenantID).Delete()
		dao.SysUser.Ctx(ctxNoTenant).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()
		resetAuthCache()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctxNoTenant, tenantID, "Authz Deny Tenant")
		t.AssertNil(err)

		// ensure no residue user with same id
		dao.SysUser.Ctx(ctxNoTenant).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()

		_, err = dao.SysUser.Ctx(ctxTenant).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: tenantID,
			dao.SysUser.Columns().Username: "authz-deny",
			dao.SysUser.Columns().Password: "password",
			dao.SysUser.Columns().Roles:    `["role_deny"]`,
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		err = authorizeCasbin(ctxTenant, authRequest{
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
	tenantA := "104"
	tenantB := "105"
	userID := "4204"
	ctxTenantA := context.WithValue(ctx, consts.CtxKeyTenantID, tenantA)
	ctxTenantB := context.WithValue(ctx, consts.CtxKeyTenantID, tenantB)
	ctxNoTenant := dao.WithoutTenant(ctx)

	t.Cleanup(func() {
		dao.SysTenant.Ctx(ctxNoTenant).Unscoped().Where(dao.SysTenant.Columns().Id, tenantA).Delete()
		dao.SysTenant.Ctx(ctxNoTenant).Unscoped().Where(dao.SysTenant.Columns().Id, tenantB).Delete()
		dao.SysUser.Ctx(ctxNoTenant).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()
		dao.CasbinRule.Ctx(ctxNoTenant).Unscoped().Where(dao.CasbinRule.Columns().V0, roleName).Delete()
		resetAuthCache()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctxNoTenant, tenantA, "Authz Tenant A")
		t.AssertNil(err)
		err = ensureTenant(ctxNoTenant, tenantB, "Authz Tenant B")
		t.AssertNil(err)

		// ensure no residue user with same id
		dao.SysUser.Ctx(ctxNoTenant).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()

		_, err = dao.SysUser.Ctx(ctxTenantA).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: tenantA,
			dao.SysUser.Columns().Username: "authz-tenant",
			dao.SysUser.Columns().Password: "password",
			dao.SysUser.Columns().Roles:    `["role_tenant"]`,
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		err = insertPolicy(ctxTenantA, roleName, service.NormalizeDomain(tenantA), path, method)
		t.AssertNil(err)

		enforcer, err := service.Casbin(ctxTenantA)
		t.AssertNil(err)
		t.AssertNil(enforcer.LoadPolicy())

		err = authorizeCasbin(ctxTenantB, authRequest{
			userID:   userID,
			tenantID: tenantB,
			path:     path,
			method:   method,
		})
		t.AssertNE(err, nil)

		err = authorizeCasbin(ctxTenantA, authRequest{
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
	tenantID := "106"
	userID := "4205"
	ctxTenant := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)
	ctxNoTenant := dao.WithoutTenant(ctx)

	t.Cleanup(func() {
		dao.SysTenant.Ctx(ctxNoTenant).Unscoped().Where(dao.SysTenant.Columns().Id, tenantID).Delete()
		dao.SysUser.Ctx(ctxNoTenant).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()
		dao.CasbinRule.Ctx(ctxNoTenant).Unscoped().Where(dao.CasbinRule.Columns().V0, roleName).Delete()
		resetAuthCache()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctxNoTenant, tenantID, "Authz Cache Tenant")
		t.AssertNil(err)

		// ensure no residue user with same id
		dao.SysUser.Ctx(ctxNoTenant).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()

		_, err = dao.SysUser.Ctx(ctxTenant).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: tenantID,
			dao.SysUser.Columns().Username: "authz-cache",
			dao.SysUser.Columns().Password: "password",
			dao.SysUser.Columns().Roles:    `["role_cache"]`,
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		err = insertPolicy(ctxTenant, roleName, service.NormalizeDomain(tenantID), path, method)
		t.AssertNil(err)

		enforcer, err := service.Casbin(ctxTenant)
		t.AssertNil(err)
		t.AssertNil(enforcer.LoadPolicy())

		err = authorizeCasbin(ctxTenant, authRequest{
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
