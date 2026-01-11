package dept

import (
	"context"
	"testing"

	v1 "backend/api/dept/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestDeptController_CreateDept(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)

	// Cleanup function
	t.Cleanup(func() {
		dao.SysDept.Ctx(ctx).Unscoped().WhereLike(dao.SysDept.Columns().Name, "TestDept%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: Successful root department creation
		req := &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDept1",
				ParentId:  0,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		}
		res, err := ctrl.CreateDept(ctx, req)
		t.AssertNil(err)
		t.AssertNE(res.Id, "")

		// Test case 2: Successful child department creation
		childReq := &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptChild",
				ParentId:  res.Id,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		}
		childRes, err := ctrl.CreateDept(ctx, childReq)
		t.AssertNil(err)
		t.AssertNE(childRes.Id, "")

		// Test case 3: Department creation with empty name (should fail due to validation)
		reqInvalid := &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "",
				ParentId:  0,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		}
		res, err = ctrl.CreateDept(ctx, reqInvalid)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)

		// Test case 4: Department creation with invalid parent
		reqInvalidParent := &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptInvalidParent",
				ParentId:  999999,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		}
		res, err = ctrl.CreateDept(ctx, reqInvalidParent)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)
	})
}

func TestDeptController_GetDept(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysDept.Ctx(ctx).Unscoped().WhereLike(dao.SysDept.Columns().Name, "TestDept%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Prepare data: Create a department first
		createReq := &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptGet",
				ParentId:  0,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		}
		createRes, _ := ctrl.CreateDept(ctx, createReq)

		// Test case 1: Successful retrieval
		getReq := &v1.GetDeptReq{ID: createRes.Id}
		getRes, err := ctrl.GetDept(ctx, getReq)
		t.AssertNil(err)
		t.AssertNE(getRes.SysDeptGetOut.SysDept, nil)
		t.Assert(getRes.SysDeptGetOut.SysDept.Name, "TestDeptGet")

		// Test case 2: Department not found
		getReqNotFound := &v1.GetDeptReq{ID: 999999}
		getResNotFound, err := ctrl.GetDept(ctx, getReqNotFound)
		t.AssertNil(getResNotFound)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 3: Invalid ID (empty) - should fail validation
		getReqInvalid := &v1.GetDeptReq{ID: 0}
		getResInvalid, err := ctrl.GetDept(ctx, getReqInvalid)
		t.AssertNil(getResInvalid)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)
	})
}

func TestDeptController_UpdateDept(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysDept.Ctx(ctx).Unscoped().WhereLike(dao.SysDept.Columns().Name, "TestDept%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Prepare data: Create a department first
		createReq := &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptUpdate",
				ParentId:  0,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		}
		createRes, _ := ctrl.CreateDept(ctx, createReq)

		// Test case 1: Successful update
		updateReq := &v1.UpdateDeptReq{
			ID: createRes.Id,
			SysDeptUpdateIn: model.SysDeptUpdateIn{
				Name:       "TestDeptUpdated",
				ParentId:   0,
				Status:     0,
				Order:      2,
				ModifierId: 1,
			},
		}
		res, err := ctrl.UpdateDept(ctx, updateReq)
		t.AssertNil(err)
		t.AssertNE(res, nil)

		// Verify department was updated
		getRes, _ := ctrl.GetDept(ctx, &v1.GetDeptReq{ID: createRes.Id})
		t.Assert(getRes.SysDeptGetOut.SysDept.Name, "TestDeptUpdated")
		t.Assert(getRes.SysDeptGetOut.SysDept.Order, 2)

		// Test case 2: Update non-existent department
		updateReqNotFound := &v1.UpdateDeptReq{
			ID: 999999,
			SysDeptUpdateIn: model.SysDeptUpdateIn{
				Name:       "NonExistentDept",
				ParentId:   0,
				Status:     1,
				Order:      1,
				ModifierId: 1,
			},
		}
		res, err = ctrl.UpdateDept(ctx, updateReqNotFound)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 3: Update with empty name (should fail due to validation)
		updateReqInvalid := &v1.UpdateDeptReq{
			ID: createRes.Id,
			SysDeptUpdateIn: model.SysDeptUpdateIn{
				Name:       "",
				ParentId:   0,
				Status:     1,
				Order:      1,
				ModifierId: 1,
			},
		}
		res, err = ctrl.UpdateDept(ctx, updateReqInvalid)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)

		// Test case 4: Update department to be its own parent
		updateReqSelfParent := &v1.UpdateDeptReq{
			ID: createRes.Id,
			SysDeptUpdateIn: model.SysDeptUpdateIn{
				Name:       "TestDeptSelfParent",
				ParentId:   createRes.Id,
				Status:     1,
				Order:      1,
				ModifierId: 1,
			},
		}
		res, err = ctrl.UpdateDept(ctx, updateReqSelfParent)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)
	})
}

func TestDeptController_DeleteDept(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysDept.Ctx(ctx).Unscoped().WhereLike(dao.SysDept.Columns().Name, "TestDept%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Prepare data: Create a department first
		createReq := &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptDelete",
				ParentId:  0,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		}
		createRes, _ := ctrl.CreateDept(ctx, createReq)

		// Test case 1: Successful deletion
		deleteReq := &v1.DeleteDeptReq{ID: createRes.Id}
		res, err := ctrl.DeleteDept(ctx, deleteReq)
		t.AssertNil(err)
		t.AssertNE(res, nil)

		// Verify department was deleted
		getRes, err := ctrl.GetDept(ctx, &v1.GetDeptReq{ID: createRes.Id})
		t.AssertNil(getRes)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 2: Delete non-existent department
		deleteReqNotFound := &v1.DeleteDeptReq{ID: 999999}
		res, err = ctrl.DeleteDept(ctx, deleteReqNotFound)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeNotFound)

		// Test case 3: Delete with invalid ID (empty) - should fail validation
		deleteReqInvalid := &v1.DeleteDeptReq{ID: 0}
		res, err = ctrl.DeleteDept(ctx, deleteReqInvalid)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 4: Delete department with children (should fail)
		parentReq := &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptParent",
				ParentId:  0,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		}
		parentRes, _ := ctrl.CreateDept(ctx, parentReq)

		childReq := &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptChild",
				ParentId:  parentRes.Id,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		}
		_, _ = ctrl.CreateDept(ctx, childReq)

		deleteReq := &v1.DeleteDeptReq{ID: parentRes.Id}
		res, err := ctrl.DeleteDept(ctx, deleteReq)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)
	})
}

func TestDeptController_GetDeptList(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysDept.Ctx(ctx).Unscoped().WhereLike(dao.SysDept.Columns().Name, "TestDept%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Prepare data: Create multiple departments
		_, _ = ctrl.CreateDept(ctx, &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptList1",
				ParentId:  0,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		})
		_, _ = ctrl.CreateDept(ctx, &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptList2",
				ParentId:  0,
				Status:    1,
				Order:     2,
				CreatorId: 1,
			},
		})

		// Test case 1: Get all departments
		req := &v1.GetDeptListReq{}
		res, err := ctrl.GetDeptList(ctx, req)
		t.AssertNil(err)
		t.AssertGE(res.Total, 2)
		t.AssertGE(len(res.List), 2)

		// Test case 2: Filter by name
		reqByName := &v1.GetDeptListReq{Name: "TestDeptList1"}
		resByName, err := ctrl.GetDeptList(ctx, reqByName)
		t.AssertNil(err)
		t.AssertGE(resByName.Total, 1)
		t.Assert(resByName.List[0].SysDept.Name, "TestDeptList1")

		// Test case 3: Filter by status
		reqByStatus := &v1.GetDeptListReq{Status: "1"}
		resByStatus, err := ctrl.GetDeptList(ctx, reqByStatus)
		t.AssertNil(err)
		t.AssertGE(resByStatus.Total, 2)
	})
}

func TestDeptController_GetDeptTree(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	// Cleanup function
	t.Cleanup(func() {
		dao.SysDept.Ctx(ctx).Unscoped().WhereLike(dao.SysDept.Columns().Name, "TestDept%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		// Prepare data: Create hierarchical departments
		parentRes, _ := ctrl.CreateDept(ctx, &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptTreeRoot",
				ParentId:  0,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		})

		_, _ = ctrl.CreateDept(ctx, &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptTreeChild",
				ParentId:  parentRes.Id,
				Status:    1,
				Order:     1,
				CreatorId: 1,
			},
		})

		// Test case 1: Get department tree
		req := &v1.GetDeptTreeReq{}
		res, err := ctrl.GetDeptTree(ctx, req)
		t.AssertNil(err)
		t.AssertGE(len(res.List), 1)

		// Find root in tree
		var root *model.SysDeptTreeOut
		for _, item := range res.List {
			if item.Id == parentRes.Id {
				root = item
				break
			}
		}
		t.AssertNE(root, nil)
		t.Assert(root.Name, "TestDeptTreeRoot")
		t.AssertGE(len(root.Children), 1)

		// Test case 2: Inactive departments should not appear in tree
		inactiveRes, _ := ctrl.CreateDept(ctx, &v1.CreateDeptReq{
			SysDeptCreateIn: model.SysDeptCreateIn{
				Name:      "TestDeptInactive",
				ParentId:  0,
				Status:    0,
				Order:     1,
				CreatorId: 1,
			},
		})

		res, _ = ctrl.GetDeptTree(ctx, req)
		// Verify inactive dept is not in tree
		foundInactive := false
		for _, item := range res.List {
			if item.Id == inactiveRes.Id {
				foundInactive = true
				break
			}
		}
		t.Assert(foundInactive, false)
	})
}
