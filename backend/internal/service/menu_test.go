package service

import (
	"context"
	"testing"

	"backend/api/menu/v1"
	"backend/internal/dao"
	"backend/internal/model"

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
				t.Assert(systemMenuChild.Component, "/system/menu/list")
			}

			// Verify System Dept child exists
			systemDeptChild := findMenuByPath(systemMenu.Children, "/system/dept")
			if systemDeptChild != nil {
				t.Assert(systemDeptChild.Name, "SystemDept")
				t.Assert(systemDeptChild.Component, "/system/dept/list")
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
		in := model.SysMenuCreateIn{
			Name:  "menu-name",
			Path:  "/menu-path",
			Type:  "menu",
			Order: 7,
		}
		data := buildMenuCreateData(in)
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
			ID:    1,
			Name:  "menu-name",
			Path:  "/menu-path",
			Type:  "catalog",
			Order: 11,
		}
		data := buildMenuUpdateData(in)
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
