package v1_test

import (
	"context"
	"testing"

	"backend/api/menu/v1"
	"backend/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestGetMenuReq_Validation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test case 1: ID is missing (should fail)
		req := &v1.GetMenuReq{}
		err := g.Validator().Data(req).Run(context.Background())
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "ID不能为空")

		// Test case 2: ID is provided (should pass)
		req.ID = "123"
		err = g.Validator().Data(req).Run(context.Background())
		t.AssertNil(err)
	})
}

func TestCreateMenuReq_Validation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Name and Type are missing (should fail)
		req := &v1.CreateMenuReq{
			SysMenuCreateIn: model.SysMenuCreateIn{},
		}
		err := g.Validator().Data(req).Run(context.Background())
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "名称不能为空") // Assuming "名称不能为空" is the first error

		// Test case 2: Name is provided, Type is missing (should fail)
		req.SysMenuCreateIn.Name = "Test Menu"
		err = g.Validator().Data(req).Run(context.Background())
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "类型不能为空")

		// Test case 3: Name and Type are provided (should pass)
		req.SysMenuCreateIn.Name = "Test Menu"
		req.SysMenuCreateIn.Type = "menu"
		err = g.Validator().Data(req).Run(context.Background())
		t.AssertNil(err)
	})
}

func TestUpdateMenuReq_Validation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test case 1: ID, Name and Type are missing (should fail)
		req := &v1.UpdateMenuReq{
			SysMenuUpdateIn: model.SysMenuUpdateIn{},
		}
		err := g.Validator().Data(req).Run(context.Background())
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "ID不能为空") // Assuming "ID不能为空" is the first error

		// Test case 2: ID is provided, Name and Type are missing (should fail)
		req.SysMenuUpdateIn.ID = "123"
		err = g.Validator().Data(req).Run(context.Background())
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "名称不能为空")

		// Test case 3: ID and Name are provided, Type is missing (should fail)
		req.SysMenuUpdateIn.ID = "123"
		req.SysMenuUpdateIn.Name = "Updated Menu"
		err = g.Validator().Data(req).Run(context.Background())
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "类型不能为空")

		// Test case 4: All required fields are provided (should pass)
		req.SysMenuUpdateIn.ID = "123"
		req.SysMenuUpdateIn.Name = "Updated Menu"
		req.SysMenuUpdateIn.Type = "menu"
		err = g.Validator().Data(req).Run(context.Background())
		t.AssertNil(err)
	})
}
