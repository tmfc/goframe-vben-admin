package service

import (
	"context"
	"testing"

	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSysDictType_CRUD(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	t.Cleanup(func() {
		dao.SysDictType.Ctx(ctx).Unscoped().WhereLike(dao.SysDictType.Columns().TypeCode, "test_dict_type%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		createIn := model.SysDictTypeCreateIn{
			TypeCode:    "test_dict_type_1",
			TypeName:    "Test Dict Type",
			Description: "Test",
			Status:      1,
			Sort:        1,
		}
		id, err := SysDictType().CreateDictType(ctx, createIn)
		t.AssertNil(err)
		t.AssertNE(id, 0)

		_, err = SysDictType().CreateDictType(ctx, createIn)
		t.AssertNE(err, nil)

		getOut, err := SysDictType().GetDictType(ctx, model.SysDictTypeGetIn{Id: id})
		t.AssertNil(err)
		t.Assert(getOut.SysDictType.TypeCode, "test_dict_type_1")

		updateIn := model.SysDictTypeUpdateIn{
			Id:          id,
			TypeCode:    "test_dict_type_1_updated",
			TypeName:    "Test Dict Type Updated",
			Description: "Updated",
			Status:      1,
			Sort:        2,
		}
		err = SysDictType().UpdateDictType(ctx, updateIn)
		t.AssertNil(err)

		updated, err := SysDictType().GetDictType(ctx, model.SysDictTypeGetIn{Id: id})
		t.AssertNil(err)
		t.Assert(updated.SysDictType.TypeName, "Test Dict Type Updated")

		listOut, err := SysDictType().GetDictTypeList(ctx, model.SysDictTypeListIn{TypeCode: "test_dict_type"})
		t.AssertNil(err)
		t.AssertGT(listOut.Total, 0)

		err = SysDictType().DeleteDictType(ctx, model.SysDictTypeDeleteIn{Id: id})
		t.AssertNil(err)
	})
}

func TestSysDictType_DeleteSystem(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()

	t.Cleanup(func() {
		dao.SysDictType.Ctx(ctx).Unscoped().WhereLike(dao.SysDictType.Columns().TypeCode, "test_dict_type_system%").Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		id, err := SysDictType().CreateDictType(ctx, model.SysDictTypeCreateIn{
			TypeCode: "test_dict_type_system",
			TypeName: "System Dict",
			IsSystem: true,
			Status:   1,
		})
		t.AssertNil(err)

		var existing entity.SysDictType
		err = dao.SysDictType.Ctx(ctx).Where(dao.SysDictType.Columns().Id, id).Scan(&existing)
		t.AssertNil(err)
		t.Assert(existing.IsSystem, true)

		err = SysDictType().DeleteDictType(ctx, model.SysDictTypeDeleteIn{Id: id})
		t.AssertNE(err, nil)
	})
}
