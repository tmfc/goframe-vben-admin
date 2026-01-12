package service

import (
	"context"
	"strings"

	"backend/internal/dao"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func fetchDictTypeByID(ctx context.Context, id int64) (*entity.SysDictType, error) {
	var dictType entity.SysDictType
	if err := dao.SysDictType.Ctx(ctx).
		Where(dao.SysDictType.Columns().Id, id).
		Scan(&dictType); err != nil {
		return nil, err
	}
	if dictType.Id == 0 {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Dict type with ID %d not found", id)
	}
	return &dictType, nil
}

func fetchDictTypeByCode(ctx context.Context, typeCode string) (*entity.SysDictType, error) {
	typeCode = strings.TrimSpace(typeCode)
	if typeCode == "" {
		return nil, gerror.NewCodef(gcode.CodeValidationFailed, "Dict type code cannot be empty")
	}
	var dictType entity.SysDictType
	if err := dao.SysDictType.Ctx(ctx).
		Where(dao.SysDictType.Columns().TypeCode, typeCode).
		Scan(&dictType); err != nil {
		return nil, err
	}
	if dictType.Id == 0 {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Dict type code '%s' not found", typeCode)
	}
	return &dictType, nil
}
