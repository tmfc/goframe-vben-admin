package sys_role

import (
	"context"
	"testing"

	"backend/api/sys_role/v1"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/service"

	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSysRoleController_CreateRole(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRole%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Successful role creation
		req := &v1.CreateRoleReq{
			SysRoleCreateIn: model.SysRoleCreateIn{
				Name:        "TestRole1",
				Description: "Description for TestRole1",
				ParentId:    0,
				Status:      1,
				CreatorId:   1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		res, err := ctrl.CreateRole(ctx, req)
		t.AssertNil(err)
		t.AssertNE(res.Id, 0)

		// Test case 2: Role creation with duplicate name (should fail)
		res, err = ctrl.CreateRole(ctx, req)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)

		// Test case 3: Role creation with empty name (should fail due to validation)
		reqInvalid := &v1.CreateRoleReq{
			SysRoleCreateIn: model.SysRoleCreateIn{
				Name:        "",
				Description: "Description for Invalid Role",
				ParentId:    0,
				Status:      1,
				CreatorId:   1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		res, err = ctrl.CreateRole(ctx, reqInvalid)
		t.AssertNil(res)
		t.AssertNE(err, nil)
	})
}

func TestSysRoleController_GetRole(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRole%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Prepare data: Create a role first
		createReq := &v1.CreateRoleReq{
			SysRoleCreateIn: model.SysRoleCreateIn{
				Name:        "TestRoleGet",
				Description: "Description for TestRoleGet",
				ParentId:    0,
				Status:      1,
				CreatorId:   1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		createRes, _ := ctrl.CreateRole(ctx, createReq)

		// Test case 1: Successful retrieval
		getReq := &v1.GetRoleReq{ID: createRes.Id}
		getRes, err := ctrl.GetRole(ctx, getReq)
		t.AssertNil(err)
		t.AssertNE(getRes.SysRoleGetOut.SysRole, nil)
		t.Assert(getRes.SysRoleGetOut.SysRole.Name, "TestRoleGet")

		// Test case 2: Role not found
		getReqNotFound := &v1.GetRoleReq{ID: 99999}
		getResNotFound, err := ctrl.GetRole(ctx, getReqNotFound)
		t.AssertNil(getResNotFound)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 3: Invalid ID (empty) - should fail validation
		getReqInvalid := &v1.GetRoleReq{ID: 0}
		getResInvalid, err := ctrl.GetRole(ctx, getReqInvalid)
		t.AssertNil(getResInvalid)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)
	})
}

func TestSysRoleController_UpdateRole(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRole%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Prepare data: Create a role first
		createReq := &v1.CreateRoleReq{
			SysRoleCreateIn: model.SysRoleCreateIn{
				Name:        "TestRoleUpdate",
				Description: "Description for TestRoleUpdate",
				ParentId:    0,
				Status:      1,
				CreatorId:   1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		createRes, _ := ctrl.CreateRole(ctx, createReq)

		// Test case 1: Successful update
		updateReq := &v1.UpdateRoleReq{
			ID: createRes.Id,
			SysRoleUpdateIn: model.SysRoleUpdateIn{
				Name:        "TestRoleUpdated",
				Description: "Updated Description",
				ParentId:    0,
				Status:      0,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		res, err := ctrl.UpdateRole(ctx, updateReq)
		t.AssertNil(err)
		t.AssertNE(res, nil)

		// Verify role was updated
		getRes, _ := ctrl.GetRole(ctx, &v1.GetRoleReq{ID: createRes.Id})
		t.Assert(getRes.SysRoleGetOut.SysRole.Name, "TestRoleUpdated")
		t.Assert(getRes.SysRoleGetOut.SysRole.Description, "Updated Description")

		// Test case 2: Update non-existent role
		updateReqNotFound := &v1.UpdateRoleReq{
			ID: 99999,
			SysRoleUpdateIn: model.SysRoleUpdateIn{
				Name:        "NonExistentRole",
				Description: "Description",
				ParentId:    0,
				Status:      1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		res, err = ctrl.UpdateRole(ctx, updateReqNotFound)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 3: Update with empty name (should fail due to validation)
		updateReqInvalid := &v1.UpdateRoleReq{
			ID: createRes.Id,
			SysRoleUpdateIn: model.SysRoleUpdateIn{
				Name:        "",
				Description: "Description",
				ParentId:    0,
				Status:      1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		res, err = ctrl.UpdateRole(ctx, updateReqInvalid)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)
	})
}

func TestSysRoleController_DeleteRole(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRole%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Prepare data: Create a role first
		createReq := &v1.CreateRoleReq{
			SysRoleCreateIn: model.SysRoleCreateIn{
				Name:        "TestRoleDelete",
				Description: "Description for TestRoleDelete",
				ParentId:    0,
				Status:      1,
				CreatorId:   1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		createRes, _ := ctrl.CreateRole(ctx, createReq)

		// Test case 1: Successful deletion
		deleteReq := &v1.DeleteRoleReq{ID: createRes.Id}
		res, err := ctrl.DeleteRole(ctx, deleteReq)
		t.AssertNil(err)
		t.AssertNE(res, nil)

		// Verify role was deleted
		getRes, err := ctrl.GetRole(ctx, &v1.GetRoleReq{ID: createRes.Id})
		t.AssertNil(getRes)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 2: Delete non-existent role
		deleteReqNotFound := &v1.DeleteRoleReq{ID: 99999}
		res, err = ctrl.DeleteRole(ctx, deleteReqNotFound)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 3: Delete with invalid ID (empty) - should fail validation
		deleteReqInvalid := &v1.DeleteRoleReq{ID: 0}
		res, err = ctrl.DeleteRole(ctx, deleteReqInvalid)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)
	})
}

func TestSysRoleController_AssignAndRemoveUsers(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	userID1 := "550e8400-e29b-41d4-a716-446655440010"
	userID2 := "550e8400-e29b-41d4-a716-446655440011"

	t.Cleanup(func() {
		dao.SysUserRole.Ctx(ctx).Unscoped().WhereIn(dao.SysUserRole.Columns().UserId, []string{userID1, userID2}).Delete()
		dao.SysUser.Ctx(ctx).Unscoped().WhereLike(dao.SysUser.Columns().Username, "TestRoleAssignUsers%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleAssignUsers%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		_, err := dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID1,
			dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
			dao.SysUser.Columns().Username: "TestRoleAssignUsers1",
			dao.SysUser.Columns().Password: "password123",
			dao.SysUser.Columns().RealName: "Test User 1",
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID2,
			dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
			dao.SysUser.Columns().Username: "TestRoleAssignUsers2",
			dao.SysUser.Columns().Password: "password123",
			dao.SysUser.Columns().RealName: "Test User 2",
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		roleRes, err := ctrl.CreateRole(ctx, &v1.CreateRoleReq{
			SysRoleCreateIn: model.SysRoleCreateIn{
				Name:        "TestRoleAssignUsers",
				Description: "Role for assigning users",
				Status:      1,
				CreatorId:   1,
			},
		})
		t.AssertNil(err)

		assignRes, err := ctrl.AssignUsersToRole(ctx, &v1.AssignUsersToRoleReq{
			ID:        roleRes.Id,
			UserIds:   []string{userID1, userID2},
			CreatedBy: 1,
		})
		t.AssertNil(err)
		t.Assert(assignRes.Success, true)

		count, err := dao.SysUserRole.Ctx(ctx).
			Where(dao.SysUserRole.Columns().RoleId, roleRes.Id).
			WhereIn(dao.SysUserRole.Columns().UserId, []string{userID1, userID2}).
			Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		removeRes, err := ctrl.RemoveUsersFromRole(ctx, &v1.RemoveUsersFromRoleReq{
			ID:      roleRes.Id,
			UserIds: []string{userID1, userID2},
		})
		t.AssertNil(err)
		t.Assert(removeRes.Success, true)

		count, err = dao.SysUserRole.Ctx(ctx).
			Where(dao.SysUserRole.Columns().RoleId, roleRes.Id).
			WhereIn(dao.SysUserRole.Columns().UserId, []string{userID1, userID2}).
			Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

func TestSysRoleController_GetRoleUsers(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	userID1 := "550e8400-e29b-41d4-a716-446655440012"
	userID2 := "550e8400-e29b-41d4-a716-446655440013"

	t.Cleanup(func() {
		dao.SysUserRole.Ctx(ctx).Unscoped().WhereIn(dao.SysUserRole.Columns().UserId, []string{userID1, userID2}).Delete()
		dao.SysUser.Ctx(ctx).Unscoped().WhereLike(dao.SysUser.Columns().Username, "TestRoleGetUsers%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleGetUsers%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		_, err := dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID1,
			dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
			dao.SysUser.Columns().Username: "TestRoleGetUsers1",
			dao.SysUser.Columns().Password: "password123",
			dao.SysUser.Columns().RealName: "Test User 1",
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID2,
			dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
			dao.SysUser.Columns().Username: "TestRoleGetUsers2",
			dao.SysUser.Columns().Password: "password123",
			dao.SysUser.Columns().RealName: "Test User 2",
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		roleRes, err := ctrl.CreateRole(ctx, &v1.CreateRoleReq{
			SysRoleCreateIn: model.SysRoleCreateIn{
				Name:        "TestRoleGetUsers",
				Description: "Role for getting users",
				Status:      1,
				CreatorId:   1,
			},
		})
		t.AssertNil(err)

		_, err = ctrl.AssignUsersToRole(ctx, &v1.AssignUsersToRoleReq{
			ID:      roleRes.Id,
			UserIds: []string{userID1, userID2},
		})
		t.AssertNil(err)

		usersRes, err := ctrl.GetRoleUsers(ctx, &v1.GetRoleUsersReq{ID: roleRes.Id})
		t.AssertNil(err)
		t.Assert(len(usersRes.Users), 2)

		userIDs := map[string]bool{}
		for _, user := range usersRes.Users {
			userIDs[user.Id] = true
		}
		t.Assert(userIDs[userID1], true)
		t.Assert(userIDs[userID2], true)
	})
}

func TestSysRoleController_AssignAndRemovePermissions(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	t.Cleanup(func() {
		dao.SysRolePermission.Ctx(ctx).Unscoped().Where("1=1").Delete()
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestRolePerm%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRolePerms%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		roleRes, err := ctrl.CreateRole(ctx, &v1.CreateRoleReq{
			SysRoleCreateIn: model.SysRoleCreateIn{
				Name:        "TestRolePerms",
				Description: "Role for permissions",
				Status:      1,
				CreatorId:   1,
			},
		})
		t.AssertNil(err)

		perm1ID, err := service.SysPermission().CreatePermission(ctx, model.SysPermissionCreateIn{Name: "TestRolePerm1"})
		t.AssertNil(err)
		perm2ID, err := service.SysPermission().CreatePermission(ctx, model.SysPermissionCreateIn{Name: "TestRolePerm2"})
		t.AssertNil(err)

		assignRes, err := ctrl.AssignPermissionsToRole(ctx, &v1.AssignPermissionsToRoleReq{
			ID:            roleRes.Id,
			PermissionIds: []uint{perm1ID, perm2ID},
		})
		t.AssertNil(err)
		t.Assert(assignRes.Success, true)

		count, err := dao.SysRolePermission.Ctx(ctx).
			Where(dao.SysRolePermission.Columns().RoleId, roleRes.Id).
			Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		permsRes, err := ctrl.GetRolePermissions(ctx, &v1.GetRolePermissionsReq{ID: roleRes.Id})
		t.AssertNil(err)
		t.Assert(len(permsRes.PermissionIDs), 2)

		removeRes, err := ctrl.RemovePermissionsFromRole(ctx, &v1.RemovePermissionsFromRoleReq{
			ID:            roleRes.Id,
			PermissionIds: []uint{perm1ID},
		})
		t.AssertNil(err)
		t.Assert(removeRes.Success, true)

		count, err = dao.SysRolePermission.Ctx(ctx).
			Where(dao.SysRolePermission.Columns().RoleId, roleRes.Id).
			Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}
