package sys_dict

import (
	"context"
	"testing"

	v1 "backend/api/sys_dict/v1"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/service"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSysDictController_PublicEndpoints(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	t.Cleanup(func() {
		dao.SysDictData.Ctx(ctx).Unscoped().WhereLike(dao.SysDictData.Columns().Value, "test_dict_public_ctrl%").Delete()
		dao.SysDictType.Ctx(ctx).Unscoped().WhereLike(dao.SysDictType.Columns().TypeCode, "test_dict_public_ctrl%").Delete()
	})

	typeID, err := service.SysDictType().CreateDictType(ctx, model.SysDictTypeCreateIn{
		TypeCode: "test_dict_public_ctrl",
		TypeName: "Test Dict Public Ctrl",
		Status:   1,
	})
	if err != nil {
		t.Fatalf("failed to create dict type: %v", err)
	}

	_, err = service.SysDictData().CreateDictData(ctx, model.SysDictDataCreateIn{
		DictTypeId: typeID,
		Label:      "中文",
		LabelI18n:  map[string]string{"en": "English"},
		Value:      "test_dict_public_ctrl_1",
		Status:     1,
		Sort:       1,
	})
	if err != nil {
		t.Fatalf("failed to create dict data: %v", err)
	}

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		optionsRes, err := ctrl.GetDictOptions(ctx, &v1.GetDictOptionsReq{TypeCode: "test_dict_public_ctrl"})
		t.AssertNil(err)
		t.Assert(len(optionsRes.List) > 0, true)
		t.Assert(optionsRes.List[0].Label, "中文")

		labelRes, err := ctrl.GetDictLabel(ctx, &v1.GetDictLabelReq{TypeCode: "test_dict_public_ctrl", Value: "test_dict_public_ctrl_1"})
		t.AssertNil(err)
		t.Assert(labelRes.Label, "中文")
	})
}
