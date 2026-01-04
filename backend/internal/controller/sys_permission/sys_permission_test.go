package sys_permission

import (
	"context"
	"testing"

	"backend/api/sys_permission/v1"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/service"

	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func TestSysPermissionController_CreatePermission(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermission%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Successful permission creation
		req := &v1.CreatePermissionReq{
			SysPermissionCreateIn: model.SysPermissionCreateIn{
				Name:        "TestPermission1",
				Description: "Description for TestPermission1",
				ParentId:    0,
				Status:      1,
				CreatorId:   1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		res, err := ctrl.CreatePermission(ctx, req)
		t.AssertNil(err)
		t.AssertNE(res.Id, 0)

		// Test case 2: Permission creation with duplicate name (should fail)
		res, err = ctrl.CreatePermission(ctx, req)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)

		// Test case 3: Permission creation with empty name (should fail due to validation)
		reqInvalid := &v1.CreatePermissionReq{
			SysPermissionCreateIn: model.SysPermissionCreateIn{
				Name:        "",
				Description: "Description for Invalid Permission",
				ParentId:    0,
				Status:      1,
				CreatorId:   1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		res, err = ctrl.CreatePermission(ctx, reqInvalid)
		t.AssertNil(res)
		t.AssertNE(err, nil)
	})
}

func TestSysPermissionController_GetPermission(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermission%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Prepare data: Create a permission first
		createReq := &v1.CreatePermissionReq{
			SysPermissionCreateIn: model.SysPermissionCreateIn{
				Name:        "TestPermissionGet",
				Description: "Description for TestPermissionGet",
				ParentId:    0,
				Status:      1,
				CreatorId:   1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		createRes, _ := ctrl.CreatePermission(ctx, createReq)

		// Test case 1: Successful retrieval
		getReq := &v1.GetPermissionReq{ID: createRes.Id}
		getRes, err := ctrl.GetPermission(ctx, getReq)
		t.AssertNil(err)
		t.AssertNE(getRes.SysPermissionGetOut.SysPermission, nil)
		t.Assert(getRes.SysPermissionGetOut.SysPermission.Name, "TestPermissionGet")

		// Test case 2: Permission not found
		getReqNotFound := &v1.GetPermissionReq{ID: 99999}
		getResNotFound, err := ctrl.GetPermission(ctx, getReqNotFound)
		t.AssertNil(getResNotFound)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 3: Invalid ID (empty) - should fail validation
		getReqInvalid := &v1.GetPermissionReq{ID: 0}
		getResInvalid, err := ctrl.GetPermission(ctx, getReqInvalid)
		t.AssertNil(getResInvalid)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)
	})
}

func TestSysPermissionController_UpdatePermission(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermission%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Prepare data: Create a permission first
		createReq := &v1.CreatePermissionReq{
			SysPermissionCreateIn: model.SysPermissionCreateIn{
				Name:        "TestPermissionUpdate",
				Description: "Description for TestPermissionUpdate",
				ParentId:    0,
				Status:      1,
				CreatorId:   1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		createRes, _ := ctrl.CreatePermission(ctx, createReq)

		// Test case 1: Successful update
		updateReq := &v1.UpdatePermissionReq{
			ID: createRes.Id,
			SysPermissionUpdateIn: model.SysPermissionUpdateIn{
				Name:        "TestPermissionUpdated",
				Description: "Updated Description",
				ParentId:    0,
				Status:      0,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		res, err := ctrl.UpdatePermission(ctx, updateReq)
		t.AssertNil(err)
		t.AssertNE(res, nil)

		// Verify permission was updated
		getRes, _ := ctrl.GetPermission(ctx, &v1.GetPermissionReq{ID: createRes.Id})
		t.Assert(getRes.SysPermissionGetOut.SysPermission.Name, "TestPermissionUpdated")
		t.Assert(getRes.SysPermissionGetOut.SysPermission.Description, "Updated Description")

		// Test case 2: Update non-existent permission
		updateReqNotFound := &v1.UpdatePermissionReq{
			ID: 99999,
			SysPermissionUpdateIn: model.SysPermissionUpdateIn{
				Name:        "NonExistentPermission",
				Description: "Description",
				ParentId:    0,
				Status:      1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		res, err = ctrl.UpdatePermission(ctx, updateReqNotFound)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 3: Update with empty name (should fail due to validation)
		updateReqInvalid := &v1.UpdatePermissionReq{
			ID: createRes.Id,
			SysPermissionUpdateIn: model.SysPermissionUpdateIn{
				Name:        "",
				Description: "Description",
				ParentId:    0,
				Status:      1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		res, err = ctrl.UpdatePermission(ctx, updateReqInvalid)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)
	})
}

func TestSysPermissionController_DeletePermission(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermission%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Prepare data: Create a permission first
		createReq := &v1.CreatePermissionReq{
			SysPermissionCreateIn: model.SysPermissionCreateIn{
				Name:        "TestPermissionDelete",
				Description: "Description for TestPermissionDelete",
				ParentId:    0,
				Status:      1,
				CreatorId:   1,
				ModifierId:  1,
				DeptId:      1,
			},
		}
		createRes, _ := ctrl.CreatePermission(ctx, createReq)

		// Test case 1: Successful deletion
		deleteReq := &v1.DeletePermissionReq{ID: createRes.Id}
		res, err := ctrl.DeletePermission(ctx, deleteReq)
		t.AssertNil(err)
		t.AssertNE(res, nil)

		// Verify permission was deleted
		getRes, err := ctrl.GetPermission(ctx, &v1.GetPermissionReq{ID: createRes.Id})
		t.AssertNil(getRes)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 2: Delete non-existent permission
		deleteReqNotFound := &v1.DeletePermissionReq{ID: 99999}
		res, err = ctrl.DeletePermission(ctx, deleteReqNotFound)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 3: Delete with invalid ID (empty) - should fail validation
		deleteReqInvalid := &v1.DeletePermissionReq{ID: 0}
		res, err = ctrl.DeletePermission(ctx, deleteReqInvalid)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)
	})
}

func TestSysPermissionController_GetPermissionsByUser(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().WhereLike(dao.SysUser.Columns().Username, "testuserperm%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestPermByUser%").Delete()
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermByUser%").Delete()
		dao.SysRolePermission.Ctx(ctx).Unscoped().Where("1=1").Delete()
		dao.CasbinRule.Ctx(ctx).Unscoped().Where("1=1").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		userID, err := service.User().Create(ctx, model.UserCreateIn{
			Username: "testuserperm1",
			Password: "password",
		})
		t.AssertNil(err)

		role1ID, err := service.SysRole().CreateRole(ctx, model.SysRoleCreateIn{Name: "TestPermByUserRole1"})
		t.AssertNil(err)
		role2ID, err := service.SysRole().CreateRole(ctx, model.SysRoleCreateIn{Name: "TestPermByUserRole2"})
		t.AssertNil(err)

		perm1ID, err := service.SysPermission().CreatePermission(ctx, model.SysPermissionCreateIn{Name: "TestPermByUser1"})
		t.AssertNil(err)
		perm2ID, err := service.SysPermission().CreatePermission(ctx, model.SysPermissionCreateIn{Name: "TestPermByUser2"})
		t.AssertNil(err)

		err = service.SysRole().AssignPermissionToRole(ctx, model.SysRolePermissionIn{RoleID: role1ID, PermissionID: perm1ID})
		t.AssertNil(err)
		err = service.SysRole().AssignPermissionToRole(ctx, model.SysRolePermissionIn{RoleID: role2ID, PermissionID: perm2ID})
		t.AssertNil(err)

		enforcer, err := service.Casbin(ctx)
		t.AssertNil(err)
		_, err = enforcer.AddGroupingPolicy("u_"+userID, gconv.String(role1ID), "00000000-0000-0000-0000-000000000000")
		t.AssertNil(err)
		_, err = enforcer.AddGroupingPolicy("u_"+userID, gconv.String(role2ID), "00000000-0000-0000-0000-000000000000")
		t.AssertNil(err)

		permsRes, err := ctrl.GetPermissionsByUser(ctx, &v1.GetPermissionsByUserReq{UserID: userID})
		t.AssertNil(err)
		t.Assert(len(permsRes.PermissionIDs), 2)

		user2ID, err := service.User().Create(ctx, model.UserCreateIn{
			Username: "testuserperm2",
			Password: "password",
		})
		t.AssertNil(err)

		permsRes2, err := ctrl.GetPermissionsByUser(ctx, &v1.GetPermissionsByUserReq{UserID: user2ID})
		t.AssertNil(err)
		t.Assert(len(permsRes2.PermissionIDs), 0)
	})
}
