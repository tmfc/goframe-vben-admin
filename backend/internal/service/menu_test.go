package service

import (
	"context"
	"testing"

	"backend/api/menu/v1"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestMenu_All(t *testing.T) {
	ctx := context.TODO()
	gtest.C(t, func(t *gtest.T) {
		menus, err := Menu().All(ctx)
		t.AssertNil(err)
		t.AssertGT(len(menus), 0)

		workspace := findMenuByPath(menus, "/workspace")
		t.AssertNE(workspace, nil)
		t.Assert(workspace.Component, "/dashboard/workspace/index")
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
