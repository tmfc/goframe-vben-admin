package service

import (
	"context"
	"testing"

	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestSysRole_SyncRoleToCasbin(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	roleName := "TestRoleCasbinSync"
	permName := "TestPermissionCasbinSync"
	scope := "read"

	cleanupCasbin := func() {
		enforcer, err := Casbin(ctx)
		if err != nil {
			return
		}
		_, _ = enforcer.RemoveFilteredNamedPolicy("p", 0, roleName)
	}

	cleanupCasbin()
	dao.SysRolePermission.Ctx(ctx).Where("role_id in (select id from sys_role where name like ?)", "TestRoleCasbin%").Delete()
	dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleCasbin%").Delete()
	dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermissionCasbin%").Delete()

	t.Cleanup(func() {
		cleanupCasbin()
		dao.SysRolePermission.Ctx(ctx).Where("role_id in (select id from sys_role where name like ?)", "TestRoleCasbin%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleCasbin%").Delete()
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermissionCasbin%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		permID, err := SysPermission().CreatePermission(ctx, model.SysPermissionCreateIn{
			Name:        permName,
			Description: "Permission for Casbin sync",
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
		})
		t.AssertNil(err)

		roleID, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        roleName,
			Description: "Role for Casbin sync",
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
		})
		t.AssertNil(err)

		_, err = dao.SysRolePermission.Ctx(ctx).Data(g.Map{
			dao.SysRolePermission.Columns().RoleId:       roleID,
			dao.SysRolePermission.Columns().PermissionId: permID,
			dao.SysRolePermission.Columns().Scope:        scope,
		}).Insert()
		t.AssertNil(err)

		err = SysRole().SyncRoleToCasbin(ctx, roleID)
		t.AssertNil(err)

		var role entity.SysRole
		err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, roleID).Scan(&role)
		t.AssertNil(err)
		domain := NormalizeDomain(gconv.String(role.TenantId))

		count, err := dao.CasbinRule.Ctx(ctx).
			Where(dao.CasbinRule.Columns().Ptype, "p").
			Where(dao.CasbinRule.Columns().V0, roleName).
			Where(dao.CasbinRule.Columns().V1, domain).
			Where(dao.CasbinRule.Columns().V2, permName).
			Where(dao.CasbinRule.Columns().V3, scope).
			Count()
		t.AssertNil(err)
		t.AssertGT(count, 0)
	})
}

func TestSysRole_RemoveRoleFromCasbin(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	roleName := "TestRoleCasbinRemove"
	permName := "TestPermissionCasbinRemove"
	scope := "write"

	cleanupCasbin := func() {
		enforcer, err := Casbin(ctx)
		if err != nil {
			return
		}
		_, _ = enforcer.RemoveFilteredNamedPolicy("p", 0, roleName)
	}

	cleanupCasbin()
	dao.SysRolePermission.Ctx(ctx).Where("role_id in (select id from sys_role where name like ?)", "TestRoleCasbin%").Delete()
	dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleCasbin%").Delete()
	dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermissionCasbin%").Delete()

	t.Cleanup(func() {
		cleanupCasbin()
		dao.SysRolePermission.Ctx(ctx).Where("role_id in (select id from sys_role where name like ?)", "TestRoleCasbin%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleCasbin%").Delete()
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermissionCasbin%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		permID, err := SysPermission().CreatePermission(ctx, model.SysPermissionCreateIn{
			Name:        permName,
			Description: "Permission for Casbin remove",
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
		})
		t.AssertNil(err)

		roleID, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        roleName,
			Description: "Role for Casbin remove",
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
		})
		t.AssertNil(err)

		_, err = dao.SysRolePermission.Ctx(ctx).Data(g.Map{
			dao.SysRolePermission.Columns().RoleId:       roleID,
			dao.SysRolePermission.Columns().PermissionId: permID,
			dao.SysRolePermission.Columns().Scope:        scope,
		}).Insert()
		t.AssertNil(err)

		err = SysRole().SyncRoleToCasbin(ctx, roleID)
		t.AssertNil(err)

		err = SysRole().RemoveRoleFromCasbin(ctx, roleID)
		t.AssertNil(err)

		var role entity.SysRole
		err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, roleID).Scan(&role)
		t.AssertNil(err)
		domain := NormalizeDomain(gconv.String(role.TenantId))

		count, err := dao.CasbinRule.Ctx(ctx).
			Where(dao.CasbinRule.Columns().Ptype, "p").
			Where(dao.CasbinRule.Columns().V0, roleName).
			Where(dao.CasbinRule.Columns().V1, domain).
			Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

func TestSysRole_SyncAllRolesToCasbin(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	roleNameA := "TestRoleCasbinAllA"
	roleNameB := "TestRoleCasbinAllB"
	permNameA := "TestPermissionCasbinAllA"
	permNameB := "TestPermissionCasbinAllB"

	cleanupCasbin := func() {
		enforcer, err := Casbin(ctx)
		if err != nil {
			return
		}
		_, _ = enforcer.RemoveFilteredNamedPolicy("p", 0, roleNameA)
		_, _ = enforcer.RemoveFilteredNamedPolicy("p", 0, roleNameB)
	}

	cleanupCasbin()
	dao.SysRolePermission.Ctx(ctx).Where("role_id in (select id from sys_role where name like ?)", "TestRoleCasbin%").Delete()
	dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleCasbin%").Delete()
	dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermissionCasbin%").Delete()

	t.Cleanup(func() {
		cleanupCasbin()
		dao.SysRolePermission.Ctx(ctx).Where("role_id in (select id from sys_role where name like ?)", "TestRoleCasbin%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleCasbin%").Delete()
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermissionCasbin%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		permA, err := SysPermission().CreatePermission(ctx, model.SysPermissionCreateIn{
			Name:        permNameA,
			Description: "Permission A for Casbin all",
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
		})
		t.AssertNil(err)
		permB, err := SysPermission().CreatePermission(ctx, model.SysPermissionCreateIn{
			Name:        permNameB,
			Description: "Permission B for Casbin all",
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
		})
		t.AssertNil(err)

		roleA, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        roleNameA,
			Description: "Role A for Casbin all",
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
		})
		t.AssertNil(err)
		roleB, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        roleNameB,
			Description: "Role B for Casbin all",
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
		})
		t.AssertNil(err)

		_, err = dao.SysRolePermission.Ctx(ctx).Data(g.Map{
			dao.SysRolePermission.Columns().RoleId:       roleA,
			dao.SysRolePermission.Columns().PermissionId: permA,
		}).Insert()
		t.AssertNil(err)
		_, err = dao.SysRolePermission.Ctx(ctx).Data(g.Map{
			dao.SysRolePermission.Columns().RoleId:       roleB,
			dao.SysRolePermission.Columns().PermissionId: permB,
		}).Insert()
		t.AssertNil(err)

		err = SysRole().SyncAllRolesToCasbin(ctx)
		t.AssertNil(err)

		var role entity.SysRole
		err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, roleA).Scan(&role)
		t.AssertNil(err)
		domain := NormalizeDomain(gconv.String(role.TenantId))

		countA, err := dao.CasbinRule.Ctx(ctx).
			Where(dao.CasbinRule.Columns().Ptype, "p").
			Where(dao.CasbinRule.Columns().V0, roleNameA).
			Where(dao.CasbinRule.Columns().V1, domain).
			Where(dao.CasbinRule.Columns().V2, permNameA).
			Count()
		t.AssertNil(err)
		t.AssertGT(countA, 0)

		countB, err := dao.CasbinRule.Ctx(ctx).
			Where(dao.CasbinRule.Columns().Ptype, "p").
			Where(dao.CasbinRule.Columns().V0, roleNameB).
			Where(dao.CasbinRule.Columns().V1, domain).
			Where(dao.CasbinRule.Columns().V2, permNameB).
			Count()
		t.AssertNil(err)
		t.AssertGT(countB, 0)
	})
}
