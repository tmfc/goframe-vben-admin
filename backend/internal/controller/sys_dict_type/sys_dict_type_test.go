package sys_dict_type

import (
	"context"
	"testing"

	"backend/api/sys_dict_type/v1"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSysDictTypeController_CRUD(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	t.Cleanup(func() {
		dao.SysDictType.Ctx(ctx).Unscoped().WhereLike(dao.SysDictType.Columns().TypeCode, "test_dict_type_ctrl%").Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		createReq := &v1.CreateDictTypeReq{
			SysDictTypeCreateIn: model.SysDictTypeCreateIn{
				TypeCode: "test_dict_type_ctrl",
				TypeName: "Test Dict Type Ctrl",
				Status:   1,
				Sort:     1,
			},
		}
		createRes, err := ctrl.CreateDictType(ctx, createReq)
		t.AssertNil(err)
		t.AssertNE(createRes.Id, 0)

		_, err = ctrl.CreateDictType(ctx, createReq)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)

		getRes, err := ctrl.GetDictType(ctx, &v1.GetDictTypeReq{ID: createRes.Id})
		t.AssertNil(err)
		t.Assert(getRes.SysDictTypeGetOut.SysDictType.TypeCode, "test_dict_type_ctrl")

		updateReq := &v1.UpdateDictTypeReq{
			ID: createRes.Id,
			SysDictTypeUpdateIn: model.SysDictTypeUpdateIn{
				TypeCode: "test_dict_type_ctrl_updated",
				TypeName: "Test Dict Type Ctrl Updated",
				Status:   1,
				Sort:     2,
			},
		}
		_, err = ctrl.UpdateDictType(ctx, updateReq)
		t.AssertNil(err)

		listRes, err := ctrl.GetDictTypeList(ctx, &v1.GetDictTypeListReq{TypeCode: "test_dict_type_ctrl"})
		t.AssertNil(err)
		t.AssertGT(listRes.Total, 0)

		_, err = ctrl.DeleteDictType(ctx, &v1.DeleteDictTypeReq{ID: createRes.Id})
		t.AssertNil(err)
	})
}
