package service

import (
	"context"
	"testing"

	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestDept_CreateDept(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	testTenantId := "00000000-0000-0000-0000-000000000000"

	// Ensure tenant exists
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Create root department
		in := model.SysDeptCreateIn{
			Name:      "Engineering",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		id, err := Dept().CreateDept(ctx, in)
		t.AssertNil(err)
		t.AssertNE(id, "")

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, id).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 2: Create child department
		// First create parent
		parentIn := model.SysDeptCreateIn{
			Name:      "Parent Dept",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		parentId, err := Dept().CreateDept(ctx, parentIn)
		t.AssertNil(err)
		t.AssertNE(parentId, "")

		// Create child
		childIn := model.SysDeptCreateIn{
			Name:      "Child Dept",
			ParentId:  parentId,
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		childId, err := Dept().CreateDept(ctx, childIn)
		t.AssertNil(err)
		t.AssertNE(childId, "")

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, childId).Delete()
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, parentId).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 3: Create department with invalid parent
		in := model.SysDeptCreateIn{
			Name:      "Invalid Parent Dept",
			ParentId:  "00000000-0000-0000-0000-000000000999",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		_, err := Dept().CreateDept(ctx, in)
		t.AssertNE(err, nil)
		t.Assert(err.(*gerror.Error).Code(), gcode.CodeValidationFailed)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 4: Create department with empty name
		in := model.SysDeptCreateIn{
			Name:      "",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		_, err := Dept().CreateDept(ctx, in)
		t.AssertNE(err, nil)
	})
}

func TestDept_GetDept(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	testTenantId := "00000000-0000-0000-0000-000000000000"

	// Ensure tenant exists
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Get existing department
		in := model.SysDeptCreateIn{
			Name:      "Test Dept",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		id, err := Dept().CreateDept(ctx, in)
		t.AssertNil(err)

		out, err := Dept().GetDept(ctx, id)
		t.AssertNil(err)
		t.AssertNE(out, nil)
		t.Assert(out.SysDept.Name, "Test Dept")

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, id).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 2: Get non-existent department
		_, err := Dept().GetDept(ctx, "99999")
		t.AssertNE(err, nil)
	})
}

func TestDept_UpdateDept(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	testTenantId := "00000000-0000-0000-0000-000000000000"

	// Ensure tenant exists
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Update department
		in := model.SysDeptCreateIn{
			Name:      "Original Name",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		id, err := Dept().CreateDept(ctx, in)
		t.AssertNil(err)

		updateIn := model.SysDeptUpdateIn{
			ID:         id,
			Name:       "Updated Name",
			ParentId:   "0",
			Status:     1,
			Order:      2,
			ModifierId: 2,
		}
		err = Dept().UpdateDept(ctx, updateIn)
		t.AssertNil(err)

		// Verify update
		out, err := Dept().GetDept(ctx, id)
		t.AssertNil(err)
		t.Assert(out.SysDept.Name, "Updated Name")
		t.Assert(out.SysDept.Order, 2)

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, id).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 2: Update department with invalid parent
		in := model.SysDeptCreateIn{
			Name:      "Test Dept",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		id, err := Dept().CreateDept(ctx, in)
		t.AssertNil(err)

		updateIn := model.SysDeptUpdateIn{
			ID:         id,
			Name:       "Test Dept",
			ParentId:   "99999",
			Status:     1,
			Order:      1,
			ModifierId: 1,
		}
		err = Dept().UpdateDept(ctx, updateIn)
		t.AssertNE(err, nil)

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, id).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 3: Update department to be its own parent
		in := model.SysDeptCreateIn{
			Name:      "Test Dept",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		id, err := Dept().CreateDept(ctx, in)
		t.AssertNil(err)

		updateIn := model.SysDeptUpdateIn{
			ID:         id,
			Name:       "Test Dept",
			ParentId:   id,
			Status:     1,
			Order:      1,
			ModifierId: 1,
		}
		err = Dept().UpdateDept(ctx, updateIn)
		t.AssertNE(err, nil)

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, id).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 4: Update non-existent department
		updateIn := model.SysDeptUpdateIn{
			ID:         "99999",
			Name:       "Non-existent",
			ParentId:   "0",
			Status:     1,
			Order:      1,
			ModifierId: 1,
		}
		err := Dept().UpdateDept(ctx, updateIn)
		t.AssertNE(err, nil)
	})
}

func TestDept_DeleteDept(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	testTenantId := "00000000-0000-0000-0000-000000000000"

	// Ensure tenant exists
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Delete department
		in := model.SysDeptCreateIn{
			Name:      "To Delete",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		id, err := Dept().CreateDept(ctx, in)
		t.AssertNil(err)

		err = Dept().DeleteDept(ctx, id)
		t.AssertNil(err)

		// Verify deletion
		_, err = Dept().GetDept(ctx, id)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 2: Delete department with children
		// Create parent
		parentIn := model.SysDeptCreateIn{
			Name:      "Parent",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		parentId, err := Dept().CreateDept(ctx, parentIn)
		t.AssertNil(err)

		// Create child
		childIn := model.SysDeptCreateIn{
			Name:      "Child",
			ParentId:  parentId,
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		childId, err := Dept().CreateDept(ctx, childIn)
		t.AssertNil(err)

		// Try to delete parent
		err = Dept().DeleteDept(ctx, parentId)
		t.AssertNE(err, nil)

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, childId).Delete()
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, parentId).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 3: Delete non-existent department
		err := Dept().DeleteDept(ctx, "00000000-0000-0000-0000-000000000999")
		t.AssertNE(err, nil)
	})
}

func TestDept_GetDeptList(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	testTenantId := "00000000-0000-0000-0000-000000000000"

	// Ensure tenant exists
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	// Clean up test data
	dao.SysDept.Ctx(ctx).Unscoped().WhereLike(dao.SysDept.Columns().Name, "List Test%").Delete()

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Get all departments
		in := model.SysDeptCreateIn{
			Name:      "List Test 1",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		_, err := Dept().CreateDept(ctx, in)
		t.AssertNil(err)

		in2 := model.SysDeptCreateIn{
			Name:      "List Test 2",
			ParentId:  "0",
			Status:    1,
			Order:     2,
			CreatorId: 1,
		}
		_, err = Dept().CreateDept(ctx, in2)
		t.AssertNil(err)

		listIn := model.SysDeptGetListIn{}
		out, err := Dept().GetDeptList(ctx, listIn)
		t.AssertNil(err)
		t.AssertGE(out.Total, 2)

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().WhereLike(dao.SysDept.Columns().Name, "List Test%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 2: Filter by name
		in := model.SysDeptCreateIn{
			Name:      "Engineering",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		_, err := Dept().CreateDept(ctx, in)
		t.AssertNil(err)

		listIn := model.SysDeptGetListIn{
			Name: "Engineering",
		}
		out, err := Dept().GetDeptList(ctx, listIn)
		t.AssertNil(err)
		t.AssertGE(out.Total, 1)
		t.Assert(out.List[0].SysDept.Name, "Engineering")

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Name, "Engineering").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 3: Filter by status
		in := model.SysDeptCreateIn{
			Name:      "Active Dept",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		_, err := Dept().CreateDept(ctx, in)
		t.AssertNil(err)

		listIn := model.SysDeptGetListIn{
			Status: "1",
		}
		out, err := Dept().GetDeptList(ctx, listIn)
		t.AssertNil(err)
		t.AssertGE(out.Total, 1)

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Name, "Active Dept").Delete()
	})
}

func TestDept_GetDeptTree(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	testTenantId := "00000000-0000-0000-0000-000000000000"

	// Ensure tenant exists
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	// Clean up test data
	dao.SysDept.Ctx(ctx).Unscoped().WhereLike(dao.SysDept.Columns().Name, "Tree Test%").Delete()

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Get department tree with hierarchical structure
		// Create root
		rootIn := model.SysDeptCreateIn{
			Name:      "Tree Test Root",
			ParentId:  "0",
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		rootId, err := Dept().CreateDept(ctx, rootIn)
		t.AssertNil(err)

		// Create child
		childIn := model.SysDeptCreateIn{
			Name:      "Tree Test Child",
			ParentId:  rootId,
			Status:    1,
			Order:     1,
			CreatorId: 1,
		}
		childId, err := Dept().CreateDept(ctx, childIn)
		t.AssertNil(err)

		// Get tree
		tree, err := Dept().GetDeptTree(ctx)
		t.AssertNil(err)
		t.AssertGT(len(tree), 0)

		// Find root in tree
		var root *model.SysDeptTreeOut
		for _, item := range tree {
			if item.Name == "Tree Test Root" {
				root = item
				break
			}
		}
		t.AssertNE(root, nil)
		t.AssertGE(len(root.Children), 1)
		t.Assert(root.Children[0].Name, "Tree Test Child")

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, childId).Delete()
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, rootId).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 2: Inactive departments should not appear in tree
		in := model.SysDeptCreateIn{
			Name:      "Inactive Dept",
			ParentId:  "0",
			Status:    0,
			Order:     1,
			CreatorId: 1,
		}
		id, err := Dept().CreateDept(ctx, in)
		t.AssertNil(err)

		tree, err := Dept().GetDeptTree(ctx)
		t.AssertNil(err)

		// Verify inactive dept is not in tree
		for _, item := range tree {
			t.AssertNE(item.Name, "Inactive Dept")
		}

		// Clean up
		dao.SysDept.Ctx(ctx).Unscoped().Where(dao.SysDept.Columns().Id, id).Delete()
	})
}
