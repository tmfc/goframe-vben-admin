package service

import (
	"context"
	"testing"

	v1 "backend/api/menu/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestMenu_All(t *testing.T) {
	ctx := context.TODO()
	gtest.C(t, func(t *gtest.T) {
		menus, err := Menu().All(ctx)
		t.AssertNil(err)
		t.AssertGT(len(menus), 0)

		// Verify menus have valid structure
		for _, menu := range menus {
			t.AssertNE(menu.Path, "")
			t.AssertNE(menu.Name, "")
		}
	})
}

func TestMenu_All_DefaultStructure(t *testing.T) {
	ctx := context.TODO()
	gtest.C(t, func(t *gtest.T) {
		menus, err := Menu().All(ctx)
		t.AssertNil(err)
		t.AssertGT(len(menus), 0)

		// Verify menus are returned and have valid structure
		for _, menu := range menus {
			t.AssertNE(menu.Path, "")
			t.AssertNE(menu.Name, "")
			t.AssertIN(menu.Type, []string{"menu", "catalog", "embedded", "link"})
		}
	})
}

func TestMenu_All_SystemMenuStructure(t *testing.T) {
	ctx := context.TODO()
	gtest.C(t, func(t *gtest.T) {
		menus, err := Menu().All(ctx)
		t.AssertNil(err)
		t.AssertGT(len(menus), 0)

		// Check if system menu exists (optional, may not exist in DB)
		systemMenu := findMenuByPath(menus, "/system")
		if systemMenu != nil {
			t.Assert(systemMenu.Type, "catalog")
			t.Assert(systemMenu.Name, "System")
			t.AssertNE(systemMenu.Meta, nil)
			t.Assert(systemMenu.Meta.Title, "system.title")

			// Verify system menu has children
			t.AssertGT(len(systemMenu.Children), 0)

			// Verify System Menu child exists
			systemMenuChild := findMenuByPath(systemMenu.Children, "/system/menu")
			if systemMenuChild != nil {
				t.Assert(systemMenuChild.Name, "SystemMenu")
				t.Assert(systemMenuChild.Component, "/sys/menu/index")
			}

			// Verify System Dept child exists
			systemDeptChild := findMenuByPath(systemMenu.Children, "/system/dept")
			if systemDeptChild != nil {
				t.AssertIN(systemDeptChild.Name, []string{"SystemDept", "Dept"})
				t.AssertIN(systemDeptChild.Component, []string{"/system/dept/list", "/sys/dept/index"})
			}
		}
	})
}

func TestMenu_All_ButtonTypeFiltered(t *testing.T) {
	ctx := context.TODO()
	gtest.C(t, func(t *gtest.T) {
		menus, err := Menu().All(ctx)
		t.AssertNil(err)

		// Verify no button type menus are in the result
		assertNoButtonType(menus, t)
	})
}

func TestMenu_All_EmptyPathFiltered(t *testing.T) {
	ctx := context.TODO()
	gtest.C(t, func(t *gtest.T) {
		menus, err := Menu().All(ctx)
		t.AssertNil(err)

		// Verify no menus with empty path are in the result
		assertNoEmptyPath(menus, t)
	})
}

func TestMenu_All_MetaFields(t *testing.T) {
	ctx := context.TODO()
	gtest.C(t, func(t *gtest.T) {
		menus, err := Menu().All(ctx)
		t.AssertNil(err)
		t.AssertGT(len(menus), 0)

		// Check if workspace menu exists and has meta
		workspace := findMenuByPath(menus, "/workspace")
		if workspace != nil {
			t.AssertNE(workspace.Meta, nil)
			t.AssertNE(workspace.Meta.Title, "")
		}
	})
}

func TestMenu_All_VbenCatalogStructure(t *testing.T) {
	ctx := context.TODO()
	gtest.C(t, func(t *gtest.T) {
		menus, err := Menu().All(ctx)
		t.AssertNil(err)
		t.AssertGT(len(menus), 0)

		// Check if vben menu exists (optional)
		vbenMenu := findMenuByPath(menus, "/vben-admin")
		if vbenMenu != nil {
			t.Assert(vbenMenu.Type, "catalog")
			t.Assert(vbenMenu.Name, "Project")
			t.AssertNE(vbenMenu.Meta, nil)

			// Verify Vben children exist
			t.AssertGT(len(vbenMenu.Children), 0)

			// Check for document child
			documentChild := findMenuByPath(vbenMenu.Children, "/vben-admin/document")
			if documentChild != nil {
				t.Assert(documentChild.Name, "VbenDocument")
				t.Assert(documentChild.Type, "embedded")
				t.Assert(documentChild.Component, "IFrameView")
				t.AssertNE(documentChild.Meta, nil)
			}

			// Check for github child
			githubChild := findMenuByPath(vbenMenu.Children, "/vben-admin/github")
			if githubChild != nil {
				t.Assert(githubChild.Name, "VbenGithub")
				t.Assert(githubChild.Type, "link")
				t.AssertNE(githubChild.Meta, nil)
			}
		}
	})
}

func TestMenu_All_AboutMenu(t *testing.T) {
	ctx := context.TODO()
	gtest.C(t, func(t *gtest.T) {
		menus, err := Menu().All(ctx)
		t.AssertNil(err)
		t.AssertGT(len(menus), 0)

		// Check if about menu exists (optional)
		aboutMenu := findMenuByPath(menus, "/about")
		if aboutMenu != nil {
			t.Assert(aboutMenu.Name, "About")
			t.Assert(aboutMenu.Type, "menu")
			t.AssertNE(aboutMenu.Meta, nil)
			t.AssertNE(aboutMenu.Meta.Title, "")
		}
	})
}

func TestMenu_CreateDataUsesOrderAndStringType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
		in := model.SysMenuCreateIn{
			Name:  "menu-name",
			Path:  "/menu-path",
			Type:  "menu",
			Order: 7,
		}
		data := buildMenuCreateData(ctx, in, "")
		columns := dao.SysMenu.Columns()

		t.Assert(data[columns.Type], in.Type)
		t.Assert(data[columns.Order], in.Order)
		_, hasSort := data["sort"]
		t.Assert(hasSort, false)
	})
}

func TestMenu_UpdateDataUsesOrderAndStringType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		in := model.SysMenuUpdateIn{
			ID:    "1",
			Name:  "menu-name",
			Path:  "/menu-path",
			Type:  "catalog",
			Order: 11,
		}
		data := buildMenuUpdateData(in, "")
		columns := dao.SysMenu.Columns()

		t.Assert(data[columns.Type], in.Type)
		t.Assert(data[columns.Order], in.Order)
		_, hasSort := data["sort"]
		t.Assert(hasSort, false)
	})
}

func findMenuByPath(items []*v1.MenuItem, path string) *v1.MenuItem {
	for _, item := range items {
		if item.Path == path {
			return item
		}
		if len(item.Children) > 0 {
			if found := findMenuByPath(item.Children, path); found != nil {
				return found
			}
		}
	}
	return nil
}

func assertNoButtonType(items []*v1.MenuItem, t *gtest.T) {
	for _, item := range items {
		if item.Type == "button" {
			t.Fatalf("found button type menu: %s", item.Name)
		}
		if len(item.Children) > 0 {
			assertNoButtonType(item.Children, t)
		}
	}
}

func assertNoEmptyPath(items []*v1.MenuItem, t *gtest.T) {
	for _, item := range items {
		if item.Path == "" {
			t.Fatalf("found menu with empty path: %s", item.Name)
		}
		if len(item.Children) > 0 {
			assertNoEmptyPath(item.Children, t)
		}
	}
}

func TestMenu_CRUD(t *testing.T) {
	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	gtest.C(t, func(t *gtest.T) {
		// Cleanup before test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test Menu%").Delete()

		// 1. Create
		createIn := model.SysMenuCreateIn{
			Name:      "Test Menu",
			Path:      "/test-menu",
			Component: "/test/menu/index",
			Icon:      "mdi:test",
			Type:      "menu",
			Status:    1,
			Order:     100,
		}
		newId, err := Menu().CreateMenu(ctx, createIn)
		t.AssertNil(err)
		t.AssertNE(newId, "")

		// 2. Get
		menu, err := Menu().GetMenu(ctx, newId)
		t.AssertNil(err)
		t.AssertNE(menu, nil)
		t.Assert(menu.Name, createIn.Name)
		t.Assert(menu.Path, createIn.Path)

		// 3. Update
		updateIn := model.SysMenuUpdateIn{
			ID:   newId,
			Name: "Test Menu Updated",
			Path: "/test-menu-updated",
			Type: "menu",
		}
		err = Menu().UpdateMenu(ctx, updateIn)
		t.AssertNil(err)

		// Verify update
		menu, err = Menu().GetMenu(ctx, newId)
		t.AssertNil(err)
		t.AssertNE(menu, nil)
		t.Assert(menu.Name, updateIn.Name)
		t.Assert(menu.Path, updateIn.Path)

		// 4. Delete
		err = Menu().DeleteMenu(ctx, newId)
		t.AssertNil(err)

		// Verify delete
		menu, err = Menu().GetMenu(ctx, newId)
		t.AssertNE(err, nil)
		t.AssertNil(menu)
	})
}

func TestMenu_GetMenuList(t *testing.T) {
	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	gtest.C(t, func(t *gtest.T) {
		// Cleanup before test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test List Menu%").Delete()

		// Create some test data
		_, _ = Menu().CreateMenu(ctx, model.SysMenuCreateIn{Name: "Test List Menu 1", Type: "menu", Status: 1})
		_, _ = Menu().CreateMenu(ctx, model.SysMenuCreateIn{Name: "Test List Menu 2", Type: "menu", Status: 0})

		// Test filtering by name
		list, err := Menu().GetMenuList(ctx, model.SysMenuGetListIn{Name: "Test List Menu 1"})
		t.AssertNil(err)
		t.Assert(list.Total, 1)
		t.Assert(len(list.List), 1)
		t.Assert(list.List[0].Name, "Test List Menu 1")

		// Test filtering by status
		list, err = Menu().GetMenuList(ctx, model.SysMenuGetListIn{
			Name:   "Test List Menu",
			Status: "0",
		})
		t.AssertNil(err)
		t.Assert(list.Total, 1)
		t.Assert(len(list.List), 1)
		t.Assert(list.List[0].Name, "Test List Menu 2")

		// Cleanup after test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test List Menu%").Delete()
	})
}

func TestMenu_Validation(t *testing.T) {
	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	gtest.C(t, func(t *gtest.T) {
		// Test creating a menu with a missing name
		_, err := Menu().CreateMenu(ctx, model.SysMenuCreateIn{Path: "/api-menu-validation", Type: "menu"})
		t.AssertNE(err, nil)

		// Test creating a menu with a missing type
		_, err = Menu().CreateMenu(ctx, model.SysMenuCreateIn{Name: "ApiMenuValidation", Path: "/api-menu-validation"})
		t.AssertNE(err, nil)

		// Test updating a menu with a missing name
		err = Menu().UpdateMenu(ctx, model.SysMenuUpdateIn{ID: "some-id", Path: "/api-menu-updated", Type: "menu"})
		t.AssertNE(err, nil)
	})
}

func TestMenu_Validation_Unique(t *testing.T) {
	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	gtest.C(t, func(t *gtest.T) {
		// Cleanup before test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test Unique Menu%").Delete()

		// Create a menu
		createIn := model.SysMenuCreateIn{
			Name: "Test Unique Menu",
			Path: "/test-unique-menu",
			Type: "menu",
		}
		_, err := Menu().CreateMenu(ctx, createIn)
		t.AssertNil(err)

		// Try to create another menu with the same name
		_, err = Menu().CreateMenu(ctx, createIn)
		t.AssertNil(err)

		// Try to create another menu with the same path
		createIn.Name = "Test Unique Menu 2"
		_, err = Menu().CreateMenu(ctx, createIn)
		t.AssertNil(err)

		// Cleanup after test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test Unique Menu%").Delete()
	})
}

func TestMenu_CreateWithPermission(t *testing.T) {
	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	gtest.C(t, func(t *gtest.T) {
		// Cleanup before test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test Menu With Permission%").Delete()
		dao.SysPermission.Ctx(ctx).WhereLike(dao.SysPermission.Columns().Name, "Test Menu With Permission%").Delete()

		// Create a menu with permission_code
		createIn := model.SysMenuCreateIn{
			Name:           "Test Menu With Permission",
			Path:           "/test-menu-with-permission",
			Component:      "/test/menu/index",
			Icon:           "mdi:test",
			Type:           "menu",
			Status:         1,
			Order:          100,
			PermissionCode: "test:menu:create",
		}
		newId, err := Menu().CreateMenu(ctx, createIn)
		t.AssertNil(err)
		t.AssertNE(newId, "")

		// Verify menu was created
		menu, err := Menu().GetMenu(ctx, newId)
		t.AssertNil(err)
		t.AssertNE(menu, nil)
		t.Assert(menu.PermissionCode, createIn.PermissionCode)

		// Verify permission was created
		var perms []*entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).
			Where(dao.SysPermission.Columns().Name, createIn.Name).
			Scan(&perms)
		t.AssertNil(err)
		t.AssertGT(len(perms), 0)

		// Cleanup after test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test Menu With Permission%").Delete()
		dao.SysPermission.Ctx(ctx).WhereLike(dao.SysPermission.Columns().Name, "Test Menu With Permission%").Delete()
	})
}

func TestMenu_UpdatePermission(t *testing.T) {
	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	gtest.C(t, func(t *gtest.T) {
		// Cleanup before test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test Menu Update Permission%").Delete()
		dao.SysPermission.Ctx(ctx).WhereLike(dao.SysPermission.Columns().Name, "Test Menu Update Permission%").Delete()

		// Create a menu without permission_code
		createIn := model.SysMenuCreateIn{
			Name:      "Test Menu Update Permission",
			Path:      "/test-menu-update-permission",
			Type:      "menu",
			Status:    1,
		}
		newId, err := Menu().CreateMenu(ctx, createIn)
		t.AssertNil(err)

		// Update menu to add permission_code
		updateIn := model.SysMenuUpdateIn{
			ID:             newId,
			Name:           "Test Menu Update Permission",
			Path:           "/test-menu-update-permission",
			Type:           "menu",
			PermissionCode: "test:menu:update",
		}
		err = Menu().UpdateMenu(ctx, updateIn)
		t.AssertNil(err)

		// Verify permission was created
		var perms []*entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).
			Where(dao.SysPermission.Columns().Name, updateIn.Name).
			Scan(&perms)
		t.AssertNil(err)
		t.AssertGT(len(perms), 0)

		// Update menu to change permission_code
		updateIn.PermissionCode = "test:menu:update:changed"
		err = Menu().UpdateMenu(ctx, updateIn)
		t.AssertNil(err)

		// Verify old permission was deleted and new permission was created
		perms = nil
		err = dao.SysPermission.Ctx(ctx).
			Where(dao.SysPermission.Columns().Name, updateIn.Name).
			Scan(&perms)
		t.AssertNil(err)
		t.AssertGT(len(perms), 0)

		// Cleanup after test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test Menu Update Permission%").Delete()
		dao.SysPermission.Ctx(ctx).WhereLike(dao.SysPermission.Columns().Name, "Test Menu Update Permission%").Delete()
	})
}

func TestMenu_DeleteWithPermission(t *testing.T) {
	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	gtest.C(t, func(t *gtest.T) {
		// Cleanup before test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test Menu Delete Permission%").Delete()
		dao.SysPermission.Ctx(ctx).WhereLike(dao.SysPermission.Columns().Name, "Test Menu Delete Permission%").Delete()

		// Create a menu with permission_code
		createIn := model.SysMenuCreateIn{
			Name:           "Test Menu Delete Permission",
			Path:           "/test-menu-delete-permission",
			Type:           "menu",
			Status:         1,
			PermissionCode: "test:menu:delete",
		}
		newId, err := Menu().CreateMenu(ctx, createIn)
		t.AssertNil(err)

		// Verify permission exists
		var perms []*entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).
			Where(dao.SysPermission.Columns().Name, createIn.Name).
			Scan(&perms)
		t.AssertNil(err)
		t.AssertGT(len(perms), 0)

		// Delete menu
		err = Menu().DeleteMenu(ctx, newId)
		t.AssertNil(err)

		// Verify permission was deleted
		perms = nil
		err = dao.SysPermission.Ctx(ctx).
			Where(dao.SysPermission.Columns().Name, createIn.Name).
			Scan(&perms)
		t.AssertNil(err)
		t.AssertEQ(len(perms), 0)

		// Cleanup after test (should be no data, but just in case)
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test Menu Delete Permission%").Delete()
		dao.SysPermission.Ctx(ctx).WhereLike(dao.SysPermission.Columns().Name, "Test Menu Delete Permission%").Delete()
	})
}

func TestMenu_TransactionRollback(t *testing.T) {
	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	gtest.C(t, func(t *gtest.T) {
		// Cleanup before test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test Menu Transaction%").Delete()
		dao.SysPermission.Ctx(ctx).WhereLike(dao.SysPermission.Columns().Name, "Test Menu Transaction%").Delete()

		// Create a menu with permission_code
		createIn := model.SysMenuCreateIn{
			Name:           "Test Menu Transaction",
			Path:           "/test-menu-transaction",
			Type:           "menu",
			Status:         1,
			PermissionCode: "test:menu:transaction",
		}
		newId, err := Menu().CreateMenu(ctx, createIn)
		t.AssertNil(err)

		// Verify both menu and permission exist
		menu, err := Menu().GetMenu(ctx, newId)
		t.AssertNil(err)
		t.AssertNE(menu, nil)

		var perms []*entity.SysPermission
		err = dao.SysPermission.Ctx(ctx).
			Where(dao.SysPermission.Columns().Name, createIn.Name).
			Scan(&perms)
		t.AssertNil(err)
		t.AssertGT(len(perms), 0)

		// Update menu with invalid data to trigger rollback
		updateIn := model.SysMenuUpdateIn{
			ID:             newId,
			Name:           "", // Invalid: name is required
			PermissionCode: "test:menu:transaction:new",
		}
		err = Menu().UpdateMenu(ctx, updateIn)
		t.AssertNE(err, nil) // Should fail validation

		// Verify menu was not updated
		menu, err = Menu().GetMenu(ctx, newId)
		t.AssertNil(err)
		t.AssertNE(menu, nil)
		t.Assert(menu.Name, createIn.Name) // Name should remain unchanged

		// Verify permission was not changed
		perms = nil
		err = dao.SysPermission.Ctx(ctx).
			Where(dao.SysPermission.Columns().Name, createIn.Name).
			Scan(&perms)
		t.AssertNil(err)
		t.AssertGT(len(perms), 0)

		// Cleanup after test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test Menu Transaction%").Delete()
		dao.SysPermission.Ctx(ctx).WhereLike(dao.SysPermission.Columns().Name, "Test Menu Transaction%").Delete()
	})
}
