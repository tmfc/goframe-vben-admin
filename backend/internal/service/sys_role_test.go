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

func TestSysRole_CreateRole(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRole%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Successful role creation
		createIn := model.SysRoleCreateIn{
			Name:        "TestRole1",
			Description: "Description for TestRole1",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		_, err := SysRole().CreateRole(ctx, createIn)
		t.AssertNil(err)

		// Verify role was created
		var role *entity.SysRole
		err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Name, "TestRole1").Scan(&role)
		t.AssertNil(err)
		t.AssertNE(role, nil)
		t.Assert(role.Name, "TestRole1")

		// Test case 2: Role creation with duplicate name (should fail)
		_, err = SysRole().CreateRole(ctx, createIn)
		t.AssertNE(err, nil)

		// Test case 3: Role creation with empty name (should fail due to validation)
		createInInvalid := model.SysRoleCreateIn{
			Name:        "",
			Description: "Description for Invalid Role",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		_, err = SysRole().CreateRole(ctx, createInInvalid)
		t.AssertNE(err, nil)
	})
}

func TestSysRole_GetRole(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRole%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			err            error
			getOut         *model.SysRoleGetOut
			getOutNotFound *model.SysRoleGetOut
		)
		// Prepare data: Create a role first
		createIn := model.SysRoleCreateIn{
			Name:        "TestRoleGet",
			Description: "Description for TestRoleGet",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		_, _ = SysRole().CreateRole(ctx, createIn) // Assuming CreateRole works

		var role *entity.SysRole
		err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Name, "TestRoleGet").Scan(&role)
		t.AssertNil(err)
		t.AssertNE(role, nil)

		// Test case 1: Successful retrieval
		getIn := model.SysRoleGetIn{Id: uint(role.Id)}
		getOut, err = SysRole().GetRole(ctx, getIn)
		t.AssertNil(err)
		t.AssertNE(getOut, nil)
		t.Assert(getOut.Name, "TestRoleGet")

		// Test case 2: Role not found
		getInNotFound := model.SysRoleGetIn{Id: 99999}
		getOutNotFound, err = SysRole().GetRole(ctx, getInNotFound)
		t.AssertNE(err, nil)
		t.AssertNil(getOutNotFound)
	})
}

func TestSysRole_UpdateRole(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRole%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
		)
		// Prepare data: Create a role first
		createIn := model.SysRoleCreateIn{
			Name:        "TestRoleUpdate",
			Description: "Description for TestRoleUpdate",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		_, _ = SysRole().CreateRole(ctx, createIn) // Assuming CreateRole works

		var role *entity.SysRole
		err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Name, "TestRoleUpdate").Scan(&role)
		t.AssertNil(err)
		t.AssertNE(role, nil)

		// Test case 1: Successful update
		updateIn := model.SysRoleUpdateIn{
			Id:          uint(role.Id),
			Name:        "TestRoleUpdated",
			Description: "Updated Description",
			ParentId:    0,
			Status:      0,
			ModifierId:  1,
			DeptId:      1,
		}
		err = SysRole().UpdateRole(ctx, updateIn)
		t.AssertNil(err)

		// Verify role was updated
		var updatedRole *entity.SysRole
		err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, role.Id).Scan(&updatedRole)
		t.AssertNil(err)
		t.AssertNE(updatedRole, nil)
		t.Assert(updatedRole.Name, "TestRoleUpdated")
		t.Assert(updatedRole.Description, "Updated Description")

		// Test case 2: Update non-existent role
		updateInNotFound := model.SysRoleUpdateIn{
			Id:          99999,
			Name:        "NonExistentRole",
			Description: "Description",
			ParentId:    0,
			Status:      1,
			ModifierId:  1,
			DeptId:      1,
		}
		err = SysRole().UpdateRole(ctx, updateInNotFound)
		t.AssertNE(err, nil)

		// Test case 3: Update with empty name (should fail due to validation)
		updateInInvalid := model.SysRoleUpdateIn{
			Id:          uint(role.Id),
			Name:        "",
			Description: "Description",
			ParentId:    0,
			Status:      1,
			ModifierId:  1,
			DeptId:      1,
		}
		err = SysRole().UpdateRole(ctx, updateInInvalid)
		t.AssertNE(err, nil)
	})
}

func TestSysRole_DeleteRole(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRole%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
		)
		// Prepare data: Create a role first
		createIn := model.SysRoleCreateIn{
			Name:        "TestRoleDelete",
			Description: "Description for TestRoleDelete",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		_, _ = SysRole().CreateRole(ctx, createIn) // Assuming CreateRole works

		var role *entity.SysRole
		err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Name, "TestRoleDelete").Scan(&role)
		t.AssertNil(err)
		t.AssertNE(role, nil)

		// Test case 1: Successful deletion
		deleteIn := model.SysRoleDeleteIn{Id: uint(role.Id)}
		err = SysRole().DeleteRole(ctx, deleteIn)
		t.AssertNil(err)

		// Verify role was deleted
		var deletedRole *entity.SysRole
		err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, role.Id).Scan(&deletedRole)
		t.AssertNil(err)
		t.AssertNil(deletedRole) // Should be nil as it's deleted

		// Test case 2: Delete non-existent role
		deleteInNotFound := model.SysRoleDeleteIn{Id: 99999}
		err = SysRole().DeleteRole(ctx, deleteInNotFound)
		t.AssertNE(err, nil)
	})
}
