package service

import (
	"strings"
	"sync"
)

type dictDataCacheItem struct {
	Label     string
	LabelI18n map[string]string
	Value     string
	Color     string
	Icon      string
	IsDefault bool
}

var dictDataCache sync.Map

func dictCacheKey(tenantID, typeCode string) string {
	return strings.TrimSpace(tenantID) + ":" + strings.ToLower(strings.TrimSpace(typeCode))
}

func getDictCache(key string) ([]dictDataCacheItem, bool) {
	value, ok := dictDataCache.Load(key)
	if !ok {
		return nil, false
	}
	items, ok := value.([]dictDataCacheItem)
	return items, ok
}

func setDictCache(key string, items []dictDataCacheItem) {
	dictDataCache.Store(key, items)
}

func clearDictCache(key string) {
	dictDataCache.Delete(key)
}
