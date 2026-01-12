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

var localSysDictData ISysDictData

// SysDictData returns the sys_dict_data service.
func SysDictData() ISysDictData {
	return localSysDictData
}

// RegisterSysDictData registers a sys_dict_data service implementation.
func RegisterSysDictData(i ISysDictData) {
	localSysDictData = i
}

type sSysDictData struct{}

func init() {
	RegisterSysDictData(NewSysDictData())
}

// NewSysDictData creates a new sys_dict_data service.
func NewSysDictData() *sSysDictData {
	return &sSysDictData{}
}

// ISysDictData defines the interface for sys_dict_data service.
type ISysDictData interface {
	CreateDictData(ctx context.Context, in model.SysDictDataCreateIn) (id int64, err error)
	GetDictData(ctx context.Context, in model.SysDictDataGetIn) (out *model.SysDictDataGetOut, err error)
	UpdateDictData(ctx context.Context, in model.SysDictDataUpdateIn) (err error)
	DeleteDictData(ctx context.Context, in model.SysDictDataDeleteIn) (err error)
	GetDictDataList(ctx context.Context, in model.SysDictDataListIn) (out *model.SysDictDataListOut, err error)
}

// CreateDictData creates a new dict data.
func (s *sSysDictData) CreateDictData(ctx context.Context, in model.SysDictDataCreateIn) (id int64, err error) {
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return 0, err
	}

	dictType, err := fetchDictTypeByID(ctx, in.DictTypeId)
	if err != nil {
		return 0, err
	}

	columns := dao.SysDictData.Columns()
	count, err := dao.SysDictData.Ctx(ctx).
		Where(columns.DictTypeId, in.DictTypeId).
		Where(columns.Value, in.Value).
		Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, gerror.NewCodef(gcode.CodeValidationFailed, "Dict value '%s' already exists", in.Value)
	}

	labelI18n, err := encodeLabelI18n(in.LabelI18n)
	if err != nil {
		return 0, err
	}
	var labelI18nValue any = labelI18n
	if strings.TrimSpace(labelI18n) == "" {
		labelI18nValue = nil
	}

	creatorID := in.CreatorId
	if creatorID == 0 {
		creatorID = resolveActorID(ctx)
	}
	modifierID := in.ModifierId
	if modifierID == 0 {
		modifierID = creatorID
	}
	result, err := dao.SysDictData.Ctx(ctx).Data(g.Map{
		columns.DictTypeId:  in.DictTypeId,
		columns.Label:       in.Label,
		columns.LabelI18N:   labelI18nValue,
		columns.Value:       in.Value,
		columns.Description: in.Description,
		columns.Color:       in.Color,
		columns.Icon:        in.Icon,
		columns.CssClass:    in.CssClass,
		columns.Status:      in.Status,
		columns.Sort:        in.Sort,
		columns.IsDefault:   in.IsDefault,
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

	clearDictCache(dictCacheKey(resolveTenantID(ctx), dictType.TypeCode))
	return lastInsertId, nil
}

// GetDictData retrieves a dict data by ID.
func (s *sSysDictData) GetDictData(ctx context.Context, in model.SysDictDataGetIn) (out *model.SysDictDataGetOut, err error) {
	var data entity.SysDictData
	if err = dao.SysDictData.Ctx(ctx).
		Where(dao.SysDictData.Columns().Id, in.Id).
		Scan(&data); err != nil {
		return nil, err
	}
	if data.Id == 0 {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Dict data with ID %d not found", in.Id)
	}
	out = &model.SysDictDataGetOut{SysDictDataDetail: toDictDataDetail(&data)}
	return out, nil
}

// UpdateDictData updates a dict data.
func (s *sSysDictData) UpdateDictData(ctx context.Context, in model.SysDictDataUpdateIn) (err error) {
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}

	columns := dao.SysDictData.Columns()
	var existing entity.SysDictData
	if err = dao.SysDictData.Ctx(ctx).
		Where(columns.Id, in.Id).
		Scan(&existing); err != nil {
		return err
	}
	if existing.Id == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Dict data with ID %d not found", in.Id)
	}

	newType, err := fetchDictTypeByID(ctx, in.DictTypeId)
	if err != nil {
		return err
	}

	count, err := dao.SysDictData.Ctx(ctx).
		Where(columns.DictTypeId, in.DictTypeId).
		Where(columns.Value, in.Value).
		WhereNot(columns.Id, in.Id).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Dict value '%s' already exists", in.Value)
	}

	labelI18n, err := encodeLabelI18n(in.LabelI18n)
	if err != nil {
		return err
	}
	var labelI18nValue any = labelI18n
	if strings.TrimSpace(labelI18n) == "" {
		labelI18nValue = nil
	}

	modifierID := in.ModifierId
	if modifierID == 0 {
		modifierID = resolveActorID(ctx)
	}
	updateData := g.Map{
		columns.DictTypeId:  in.DictTypeId,
		columns.Label:       in.Label,
		columns.LabelI18N:   labelI18nValue,
		columns.Value:       in.Value,
		columns.Description: in.Description,
		columns.Color:       in.Color,
		columns.Icon:        in.Icon,
		columns.CssClass:    in.CssClass,
		columns.Status:      in.Status,
		columns.Sort:        in.Sort,
		columns.IsDefault:   in.IsDefault,
		columns.ModifierId:  modifierID,
	}
	if in.UpdatedAt != nil {
		updateData[columns.UpdatedAt] = in.UpdatedAt
	}

	_, err = dao.SysDictData.Ctx(ctx).
		Data(updateData).
		Where(columns.Id, in.Id).
		Update()
	if err != nil {
		return err
	}

	tenantID := resolveTenantID(ctx)
	clearDictCache(dictCacheKey(tenantID, newType.TypeCode))
	if existing.DictTypeId != in.DictTypeId {
		if oldType, err := fetchDictTypeByID(ctx, existing.DictTypeId); err == nil {
			clearDictCache(dictCacheKey(tenantID, oldType.TypeCode))
		}
	}
	return nil
}

// DeleteDictData deletes a dict data by ID.
func (s *sSysDictData) DeleteDictData(ctx context.Context, in model.SysDictDataDeleteIn) (err error) {
	columns := dao.SysDictData.Columns()
	var existing entity.SysDictData
	if err = dao.SysDictData.Ctx(ctx).
		Where(columns.Id, in.Id).
		Scan(&existing); err != nil {
		return err
	}
	if existing.Id == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Dict data with ID %d not found", in.Id)
	}

	_, err = dao.SysDictData.Ctx(ctx).
		Where(columns.Id, in.Id).
		Delete()
	if err != nil {
		return err
	}

	if dictType, err := fetchDictTypeByID(ctx, existing.DictTypeId); err == nil {
		clearDictCache(dictCacheKey(resolveTenantID(ctx), dictType.TypeCode))
	}
	return nil
}

// GetDictDataList lists dict data with pagination and filters.
func (s *sSysDictData) GetDictDataList(ctx context.Context, in model.SysDictDataListIn) (out *model.SysDictDataListOut, err error) {
	out = &model.SysDictDataListOut{}

	dictTypeId := in.DictTypeId
	if dictTypeId == 0 && strings.TrimSpace(in.TypeCode) != "" {
		typeInfo, err := fetchDictTypeByCode(ctx, in.TypeCode)
		if err != nil {
			if gerror.Code(err) == gcode.CodeNotFound {
				return &model.SysDictDataListOut{Items: []*model.SysDictDataDetail{}, Total: 0}, nil
			}
			return nil, err
		}
		dictTypeId = typeInfo.Id
	}

	m := dao.SysDictData.Ctx(ctx)
	if dictTypeId > 0 {
		m = m.Where(dao.SysDictData.Columns().DictTypeId, dictTypeId)
	}
	if strings.TrimSpace(in.Label) != "" {
		m = m.WhereLike(dao.SysDictData.Columns().Label, "%"+strings.TrimSpace(in.Label)+"%")
	}
	if strings.TrimSpace(in.Value) != "" {
		m = m.WhereLike(dao.SysDictData.Columns().Value, "%"+strings.TrimSpace(in.Value)+"%")
	}
	if strings.TrimSpace(in.Status) != "" {
		m = m.Where(dao.SysDictData.Columns().Status, in.Status)
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

	var items []entity.SysDictData
	if err = m.Order(dao.SysDictData.Columns().Sort+" asc", dao.SysDictData.Columns().Id+" asc").
		Page(in.Page, in.PageSize).
		Scan(&items); err != nil {
		return nil, err
	}

	out.Items = make([]*model.SysDictDataDetail, 0, len(items))
	for i := range items {
		out.Items = append(out.Items, toDictDataDetail(&items[i]))
	}
	return out, nil
}

func toDictDataDetail(data *entity.SysDictData) *model.SysDictDataDetail {
	if data == nil {
		return nil
	}
	return &model.SysDictDataDetail{
		Id:          data.Id,
		TenantId:    data.TenantId,
		DictTypeId:  data.DictTypeId,
		Label:       data.Label,
		LabelI18n:   parseLabelI18n(data.LabelI18N),
		Value:       data.Value,
		Description: data.Description,
		Color:       data.Color,
		Icon:        data.Icon,
		CssClass:    data.CssClass,
		Status:      data.Status,
		Sort:        data.Sort,
		IsDefault:   data.IsDefault,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
		CreatorId:   data.CreatorId,
		ModifierId:  data.ModifierId,
		DeptId:      data.DeptId,
	}
}
