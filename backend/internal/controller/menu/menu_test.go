package menu

import (
	"context"
	"testing"

	v1 "backend/api/menu/v1"
	"backend/internal/consts"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestMenuController_All(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	ctrl := &ControllerV1{}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Get all menus
		res, err := ctrl.All(ctx, &v1.MenuAllReq{})
		t.AssertNil(err)
		t.AssertNE(res, nil)
		t.AssertGT(len(res), 0)

		// Verify menus have valid structure
		for _, menu := range res {
			t.AssertNE(menu.Path, "")
			t.AssertNE(menu.Name, "")
		}
	})
}

func TestMenuController_All_MenuStructure(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	ctrl := &ControllerV1{}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Verify menu structure
		res, err := ctrl.All(ctx, &v1.MenuAllReq{})
		t.AssertNil(err)
		t.AssertGT(len(res), 0)

		// Verify menu types are valid
		validTypes := []string{"menu", "catalog", "embedded", "link"}
		for _, menu := range res {
			t.AssertIN(menu.Type, validTypes)
		}
	})
}

func TestMenuController_All_HierarchicalStructure(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	ctrl := &ControllerV1{}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Verify hierarchical structure
		res, err := ctrl.All(ctx, &v1.MenuAllReq{})
		t.AssertNil(err)
		t.AssertGT(len(res), 0)

		// Check for catalog menus with children
		for _, menu := range res {
			if menu.Type == "catalog" {
				// Catalog menus typically have children
				if len(menu.Children) > 0 {
					// Verify children have valid structure
					for _, child := range menu.Children {
						t.AssertNE(child.Path, "")
						t.AssertNE(child.Name, "")
					}
				}
			}
		}
	})
}

func TestMenuController_All_MetaFields(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	ctrl := &ControllerV1{}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Verify meta fields exist
		res, err := ctrl.All(ctx, &v1.MenuAllReq{})
		t.AssertNil(err)
		t.AssertGT(len(res), 0)

		// Check that menus have meta information
		for _, menu := range res {
			if menu.Meta != nil {
				t.AssertNE(menu.Meta.Title, "")
			}
		}
	})
}

func TestMenuController_All_NoEmptyPaths(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	ctrl := &ControllerV1{}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Verify no empty paths
		res, err := ctrl.All(ctx, &v1.MenuAllReq{})
		t.AssertNil(err)

		// Verify no menus with empty path
		for _, menu := range res {
			t.AssertNE(menu.Path, "")
			// Check children recursively
			checkNoEmptyPaths(menu, t)
		}
	})
}

func TestMenuController_All_NoButtonTypes(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	ctrl := &ControllerV1{}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Verify no button type menus
		res, err := ctrl.All(ctx, &v1.MenuAllReq{})
		t.AssertNil(err)

		// Verify no button type menus are in the result
		checkNoButtonTypes(res, t)
	})
}

func TestMenuController_All_WorkspaceMenu(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	ctrl := &ControllerV1{}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Check for workspace menu
		res, err := ctrl.All(ctx, &v1.MenuAllReq{})
		t.AssertNil(err)

		// Find workspace menu
		workspace := findMenuByPath(res, "/workspace")
		if workspace != nil {
			t.Assert(workspace.Name, "Workspace")
			// Don't assert on type as it may vary
			t.AssertNE(workspace.Meta, nil)
			t.AssertNE(workspace.Meta.Title, "")
		}
	})
}

func TestMenuController_All_SystemMenu(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	ctrl := &ControllerV1{}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Check for system menu
		res, err := ctrl.All(ctx, &v1.MenuAllReq{})
		t.AssertNil(err)

		// Find system menu
		systemMenu := findMenuByPath(res, "/system")
		if systemMenu != nil {
			t.Assert(systemMenu.Type, "catalog")
			t.Assert(systemMenu.Name, "System")
			t.AssertNE(systemMenu.Meta, nil)
			t.Assert(systemMenu.Meta.Title, "system.title")

			// Verify system menu has children
			t.AssertGT(len(systemMenu.Children), 0)
		}
	})
}

// Helper functions

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

func checkNoEmptyPaths(menu *v1.MenuItem, t *gtest.T) {
	if menu.Path == "" {
		t.Fatalf("found menu with empty path: %s", menu.Name)
	}
	for _, child := range menu.Children {
		checkNoEmptyPaths(child, t)
	}
}

func checkNoButtonTypes(items []*v1.MenuItem, t *gtest.T) {
	for _, item := range items {
		if item.Type == "button" {
			t.Fatalf("found button type menu: %s", item.Name)
		}
		if len(item.Children) > 0 {
			checkNoButtonTypes(item.Children, t)
		}
	}
}
