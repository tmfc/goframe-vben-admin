package service

import (
	"context"
	"strings"

	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var localSysDictType ISysDictType

// SysDictType returns the sys_dict_type service.
func SysDictType() ISysDictType {
	return localSysDictType
}

// RegisterSysDictType registers a sys_dict_type service implementation.
func RegisterSysDictType(i ISysDictType) {
	localSysDictType = i
}

type sSysDictType struct{}

func init() {
	RegisterSysDictType(NewSysDictType())
}

// NewSysDictType creates a new sys_dict_type service.
func NewSysDictType() *sSysDictType {
	return &sSysDictType{}
}

// ISysDictType defines the interface for sys_dict_type service.
type ISysDictType interface {
	CreateDictType(ctx context.Context, in model.SysDictTypeCreateIn) (id int64, err error)
	GetDictType(ctx context.Context, in model.SysDictTypeGetIn) (out *model.SysDictTypeGetOut, err error)
	UpdateDictType(ctx context.Context, in model.SysDictTypeUpdateIn) (err error)
	DeleteDictType(ctx context.Context, in model.SysDictTypeDeleteIn) (err error)
	GetDictTypeList(ctx context.Context, in model.SysDictTypeListIn) (out *model.SysDictTypeListOut, err error)
}

// CreateDictType creates a new dict type.
func (s *sSysDictType) CreateDictType(ctx context.Context, in model.SysDictTypeCreateIn) (id int64, err error) {
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return 0, err
	}

	columns := dao.SysDictType.Columns()
	count, err := dao.SysDictType.Ctx(ctx).
		Where(columns.TypeCode, in.TypeCode).
		Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, gerror.NewCodef(gcode.CodeValidationFailed, "Dict type code '%s' already exists", in.TypeCode)
	}

	creatorID := in.CreatorId
	if creatorID == 0 {
		creatorID = resolveActorID(ctx)
	}
	modifierID := in.ModifierId
	if modifierID == 0 {
		modifierID = creatorID
	}
	result, err := dao.SysDictType.Ctx(ctx).Data(g.Map{
		columns.TypeCode:    in.TypeCode,
		columns.TypeName:    in.TypeName,
		columns.Description: in.Description,
		columns.IsSystem:    in.IsSystem,
		columns.Status:      in.Status,
		columns.Sort:        in.Sort,
		columns.TenantId:    resolveTenantID(ctx),
		columns.CreatorId:   creatorID,
		columns.ModifierId:  modifierID,
		columns.DeptId:      in.DeptId,
	}).Insert()
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

// GetDictType retrieves a dict type by ID.
func (s *sSysDictType) GetDictType(ctx context.Context, in model.SysDictTypeGetIn) (out *model.SysDictTypeGetOut, err error) {
	out = &model.SysDictTypeGetOut{}
	err = dao.SysDictType.Ctx(ctx).
		Where(dao.SysDictType.Columns().Id, in.Id).
		Scan(&out.SysDictType)
	if err != nil {
		return nil, err
	}
	if out.SysDictType == nil {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Dict type with ID %d not found", in.Id)
	}
	return out, nil
}

// UpdateDictType updates a dict type.
func (s *sSysDictType) UpdateDictType(ctx context.Context, in model.SysDictTypeUpdateIn) (err error) {
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}

	columns := dao.SysDictType.Columns()
	var existing entity.SysDictType
	if err = dao.SysDictType.Ctx(ctx).
		Where(columns.Id, in.Id).
		Scan(&existing); err != nil {
		return err
	}
	if existing.Id == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Dict type with ID %d not found", in.Id)
	}

	if strings.TrimSpace(existing.TypeCode) != strings.TrimSpace(in.TypeCode) {
		count, err := dao.SysDictType.Ctx(ctx).
			Where(columns.TypeCode, in.TypeCode).
			WhereNot(columns.Id, in.Id).
			Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return gerror.NewCodef(gcode.CodeValidationFailed, "Dict type code '%s' already exists", in.TypeCode)
		}
	}

	modifierID := in.ModifierId
	if modifierID == 0 {
		modifierID = resolveActorID(ctx)
	}
	updateData := g.Map{
		columns.TypeCode:    in.TypeCode,
		columns.TypeName:    in.TypeName,
		columns.Description: in.Description,
		columns.IsSystem:    in.IsSystem,
		columns.Status:      in.Status,
		columns.Sort:        in.Sort,
		columns.ModifierId:  modifierID,
	}
	if in.UpdatedAt != nil {
		updateData[columns.UpdatedAt] = in.UpdatedAt
	}

	_, err = dao.SysDictType.Ctx(ctx).
		Data(updateData).
		Where(columns.Id, in.Id).
		Update()
	if err != nil {
		return err
	}

	if strings.TrimSpace(existing.TypeCode) != strings.TrimSpace(in.TypeCode) {
		tenantID := resolveTenantID(ctx)
		clearDictCache(dictCacheKey(tenantID, existing.TypeCode))
		clearDictCache(dictCacheKey(tenantID, in.TypeCode))
	} else {
		clearDictCache(dictCacheKey(resolveTenantID(ctx), existing.TypeCode))
	}
	return nil
}

// DeleteDictType deletes a dict type by ID.
func (s *sSysDictType) DeleteDictType(ctx context.Context, in model.SysDictTypeDeleteIn) (err error) {
	columns := dao.SysDictType.Columns()
	var existing entity.SysDictType
	if err = dao.SysDictType.Ctx(ctx).
		Where(columns.Id, in.Id).
		Scan(&existing); err != nil {
		return err
	}
	if existing.Id == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Dict type with ID %d not found", in.Id)
	}
	if existing.IsSystem {
		return gerror.NewCodef(gcode.CodeValidationFailed, "System dict type cannot be deleted")
	}

	_, err = dao.SysDictType.Ctx(ctx).
		Where(columns.Id, in.Id).
		Delete()
	if err != nil {
		return err
	}
	clearDictCache(dictCacheKey(resolveTenantID(ctx), existing.TypeCode))
	return nil
}

// GetDictTypeList lists dict types with pagination and filters.
func (s *sSysDictType) GetDictTypeList(ctx context.Context, in model.SysDictTypeListIn) (out *model.SysDictTypeListOut, err error) {
	out = &model.SysDictTypeListOut{}

	m := dao.SysDictType.Ctx(ctx)
	if strings.TrimSpace(in.TypeCode) != "" {
		m = m.WhereLike(dao.SysDictType.Columns().TypeCode, "%"+strings.TrimSpace(in.TypeCode)+"%")
	}
	if strings.TrimSpace(in.TypeName) != "" {
		m = m.WhereLike(dao.SysDictType.Columns().TypeName, "%"+strings.TrimSpace(in.TypeName)+"%")
	}
	if strings.TrimSpace(in.Status) != "" {
		m = m.Where(dao.SysDictType.Columns().Status, in.Status)
	}

	out.Total, err = m.Count()
	if err != nil {
		return nil, err
	}

	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	err = m.Order(dao.SysDictType.Columns().Sort+" asc", dao.SysDictType.Columns().Id+" asc").
		Page(in.Page, in.PageSize).
		Scan(&out.Items)
	if err != nil {
		return nil, err
	}
	return out, nil
}
