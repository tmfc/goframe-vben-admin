package service

import (
	"context"
	"fmt"
	"testing"

	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestRBACFlow_EndToEnd(t *testing.T) {
	testutil.RequireDatabase(t)
	ctx := context.TODO()

	tenantID := "00000000-0000-0000-0000-000000001000"
	roleName := "TestRoleE2E"
	permPath := "/api/rbac/e2e"
	method := "get"
	userID := "550e8400-e29b-41d4-a716-446655441000"

	t.Cleanup(func() {
		dao.SysUserRole.Ctx(ctx).Unscoped().Where(dao.SysUserRole.Columns().UserId, userID).Delete()
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Id, userID).Delete()
		dao.SysRolePermission.Ctx(ctx).Unscoped().Where("1=1").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().Where(dao.SysRole.Columns().Name, roleName).Delete()
		dao.SysPermission.Ctx(ctx).Unscoped().Where(dao.SysPermission.Columns().Name, permPath).Delete()
		dao.CasbinRule.Ctx(ctx).Unscoped().Where(dao.CasbinRule.Columns().V0, roleName).Delete()
		dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Id, tenantID).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctx, tenantID, "RBAC E2E Tenant")
		t.AssertNil(err)

		permID, err := insertPermission(ctx, tenantID, permPath)
		t.AssertNil(err)

		roleID, err := insertRole(ctx, tenantID, roleName)
		t.AssertNil(err)

		err = assignPermissionWithScope(ctx, roleID, permID, "get")
		t.AssertNil(err)

		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: tenantID,
			dao.SysUser.Columns().Username: "rbac-e2e-user",
			dao.SysUser.Columns().Password: "password",
			dao.SysUser.Columns().Roles:    fmt.Sprintf(`["%s"]`, roleName),
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		err = SysRole().SyncRoleToCasbin(ctx, roleID)
		t.AssertNil(err)

		enforcer, err := Casbin(ctx)
		t.AssertNil(err)
		t.AssertNil(enforcer.LoadPolicy())

		allowed, err := enforcer.Enforce(roleName, NormalizeDomain(tenantID), permPath, method)
		t.AssertNil(err)
		t.Assert(allowed, true)

		denied, err := enforcer.Enforce(roleName, NormalizeDomain(tenantID), permPath, "post")
		t.AssertNil(err)
		t.Assert(denied, false)

		_, err = enforcer.AddGroupingPolicy("u_"+userID, gconv.String(roleID), tenantID)
		t.AssertNil(err)

		perms, err := SysRole().GetPermissionsByUser(ctx, model.UserPermissionsIn{UserID: userID})
		t.AssertNil(err)
		t.Assert(len(perms.PermissionIDs), 1)
		t.Assert(perms.PermissionIDs[0], permID)
	})
}

func TestRBACFlow_TenantIsolation(t *testing.T) {
	testutil.RequireDatabase(t)
	ctx := context.TODO()

	roleName := "TestRoleTenant"
	pathA := "/api/rbac/tenant-a"
	pathB := "/api/rbac/tenant-b"
	tenantA := "00000000-0000-0000-0000-000000001001"
	tenantB := "00000000-0000-0000-0000-000000001002"

	t.Cleanup(func() {
		dao.SysRolePermission.Ctx(ctx).Unscoped().Where("1=1").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().Where(dao.SysRole.Columns().Name, roleName).Delete()
		dao.SysPermission.Ctx(ctx).Unscoped().WhereIn(dao.SysPermission.Columns().Name, []string{pathA, pathB}).Delete()
		dao.CasbinRule.Ctx(ctx).Unscoped().Where(dao.CasbinRule.Columns().V0, roleName).Delete()
		dao.SysTenant.Ctx(ctx).Unscoped().WhereIn(dao.SysTenant.Columns().Id, []string{tenantA, tenantB}).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		err := ensureTenant(ctx, tenantA, "RBAC Tenant A")
		t.AssertNil(err)
		err = ensureTenant(ctx, tenantB, "RBAC Tenant B")
		t.AssertNil(err)

		permA, err := insertPermission(ctx, tenantA, pathA)
		t.AssertNil(err)
		permB, err := insertPermission(ctx, tenantB, pathB)
		t.AssertNil(err)

		roleA, err := insertRole(ctx, tenantA, roleName)
		t.AssertNil(err)
		roleB, err := insertRole(ctx, tenantB, roleName)
		t.AssertNil(err)

		err = assignPermissionWithScope(ctx, roleA, permA, "get")
		t.AssertNil(err)
		err = assignPermissionWithScope(ctx, roleB, permB, "get")
		t.AssertNil(err)

		err = SysRole().SyncRoleToCasbin(ctx, roleA)
		t.AssertNil(err)
		err = SysRole().SyncRoleToCasbin(ctx, roleB)
		t.AssertNil(err)

		enforcer, err := Casbin(ctx)
		t.AssertNil(err)
		t.AssertNil(enforcer.LoadPolicy())

		allowedA, err := enforcer.Enforce(roleName, NormalizeDomain(tenantA), pathA, "get")
		t.AssertNil(err)
		t.Assert(allowedA, true)

		deniedA, err := enforcer.Enforce(roleName, NormalizeDomain(tenantA), pathB, "get")
		t.AssertNil(err)
		t.Assert(deniedA, false)

		allowedB, err := enforcer.Enforce(roleName, NormalizeDomain(tenantB), pathB, "get")
		t.AssertNil(err)
		t.Assert(allowedB, true)

		deniedB, err := enforcer.Enforce(roleName, NormalizeDomain(tenantB), pathA, "get")
		t.AssertNil(err)
		t.Assert(deniedB, false)
	})
}

func BenchmarkSyncRoleToCasbinLargePermissions(b *testing.B) {
	if !databaseConfigured() {
		b.Skip("skip: database.default.link not configured")
	}
	ctx := context.TODO()

	tenantID := "00000000-0000-0000-0000-000000001010"
	roleName := "BenchRoleCasbin"
	permPrefix := "/api/rbac/bench/"

	_ = ensureTenant(ctx, tenantID, "RBAC Bench Tenant")
	roleID, _ := insertRole(ctx, tenantID, roleName)
	for i := 0; i < 25; i++ {
		permID, _ := insertPermission(ctx, tenantID, fmt.Sprintf("%s%d", permPrefix, i))
		_ = assignPermissionWithScope(ctx, roleID, permID, "get")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SysRole().SyncRoleToCasbin(ctx, roleID)
	}
	b.StopTimer()

	dao.SysRolePermission.Ctx(ctx).Unscoped().Where(dao.SysRolePermission.Columns().RoleId, roleID).Delete()
	dao.SysRole.Ctx(ctx).Unscoped().Where(dao.SysRole.Columns().Id, roleID).Delete()
	dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, permPrefix+"%").Delete()
	dao.CasbinRule.Ctx(ctx).Unscoped().Where(dao.CasbinRule.Columns().V0, roleName).Delete()
	dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Id, tenantID).Delete()
}

func databaseConfigured() bool {
	cfg, err := g.Cfg().Get(context.TODO(), "database.default.link")
	if err != nil {
		return false
	}
	if cfg == nil {
		return false
	}
	return cfg.String() != ""
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

func insertRole(ctx context.Context, tenantID, name string) (uint, error) {
	result, err := dao.SysRole.Ctx(ctx).Data(g.Map{
		dao.SysRole.Columns().TenantId: tenantID,
		dao.SysRole.Columns().Name:     name,
		dao.SysRole.Columns().Status:   1,
	}).Insert()
	if err != nil {
		return 0, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(lastID), nil
}

func insertPermission(ctx context.Context, tenantID, name string) (uint, error) {
	result, err := dao.SysPermission.Ctx(ctx).Data(g.Map{
		dao.SysPermission.Columns().TenantId: tenantID,
		dao.SysPermission.Columns().Name:     name,
		dao.SysPermission.Columns().Status:   1,
	}).Insert()
	if err != nil {
		return 0, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(lastID), nil
}

func assignPermissionWithScope(ctx context.Context, roleID, permID uint, scope string) error {
	data := g.Map{
		dao.SysRolePermission.Columns().RoleId:       roleID,
		dao.SysRolePermission.Columns().PermissionId: permID,
	}
	if scope != "" {
		data[dao.SysRolePermission.Columns().Scope] = scope
	}
	_, err := dao.SysRolePermission.Ctx(ctx).Data(data).Insert()
	return err
}
