package sys_dict_data

import (
	"context"
	"testing"

	v1 "backend/api/sys_dict_data/v1"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/service"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSysDictDataController_CRUD(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	t.Cleanup(func() {
		dao.SysDictData.Ctx(ctx).Unscoped().WhereLike(dao.SysDictData.Columns().Value, "test_dict_data_ctrl%").Delete()
		dao.SysDictType.Ctx(ctx).Unscoped().WhereLike(dao.SysDictType.Columns().TypeCode, "test_dict_data_type_ctrl%").Delete()
	})

	typeID, err := service.SysDictType().CreateDictType(ctx, model.SysDictTypeCreateIn{
		TypeCode: "test_dict_data_type_ctrl",
		TypeName: "Test Dict Data Type Ctrl",
		Status:   1,
	})
	if err != nil {
		t.Fatalf("failed to create dict type: %v", err)
	}

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		createReq := &v1.CreateDictDataReq{
			SysDictDataCreateIn: model.SysDictDataCreateIn{
				DictTypeId: typeID,
				Label:      "测试",
				Value:      "test_dict_data_ctrl_1",
				Status:     1,
				Sort:       1,
			},
		}
		createRes, err := ctrl.CreateDictData(ctx, createReq)
		t.AssertNil(err)
		t.AssertNE(createRes.Id, 0)

		_, err = ctrl.CreateDictData(ctx, createReq)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeValidationFailed)

		getRes, err := ctrl.GetDictData(ctx, &v1.GetDictDataReq{ID: createRes.Id})
		t.AssertNil(err)
		t.Assert(getRes.SysDictDataGetOut.SysDictDataDetail.Value, "test_dict_data_ctrl_1")

		updateReq := &v1.UpdateDictDataReq{
			ID: createRes.Id,
			SysDictDataUpdateIn: model.SysDictDataUpdateIn{
				DictTypeId: typeID,
				Label:      "测试更新",
				Value:      "test_dict_data_ctrl_1",
				Status:     1,
				Sort:       2,
			},
		}
		_, err = ctrl.UpdateDictData(ctx, updateReq)
		t.AssertNil(err)

		listRes, err := ctrl.GetDictDataList(ctx, &v1.GetDictDataListReq{DictTypeId: typeID})
		t.AssertNil(err)
		t.AssertGT(listRes.Total, 0)

		_, err = ctrl.DeleteDictData(ctx, &v1.DeleteDictDataReq{ID: createRes.Id})
		t.AssertNil(err)
	})
}
