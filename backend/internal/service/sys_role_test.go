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

func TestSysRole_CreateRoleWithParent(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleParent%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Create parent role first
		parentIn := model.SysRoleCreateIn{
			Name:        "TestRoleParent",
			Description: "Parent role",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		parentID, err := SysRole().CreateRole(ctx, parentIn)
		t.AssertNil(err)
		t.AssertGT(parentID, uint(0))

		// Create child role with parent
		childIn := model.SysRoleCreateIn{
			Name:        "TestRoleChild",
			Description: "Child role",
			ParentId:    parentID,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		childID, err := SysRole().CreateRole(ctx, childIn)
		t.AssertNil(err)
		t.AssertGT(childID, uint(0))

		// Verify child has correct parent
		var child *entity.SysRole
		err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, childID).Scan(&child)
		t.AssertNil(err)
		t.Assert(child.ParentId, parentID)
	})
}

func TestSysRole_UpdateRoleWithInvalidParent(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleInvalidParent%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Create a role first
		createIn := model.SysRoleCreateIn{
			Name:        "TestRoleInvalidParent",
			Description: "Role for testing invalid parent",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		roleID, err := SysRole().CreateRole(ctx, createIn)
		t.AssertNil(err)

		// Try to update with non-existent parent
		updateIn := model.SysRoleUpdateIn{
			Id:          roleID,
			Name:        "TestRoleInvalidParent",
			Description: "Updated description",
			ParentId:    99999, // Non-existent parent
			Status:      1,
			ModifierId:  1,
			DeptId:      1,
		}
		err = SysRole().UpdateRole(ctx, updateIn)
		t.AssertNE(err, nil)
	})
}

func TestSysRole_DeleteRoleWithChildren(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleDeleteWithChild%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Create parent role
		parentIn := model.SysRoleCreateIn{
			Name:        "TestRoleDeleteWithChild",
			Description: "Parent role",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		parentID, err := SysRole().CreateRole(ctx, parentIn)
		t.AssertNil(err)

		// Create child role
		childIn := model.SysRoleCreateIn{
			Name:        "TestRoleDeleteWithChildChild",
			Description: "Child role",
			ParentId:    parentID,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		_, err = SysRole().CreateRole(ctx, childIn)
		t.AssertNil(err)

		// Try to delete parent role (should fail due to child roles)
		deleteIn := model.SysRoleDeleteIn{Id: parentID}
		err = SysRole().DeleteRole(ctx, deleteIn)
		t.AssertNE(err, nil)
	})
}

func TestSysRole_UpdateRoleWithSelfParent(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestRoleSelfParent%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Create a role first
		createIn := model.SysRoleCreateIn{
			Name:        "TestRoleSelfParent",
			Description: "Role for testing self parent",
			ParentId:    0,
			Status:      1,
			CreatorId:   1,
			ModifierId:  1,
			DeptId:      1,
		}
		roleID, err := SysRole().CreateRole(ctx, createIn)
		t.AssertNil(err)

		// Try to update with itself as parent
		updateIn := model.SysRoleUpdateIn{
			Id:          roleID,
			Name:        "TestRoleSelfParent",
			Description: "Updated description",
			ParentId:    roleID, // Self as parent
			Status:      1,
			ModifierId:  1,
			DeptId:      1,
		}
		err = SysRole().UpdateRole(ctx, updateIn)
		t.AssertNE(err, nil)
	})
}

func TestSysRole_AssignRoleToUser(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().WhereLike(dao.SysUser.Columns().Username, "TestUserRole%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestUserRole%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Create a test user
		userID := "550e8400-e29b-41d4-a716-446655440001"
		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
			dao.SysUser.Columns().Username: "TestUserRole1",
			dao.SysUser.Columns().Password: "password123",
			dao.SysUser.Columns().RealName: "Test User",
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		// Create a test role
		roleID, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        "TestUserRole1",
			Description: "Test role",
			Status:      1,
			CreatorId:   1,
		})
		t.AssertNil(err)

		// Test case 1: Successful assignment
		assignIn := model.AssignRoleToUserIn{
			UserId:    userID,
			RoleId:    roleID,
			CreatedBy: 1,
		}
		assignOut, err := SysRole().AssignRoleToUser(ctx, assignIn)
		t.AssertNil(err)
		t.AssertNE(assignOut, nil)
		t.Assert(assignOut.Success, true)

		// Verify assignment
		var userRole *entity.SysUserRole
		err = dao.SysUserRole.Ctx(ctx).
			Where(dao.SysUserRole.Columns().UserId, userID).
			Where(dao.SysUserRole.Columns().RoleId, roleID).
			Scan(&userRole)
		t.AssertNil(err)
		t.AssertNE(userRole, nil)

		// Test case 2: Duplicate assignment (should succeed but not create duplicate)
		assignOut2, err := SysRole().AssignRoleToUser(ctx, assignIn)
		t.AssertNil(err)
		t.Assert(assignOut2.Success, true)

		// Test case 3: Assign to non-existent user
		assignInInvalidUser := model.AssignRoleToUserIn{
			UserId:    "00000000-0000-0000-0000-000000000000",
			RoleId:    roleID,
			CreatedBy: 1,
		}
		_, err = SysRole().AssignRoleToUser(ctx, assignInInvalidUser)
		t.AssertNE(err, nil)

		// Test case 4: Assign non-existent role
		assignInInvalidRole := model.AssignRoleToUserIn{
			UserId:    userID,
			RoleId:    99999,
			CreatedBy: 1,
		}
		_, err = SysRole().AssignRoleToUser(ctx, assignInInvalidRole)
		t.AssertNE(err, nil)
	})
}

func TestSysRole_RemoveRoleFromUser(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().WhereLike(dao.SysUser.Columns().Username, "TestUserRoleRemove%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestUserRoleRemove%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Create a test user
		userID := "550e8400-e29b-41d4-a716-446655440002"
		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
			dao.SysUser.Columns().Username: "TestUserRoleRemove1",
			dao.SysUser.Columns().Password: "password123",
			
			dao.SysUser.Columns().RealName: "Test User",
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		// Create a test role
		roleID, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        "TestUserRoleRemove1",
			Description: "Test role",
			Status:      1,
			CreatorId:   1,
		})
		t.AssertNil(err)

		// Assign role to user
		_, _ = SysRole().AssignRoleToUser(ctx, model.AssignRoleToUserIn{
			UserId: userID,
			RoleId: roleID,
		})

		// Test case 1: Successful removal
		removeIn := model.RemoveRoleFromUserIn{
			UserId: userID,
			RoleId: roleID,
		}
		removeOut, err := SysRole().RemoveRoleFromUser(ctx, removeIn)
		t.AssertNil(err)
		t.AssertNE(removeOut, nil)
		t.Assert(removeOut.Success, true)

		// Verify removal
		var userRole *entity.SysUserRole
		err = dao.SysUserRole.Ctx(ctx).
			Where(dao.SysUserRole.Columns().UserId, userID).
			Where(dao.SysUserRole.Columns().RoleId, roleID).
			Scan(&userRole)
		t.AssertNil(err)
		t.AssertNil(userRole)

		// Test case 2: Remove already removed assignment (should succeed)
		removeOut2, err := SysRole().RemoveRoleFromUser(ctx, removeIn)
		t.AssertNil(err)
		t.Assert(removeOut2.Success, true)
	})
}

func TestSysRole_GetUserRoles(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().WhereLike(dao.SysUser.Columns().Username, "TestUserRoleGet%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestUserRoleGet%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Create a test user
		userID := "550e8400-e29b-41d4-a716-446655440003"
		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
			dao.SysUser.Columns().Username: "TestUserRoleGet1",
			dao.SysUser.Columns().Password: "password123",
			
			dao.SysUser.Columns().RealName: "Test User",
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		// Create test roles
		roleID1, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        "TestUserRoleGet1",
			Description: "Test role 1",
			Status:      1,
			CreatorId:   1,
		})
		t.AssertNil(err)

		roleID2, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        "TestUserRoleGet2",
			Description: "Test role 2",
			Status:      1,
			CreatorId:   1,
		})
		t.AssertNil(err)

		// Assign roles to user
		_, _ = SysRole().AssignRoleToUser(ctx, model.AssignRoleToUserIn{
			UserId: userID,
			RoleId: roleID1,
		})
		_, _ = SysRole().AssignRoleToUser(ctx, model.AssignRoleToUserIn{
			UserId: userID,
			RoleId: roleID2,
		})

		// Test case 1: Get user roles
		getIn := model.GetUserRolesIn{UserId: userID}
		getOut, err := SysRole().GetUserRoles(ctx, getIn)
		t.AssertNil(err)
		t.AssertNE(getOut, nil)
		t.Assert(len(getOut.Roles), 2)

		// Test case 2: Get roles for user with no roles
		userID2 := "550e8400-e29b-41d4-a716-446655440004"
		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID2,
			dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
			dao.SysUser.Columns().Username: "TestUserRoleGet2",
			dao.SysUser.Columns().Password: "password123",
			
			dao.SysUser.Columns().RealName: "Test User 2",
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		getIn2 := model.GetUserRolesIn{UserId: userID2}
		getOut2, err := SysRole().GetUserRoles(ctx, getIn2)
		t.AssertNil(err)
		t.AssertNE(getOut2, nil)
		t.Assert(len(getOut2.Roles), 0)
	})
}

func TestSysRole_AssignRolesToUser(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().WhereLike(dao.SysUser.Columns().Username, "TestUserRoleAssign%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestUserRoleAssign%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Create a test user
		userID := "550e8400-e29b-41d4-a716-446655440003"
		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID,
			dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
			dao.SysUser.Columns().Username: "TestUserRoleAssign1",
			dao.SysUser.Columns().Password: "password123",

			dao.SysUser.Columns().RealName: "Test User",
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		// Create test roles
		roleID1, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        "TestUserRoleAssign1",
			Description: "Test role 1",
			Status:      1,
			CreatorId:   1,
		})
		t.AssertNil(err)

		roleID2, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        "TestUserRoleAssign2",
			Description: "Test role 2",
			Status:      1,
			CreatorId:   1,
		})
		t.AssertNil(err)

		roleID3, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        "TestUserRoleAssign3",
			Description: "Test role 3",
			Status:      1,
			CreatorId:   1,
		})
		t.AssertNil(err)

		// Test case 1: Assign multiple roles
		assignIn := model.AssignRolesToUserIn{
			UserId:    userID,
			RoleIds:   []uint{roleID1, roleID2, roleID3},
			CreatedBy: 1,
		}
		assignOut, err := SysRole().AssignRolesToUser(ctx, assignIn)
		t.AssertNil(err)
		t.AssertNE(assignOut, nil)
		t.Assert(assignOut.Success, true)
		t.Assert(assignOut.Count, 3)

		// Verify assignments
		getOut, err := SysRole().GetUserRoles(ctx, model.GetUserRolesIn{UserId: userID})
		t.AssertNil(err)
		t.Assert(len(getOut.Roles), 3)

		// Test case 2: Replace roles
		assignIn2 := model.AssignRolesToUserIn{
			UserId:    userID,
			RoleIds:   []uint{roleID1},
			CreatedBy: 1,
		}
		assignOut2, err := SysRole().AssignRolesToUser(ctx, assignIn2)
		t.AssertNil(err)
		t.Assert(assignOut2.Count, 1)

		getOut2, err := SysRole().GetUserRoles(ctx, model.GetUserRolesIn{UserId: userID})
		t.AssertNil(err)
		t.Assert(len(getOut2.Roles), 1)

		// Test case 3: Empty role IDs (should fail)
		assignInInvalid := model.AssignRolesToUserIn{
			UserId:    userID,
			RoleIds:   []uint{},
			CreatedBy: 1,
		}
		_, err = SysRole().AssignRolesToUser(ctx, assignInInvalid)
		t.AssertNE(err, nil)
	})
}

func TestSysRole_GetUsersByRole(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().WhereLike(dao.SysUser.Columns().Username, "TestUsersByRole%").Delete()
		dao.SysRole.Ctx(ctx).Unscoped().WhereLike(dao.SysRole.Columns().Name, "TestUsersByRole%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		var err error
		// Create a test role
		roleID, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        "TestUsersByRole1",
			Description: "Test role",
			Status:      1,
			CreatorId:   1,
		})
		t.AssertNil(err)

		// Create test users
		userID1 := "550e8400-e29b-41d4-a716-446655440004"
		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID1,
			dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
			dao.SysUser.Columns().Username: "TestUsersByRole1",
			dao.SysUser.Columns().Password: "password123",

			dao.SysUser.Columns().RealName: "Test User 1",
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		userID2 := "550e8400-e29b-41d4-a716-446655440005"
		_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
			dao.SysUser.Columns().Id:       userID2,
			dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
			dao.SysUser.Columns().Username: "TestUsersByRole2",
			dao.SysUser.Columns().Password: "password123",

			dao.SysUser.Columns().RealName: "Test User 2",
			dao.SysUser.Columns().Status:   1,
		}).Insert()
		t.AssertNil(err)

		// Assign role to users
		_, _ = SysRole().AssignRoleToUser(ctx, model.AssignRoleToUserIn{
			UserId: userID1,
			RoleId: roleID,
		})
		_, _ = SysRole().AssignRoleToUser(ctx, model.AssignRoleToUserIn{
			UserId: userID2,
			RoleId: roleID,
		})

		// Test case 1: Get users by role
		getIn := model.GetUsersByRoleIn{RoleId: roleID}
		getOut, err := SysRole().GetUsersByRole(ctx, getIn)
		t.AssertNil(err)
		t.AssertNE(getOut, nil)
		t.Assert(len(getOut.Users), 2)

		// Test case 2: Get users for role with no users
		roleID2, err := SysRole().CreateRole(ctx, model.SysRoleCreateIn{
			Name:        "TestUsersByRole2",
			Description: "Test role 2",
			Status:      1,
			CreatorId:   1,
		})
		t.AssertNil(err)

		getIn2 := model.GetUsersByRoleIn{RoleId: roleID2}
		getOut2, err := SysRole().GetUsersByRole(ctx, getIn2)
		t.AssertNil(err)
		t.AssertNE(getOut2, nil)
		t.Assert(len(getOut2.Users), 0)

		// Test case 3: Get users for non-existent role
		getInInvalid := model.GetUsersByRoleIn{RoleId: 99999}
		_, err = SysRole().GetUsersByRole(ctx, getInInvalid)
		t.AssertNE(err, nil)
	})
}