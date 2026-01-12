package service

import (
	"context"
	"encoding/json"
	"testing"

	v1 "backend/api/auth/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"
)

func ensureTenantByCode(t *testing.T, ctx context.Context, code, name string) int64 {
	t.Helper()
	var tenant entity.SysTenant
	row, err := dao.SysTenant.Ctx(ctx).
		Where(dao.SysTenant.Columns().Code, code).
		One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if row != nil && !row.IsEmpty() {
		if err := row.Struct(&tenant); err != nil {
			t.Fatalf("failed to parse tenant: %v", err)
		}
		return tenant.Id
	}
	result, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
		dao.SysTenant.Columns().Name:   name,
		dao.SysTenant.Columns().Code:   code,
		dao.SysTenant.Columns().Status: 1,
	}).Insert()
	if err != nil {
		t.Fatalf("failed to create tenant: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to read tenant id: %v", err)
	}
	return id
}

func ensureDefaultTenant(t *testing.T, ctx context.Context) {
	t.Helper()
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, 1).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if !existing.IsEmpty() {
		return
	}
	_, err = dao.SysTenant.Ctx(ctx).Data(g.Map{
		dao.SysTenant.Columns().Id:     1,
		dao.SysTenant.Columns().Name:   "System Tenant",
		dao.SysTenant.Columns().Code:   "system",
		dao.SysTenant.Columns().Status: 1,
	}).Insert()
	if err != nil {
		t.Fatalf("failed to create default tenant: %v", err)
	}
}

func createUserForTenant(t *testing.T, ctx context.Context, tenantID int64, username, password string, roles []string) {
	t.Helper()
	rolesJSON, err := json.Marshal(roles)
	if err != nil {
		t.Fatalf("failed to marshal roles: %v", err)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	_, _ = dao.SysUser.CtxNoTenant(ctx).Unscoped().
		Where(dao.SysUser.Columns().TenantId, tenantID).
		Where(dao.SysUser.Columns().Username, username).
		Delete()
	_, err = dao.SysUser.CtxNoTenant(ctx).Data(g.Map{
		dao.SysUser.Columns().TenantId: tenantID,
		dao.SysUser.Columns().Username: username,
		dao.SysUser.Columns().Password: string(hashed),
		dao.SysUser.Columns().Roles:    string(rolesJSON),
		dao.SysUser.Columns().Status:   1,
	}).Insert()
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
}

func TestAuth_LoginMultiTenantEnabled(t *testing.T) {
	testutil.RequireDatabase(t)

	t.Setenv("APP_MULTI_TENANT", "true")
	ctx := context.TODO()

	ensureDefaultTenant(t, ctx)
	tenantCodeA := "mt_alpha"
	tenantCodeB := "mt_beta"
	tenantA := ensureTenantByCode(t, ctx, tenantCodeA, "MultiTenant Alpha")
	tenantB := ensureTenantByCode(t, ctx, tenantCodeB, "MultiTenant Beta")

	t.Cleanup(func() {
		dao.SysUser.CtxNoTenant(ctx).Unscoped().
			Where(dao.SysUser.Columns().Username, "mt_user").
			Delete()
		dao.SysUser.CtxNoTenant(ctx).Unscoped().
			Where(dao.SysUser.Columns().Username, "mt_default").
			Delete()
		dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Code, tenantCodeA).Delete()
		dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Code, tenantCodeB).Delete()
	})

	createUserForTenant(t, ctx, tenantA, "mt_user", "pass123", []string{consts.RoleAdmin})
	createUserForTenant(t, ctx, tenantB, "mt_user", "pass123", []string{consts.RoleAdmin})
	createUserForTenant(t, ctx, 1, "mt_default", "pass123", []string{consts.RoleAdmin})

	gtest.C(t, func(t *gtest.T) {
		resA, err := Auth().Login(ctx, v1.LoginReq{
			Username: "mt_user@" + tenantCodeA,
			Password: "pass123",
		})
		t.AssertNil(err)
		claimsA, _ := ParseAccessToken(resA.AccessToken)
		t.Assert(gconv.String(claimsA["tenantId"]), gconv.String(tenantA))

		resB, err := Auth().Login(ctx, v1.LoginReq{
			Username: "mt_user@" + tenantCodeB,
			Password: "pass123",
		})
		t.AssertNil(err)
		claimsB, _ := ParseAccessToken(resB.AccessToken)
		t.Assert(gconv.String(claimsB["tenantId"]), gconv.String(tenantB))

		resDefault, err := Auth().Login(ctx, v1.LoginReq{
			Username: "mt_default",
			Password: "pass123",
		})
		t.AssertNil(err)
		claimsDefault, _ := ParseAccessToken(resDefault.AccessToken)
		t.Assert(gconv.String(claimsDefault["tenantId"]), consts.DefaultTenantID)
	})
}

func TestAuth_LoginMultiTenantDisabled(t *testing.T) {
	testutil.RequireDatabase(t)

	t.Setenv("APP_MULTI_TENANT", "false")
	ctx := context.TODO()

	ensureDefaultTenant(t, ctx)
	createUserForTenant(t, ctx, 1, "single_user", "pass123", []string{consts.RoleAdmin})

	t.Cleanup(func() {
		dao.SysUser.CtxNoTenant(ctx).Unscoped().
			Where(dao.SysUser.Columns().Username, "single_user").
			Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		res, err := Auth().Login(ctx, v1.LoginReq{
			Username: "single_user@ignored",
			Password: "pass123",
		})
		t.AssertNil(err)
		claims, _ := ParseAccessToken(res.AccessToken)
		t.Assert(gconv.String(claims["tenantId"]), consts.DefaultTenantID)
	})
}

func TestAuth_SwitchTenant(t *testing.T) {
	testutil.RequireDatabase(t)

	t.Setenv("APP_MULTI_TENANT", "true")
	ctx := context.TODO()

	ensureDefaultTenant(t, ctx)
	targetCode := "switch_target"
	targetTenantID := ensureTenantByCode(t, ctx, targetCode, "Switch Target")

	t.Cleanup(func() {
		dao.SysUser.CtxNoTenant(ctx).Unscoped().
			Where(dao.SysUser.Columns().Username, "switch_super").
			Delete()
		dao.SysUser.CtxNoTenant(ctx).Unscoped().
			Where(dao.SysUser.Columns().Username, "switch_admin").
			Delete()
		dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Code, targetCode).Delete()
	})

	createUserForTenant(t, ctx, 1, "switch_super", "pass123", []string{consts.RoleSuper})
	createUserForTenant(t, ctx, 1, "switch_admin", "pass123", []string{consts.RoleAdmin})

	gtest.C(t, func(t *gtest.T) {
		superLogin, err := Auth().Login(ctx, v1.LoginReq{
			Username: "switch_super",
			Password: "pass123",
		})
		t.AssertNil(err)

		switchRes, err := Auth().SwitchTenant(ctx, v1.SwitchTenantReq{
			TenantId: targetTenantID,
			Token:    superLogin.AccessToken,
		})
		t.AssertNil(err)
		claims, _ := ParseAccessToken(switchRes.AccessToken)
		t.Assert(gconv.String(claims["tenantId"]), gconv.String(targetTenantID))

		adminLogin, err := Auth().Login(ctx, v1.LoginReq{
			Username: "switch_admin",
			Password: "pass123",
		})
		t.AssertNil(err)
		_, err = Auth().SwitchTenant(ctx, v1.SwitchTenantReq{
			TenantId: targetTenantID,
			Token:    adminLogin.AccessToken,
		})
		t.AssertNE(err, nil)
	})
}
