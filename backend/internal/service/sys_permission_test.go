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
		_, err = SysPermission().CreatePermission(ctx, createIn)
		t.AssertNil(err)

		// Verify permission was created
		var permission *entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Name, "TestPermission1").Scan(&permission)
		t.AssertNil(err)
		t.AssertNE(permission, nil)
		t.Assert(permission.Name, "TestPermission1")

		// Test case 2: Permission creation with duplicate name (should fail)
		_, err = SysPermission().CreatePermission(ctx, createIn)
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
		_, err = SysPermission().CreatePermission(ctx, createInInvalid)
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
			err            error
			getOut         *model.SysPermissionGetOut
			getOutNotFound *model.SysPermissionGetOut
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
		_, _ = SysPermission().CreatePermission(ctx, createIn) // Assuming CreatePermission works

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
		_, _ = SysPermission().CreatePermission(ctx, createIn) // Assuming CreatePermission works

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
		dao.SysMenu.Ctx(ctx).Unscoped().WhereLike(dao.SysMenu.Columns().Name, "TestMenuPermission%").Delete()
		dao.SysRolePermission.Ctx(ctx).Unscoped().Where("1=1").Delete()
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
		_, _ = SysPermission().CreatePermission(ctx, createIn) // Assuming CreatePermission works

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

		// Test case 3: Delete permission with child permissions
		parentIn := model.SysPermissionCreateIn{
			Name:        "TestPermissionParentDelete",
			Description: "Parent permission",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		parentId, err := SysPermission().CreatePermission(ctx, parentIn)
		t.AssertNil(err)
		childIn := model.SysPermissionCreateIn{
			Name:        "TestPermissionChildDelete",
			Description: "Child permission",
			ParentId:    parentId,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		_, err = SysPermission().CreatePermission(ctx, childIn)
		t.AssertNil(err)
		err = SysPermission().DeletePermission(ctx, model.SysPermissionDeleteIn{Id: parentId})
		t.AssertNE(err, nil)

		// Test case 4: Delete permission assigned to roles
		roleBindIn := model.SysPermissionCreateIn{
			Name:        "TestPermissionRoleBind",
			Description: "Role bound permission",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		roleBindId, err := SysPermission().CreatePermission(ctx, roleBindIn)
		t.AssertNil(err)
		_, err = dao.SysRolePermission.Ctx(ctx).Data(g.Map{
			dao.SysRolePermission.Columns().RoleId:       1,
			dao.SysRolePermission.Columns().PermissionId: roleBindId,
		}).Insert()
		t.AssertNil(err)
		err = SysPermission().DeletePermission(ctx, model.SysPermissionDeleteIn{Id: roleBindId})
		t.AssertNE(err, nil)

		// Test case 5: Delete permission managed by menu
		menuId, err := Menu().CreateMenu(ctx, model.SysMenuCreateIn{
			Name:           "TestMenuPermissionA",
			Type:           "menu",
			ParentId:       "0",
			Status:         1,
			Order:          1,
			PermissionCode: "perm:menu:test",
		})
		t.AssertNil(err)
		var menuPerm *entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).
			Where(dao.SysPermission.Columns().Name, "TestMenuPermissionA").
			Scan(&menuPerm)
		t.AssertNil(err)
		t.AssertNE(menuPerm, nil)
		err = SysPermission().DeletePermission(ctx, model.SysPermissionDeleteIn{Id: uint(menuPerm.Id)})
		t.AssertNE(err, nil)
		_ = Menu().DeleteMenu(ctx, menuId)
	})
}

func TestSysPermission_CreatePermissionWithParent(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermissionParent%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Create parent permission first
		parentIn := model.SysPermissionCreateIn{
			Name:        "TestPermissionParent",
			Description: "Parent permission",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		parentID, err := SysPermission().CreatePermission(ctx, parentIn)
		t.AssertNil(err)
		t.AssertGT(parentID, uint(0))

		// Create child permission with parent
		childIn := model.SysPermissionCreateIn{
			Name:        "TestPermissionChild",
			Description: "Child permission",
			ParentId:    parentID,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		childID, err := SysPermission().CreatePermission(ctx, childIn)
		t.AssertNil(err)
		t.AssertGT(childID, uint(0))

		// Verify child has correct parent
		var child *entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Id, childID).Scan(&child)
		t.AssertNil(err)
		t.Assert(child.ParentId, parentID)
	})
}

func TestSysPermission_ListPermissions(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysPermission.Ctx(ctx).Unscoped().WhereLike(dao.SysPermission.Columns().Name, "TestPermissionList%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Create multiple permissions
		for i := 1; i <= 3; i++ {
			createIn := model.SysPermissionCreateIn{
				Name:        "TestPermissionList" + string(rune('0'+i)),
				Description: "Description " + string(rune('0'+i)),
				ParentId:    0,
				Status:      1,
				CreatorId:   1,
				ModifierId:  1,
				DeptId:      1,
			}
			_, err = SysPermission().CreatePermission(ctx, createIn)
			t.AssertNil(err)
		}

		// List permissions
		permissions, err := dao.SysPermission.Ctx(ctx).
			WhereLike(dao.SysPermission.Columns().Name, "TestPermissionList%").
			All()
		t.AssertNil(err)
		t.AssertGE(len(permissions), 3)
	})
}
