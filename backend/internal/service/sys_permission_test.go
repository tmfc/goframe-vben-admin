package service

import (
	"context"
	"testing"

	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSysPermission_CreatePermission(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermission%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Test case 1: Successful permission creation
		createIn := model.SysPermissionCreateIn{
			Name:        "TestPermission1",
			Description: "Description for TestPermission1",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		err = SysPermission().CreatePermission(ctx, createIn)
		t.AssertNil(err)

		// Verify permission was created
		var permission *entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Name, "TestPermission1").Scan(&permission)
		t.AssertNil(err)
		t.AssertNE(permission, nil)
		t.Assert(permission.Name, "TestPermission1")

		// Test case 2: Permission creation with duplicate name (should fail)
		err = SysPermission().CreatePermission(ctx, createIn)
		t.AssertNE(err, nil)

		// Test case 3: Permission creation with empty name (should fail due to validation)
		createInInvalid := model.SysPermissionCreateIn{
			Name:        "",
			Description: "Description for Invalid Permission",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		err = SysPermission().CreatePermission(ctx, createInInvalid)
		t.AssertNE(err, nil)
	})
}

func TestSysPermission_GetPermission(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermission%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			err               error
			getOut            *model.SysPermissionGetOut
			getOutNotFound    *model.SysPermissionGetOut
		)
		// Prepare data: Create a permission first
		createIn := model.SysPermissionCreateIn{
			Name:        "TestPermissionGet",
			Description: "Description for TestPermissionGet",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		_ = SysPermission().CreatePermission(ctx, createIn) // Assuming CreatePermission works

		var permission *entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Name, "TestPermissionGet").Scan(&permission)
		t.AssertNil(err)
		t.AssertNE(permission, nil)

		// Test case 1: Successful retrieval
		getIn := model.SysPermissionGetIn{Id: uint(permission.Id)}
		getOut, err = SysPermission().GetPermission(ctx, getIn)
		t.AssertNil(err)
		t.AssertNE(getOut, nil)
		t.Assert(getOut.Name, "TestPermissionGet")

		// Test case 2: Permission not found
		getInNotFound := model.SysPermissionGetIn{Id: 99999}
		getOutNotFound, err = SysPermission().GetPermission(ctx, getInNotFound)
		t.AssertNE(err, nil)
		t.AssertNil(getOutNotFound)
	})
}

func TestSysPermission_UpdatePermission(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermission%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Prepare data: Create a permission first
		createIn := model.SysPermissionCreateIn{
			Name:        "TestPermissionUpdate",
			Description: "Description for TestPermissionUpdate",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		_ = SysPermission().CreatePermission(ctx, createIn) // Assuming CreatePermission works

		var permission *entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Name, "TestPermissionUpdate").Scan(&permission)
		t.AssertNil(err)
		t.AssertNE(permission, nil)

		// Test case 1: Successful update
		updateIn := model.SysPermissionUpdateIn{
			Id:          uint(permission.Id),
			Name:        "TestPermissionUpdated",
			Description: "Updated Description",
			ParentId:    0,
			Status:      0,
			ModifierId:  1,
			DeptId:      1,
		}
		err = SysPermission().UpdatePermission(ctx, updateIn)
		t.AssertNil(err)

		// Verify permission was updated
		var updatedPermission *entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Id, permission.Id).Scan(&updatedPermission)
		t.AssertNil(err)
		t.AssertNE(updatedPermission, nil)
		t.Assert(updatedPermission.Name, "TestPermissionUpdated")
		t.Assert(updatedPermission.Description, "Updated Description")

		// Test case 2: Update non-existent permission
		updateInNotFound := model.SysPermissionUpdateIn{
			Id:          99999,
			Name:        "NonExistentPermission",
			Description: "Description",
			ParentId:    0,
			Status:      1,
			ModifierId:  1,
			DeptId:      1,
		}
		err = SysPermission().UpdatePermission(ctx, updateInNotFound)
		t.AssertNE(err, nil)

		// Test case 3: Update with empty name (should fail due to validation)
		updateInInvalid := model.SysPermissionUpdateIn{
			Id:          uint(permission.Id),
			Name:        "",
			Description: "Description",
			ParentId:    0,
			Status:      1,
			ModifierId:  1,
			DeptId:      1,
		}
		err = SysPermission().UpdatePermission(ctx, updateInInvalid)
		t.AssertNE(err, nil)
	})
}

func TestSysPermission_DeletePermission(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermission%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Prepare data: Create a permission first
		createIn := model.SysPermissionCreateIn{
			Name:        "TestPermissionDelete",
			Description: "Description for TestPermissionDelete",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		_ = SysPermission().CreatePermission(ctx, createIn) // Assuming CreatePermission works

		var permission *entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Name, "TestPermissionDelete").Scan(&permission)
		t.AssertNil(err)
		t.AssertNE(permission, nil)

		// Test case 1: Successful deletion
		deleteIn := model.SysPermissionDeleteIn{Id: uint(permission.Id)}
		err = SysPermission().DeletePermission(ctx, deleteIn)
		t.AssertNil(err)

		// Verify permission was deleted
		var deletedPermission *entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Id, permission.Id).Scan(&deletedPermission)
		t.AssertNil(err)
		t.AssertNil(deletedPermission) // Should be nil as it's deleted

		// Test case 2: Delete non-existent permission
		deleteInNotFound := model.SysPermissionDeleteIn{Id: 99999}
		err = SysPermission().DeletePermission(ctx, deleteInNotFound)
		t.AssertNE(err, nil)
	})
}
