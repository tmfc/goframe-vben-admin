package service

import (
	"context"
	"strings"

	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

var localSysDict ISysDict

// SysDict returns the sys_dict service.
func SysDict() ISysDict {
	return localSysDict
}

// RegisterSysDict registers a sys_dict service implementation.
func RegisterSysDict(i ISysDict) {
	localSysDict = i
}

type sSysDict struct{}

func init() {
	RegisterSysDict(NewSysDict())
}

// NewSysDict creates a new sys_dict service.
func NewSysDict() *sSysDict {
	return &sSysDict{}
}

// ISysDict defines the interface for sys_dict service.
type ISysDict interface {
	GetDictOptions(ctx context.Context, typeCode string) (out []model.SysDictOption, err error)
	GetDictLabel(ctx context.Context, typeCode, value string) (label string, err error)
}

// GetDictOptions returns dict options by type code.
func (s *sSysDict) GetDictOptions(ctx context.Context, typeCode string) (out []model.SysDictOption, err error) {
	lang := resolveAcceptLanguage(ctx)
	tenantID := resolveTenantID(ctx)
	cacheKey := dictCacheKey(tenantID, typeCode)
	if cached, ok := getDictCache(cacheKey); ok {
		return buildDictOptions(cached, lang), nil
	}

	dictType, err := fetchDictTypeByCode(ctx, typeCode)
	if err != nil {
		if gerror.Code(err) == gcode.CodeNotFound {
			return []model.SysDictOption{}, nil
		}
		return nil, err
	}
	if dictType.Status != 1 {
		return []model.SysDictOption{}, nil
	}

	var items []entity.SysDictData
	err = dao.SysDictData.Ctx(ctx).
		Where(dao.SysDictData.Columns().DictTypeId, dictType.Id).
		Where(dao.SysDictData.Columns().Status, 1).
		Order(dao.SysDictData.Columns().Sort+" asc", dao.SysDictData.Columns().Id+" asc").
		Scan(&items)
	if err != nil {
		return nil, err
	}

	cacheItems := make([]dictDataCacheItem, 0, len(items))
	for i := range items {
		cacheItems = append(cacheItems, dictDataCacheItem{
			Label:     items[i].Label,
			LabelI18n: parseLabelI18n(items[i].LabelI18n),
			Value:     items[i].Value,
			Color:     items[i].Color,
			Icon:      items[i].Icon,
			IsDefault: items[i].IsDefault,
		})
	}
	setDictCache(cacheKey, cacheItems)
	return buildDictOptions(cacheItems, lang), nil
}

// GetDictLabel returns the dict label by type code and value.
func (s *sSysDict) GetDictLabel(ctx context.Context, typeCode, value string) (label string, err error) {
	lang := resolveAcceptLanguage(ctx)
	tenantID := resolveTenantID(ctx)
	cacheKey := dictCacheKey(tenantID, typeCode)
	cached, ok := getDictCache(cacheKey)
	if !ok {
		_, err = s.GetDictOptions(ctx, typeCode)
		if err != nil {
			return "", err
		}
		cached, _ = getDictCache(cacheKey)
	}

	value = strings.TrimSpace(value)
	for _, item := range cached {
		if strings.TrimSpace(item.Value) == value {
			return resolveDictLabel(item.Label, item.LabelI18n, lang), nil
		}
	}
	return "", gerror.NewCodef(gcode.CodeNotFound, "Dict value '%s' not found", value)
}

func buildDictOptions(items []dictDataCacheItem, lang string) []model.SysDictOption {
	options := make([]model.SysDictOption, 0, len(items))
	for _, item := range items {
		options = append(options, model.SysDictOption{
			Label:     resolveDictLabel(item.Label, item.LabelI18n, lang),
			Value:     item.Value,
			Color:     item.Color,
			Icon:      item.Icon,
			IsDefault: item.IsDefault,
		})
	}
	return options
}
