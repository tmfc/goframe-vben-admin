package service

import (
	"context"
	"testing"

	"backend/internal/consts"
	"backend/internal/model"
	"backend/internal/dao"

	"github.com/gogf/gf/v2/test/gtest"

)

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
		list, err = Menu().GetMenuList(ctx, model.SysMenuGetListIn{Status: "0"})
		t.AssertNil(err)
		t.Assert(list.Total, 1)
		t.Assert(len(list.List), 1)
		t.Assert(list.List[0].Name, "Test List Menu 2")

		// Cleanup after test
		dao.SysMenu.Ctx(ctx).WhereLike(dao.SysMenu.Columns().Name, "Test List Menu%").Delete()
	})
}
