package service

import (
	"context"
	"testing"

	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSysDictData_CRUD(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	t.Cleanup(func() {
		dao.SysDictData.Ctx(ctx).Unscoped().WhereLike(dao.SysDictData.Columns().Value, "test_dict_value%").Delete()
		dao.SysDictType.Ctx(ctx).Unscoped().WhereLike(dao.SysDictType.Columns().TypeCode, "test_dict_data_type%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		typeID, err := SysDictType().CreateDictType(ctx, model.SysDictTypeCreateIn{
			TypeCode: "test_dict_data_type",
			TypeName: "Test Dict Data Type",
			Status:   1,
		})
		t.AssertNil(err)

		createIn := model.SysDictDataCreateIn{
			DictTypeId: typeID,
			Label:      "测试",
			LabelI18n:  map[string]string{"en": "Test"},
			Value:      "test_dict_value_1",
			Status:     1,
			Sort:       1,
		}
		dataID, err := SysDictData().CreateDictData(ctx, createIn)
		t.AssertNil(err)
		t.AssertNE(dataID, 0)

		_, err = SysDictData().CreateDictData(ctx, createIn)
		t.AssertNE(err, nil)

		getOut, err := SysDictData().GetDictData(ctx, model.SysDictDataGetIn{Id: dataID})
		t.AssertNil(err)
		t.Assert(getOut.SysDictDataDetail.Label, "测试")
		t.Assert(getOut.SysDictDataDetail.LabelI18n["en"], "Test")

		updateIn := model.SysDictDataUpdateIn{
			Id:         dataID,
			DictTypeId: typeID,
			Label:      "测试更新",
			LabelI18n:  map[string]string{"en": "Test Updated"},
			Value:      "test_dict_value_1",
			Status:     1,
			Sort:       2,
		}
		err = SysDictData().UpdateDictData(ctx, updateIn)
		t.AssertNil(err)

		updated, err := SysDictData().GetDictData(ctx, model.SysDictDataGetIn{Id: dataID})
		t.AssertNil(err)
		t.Assert(updated.SysDictDataDetail.Label, "测试更新")

		listOut, err := SysDictData().GetDictDataList(ctx, model.SysDictDataListIn{DictTypeId: typeID})
		t.AssertNil(err)
		t.AssertGT(listOut.Total, 0)

		err = SysDictData().DeleteDictData(ctx, model.SysDictDataDeleteIn{Id: dataID})
		t.AssertNil(err)
	})
}

func TestSysDict_PublicOptionsAndLabel(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	t.Cleanup(func() {
		dao.SysDictData.Ctx(ctx).Unscoped().WhereLike(dao.SysDictData.Columns().Value, "test_dict_public_value%").Delete()
		dao.SysDictType.Ctx(ctx).Unscoped().WhereLike(dao.SysDictType.Columns().TypeCode, "test_dict_public%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		typeID, err := SysDictType().CreateDictType(ctx, model.SysDictTypeCreateIn{
			TypeCode: "test_dict_public",
			TypeName: "Test Dict Public",
			Status:   1,
		})
		t.AssertNil(err)

		_, err = SysDictData().CreateDictData(ctx, model.SysDictDataCreateIn{
			DictTypeId: typeID,
			Label:      "中文",
			LabelI18n:  map[string]string{"en": "English"},
			Value:      "test_dict_public_value",
			Status:     1,
			Sort:       1,
		})
		t.AssertNil(err)

		options, err := SysDict().GetDictOptions(ctx, "test_dict_public")
		t.AssertNil(err)
		t.Assert(len(options) > 0, true)
		t.Assert(options[0].Label, "中文")

		label, err := SysDict().GetDictLabel(ctx, "test_dict_public", "test_dict_public_value")
		t.AssertNil(err)
		t.Assert(label, "中文")
	})
}

func TestResolveDictLabel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		label := resolveDictLabel("默认", map[string]string{"en": "Default", "zh-tw": "預設"}, "en")
		t.Assert(label, "Default")

		label = resolveDictLabel("默认", map[string]string{"en": "Default", "zh-tw": "預設"}, "zh-TW")
		t.Assert(label, "預設")

		label = resolveDictLabel("默认", map[string]string{"en": "Default"}, "fr")
		t.Assert(label, "默认")
	})
}
