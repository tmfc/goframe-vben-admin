package model

import (
	"backend/internal/model/entity"
)

// SysMenuCreateIn is the input for creating a new menu.
type SysMenuCreateIn struct {
	Name           string `json:"name" v:"required#名称不能为空"`
	Path           string `json:"path"`
	Component      string `json:"component"`
	Icon           string `json:"icon"`
	Type           string `json:"type" v:"required#类型不能为空"`
	ParentId       string `json:"parentId"`
	Meta           string `json:"meta"`
	Status         int    `json:"status"`
	Order          int    `json:"order"`
	PermissionCode string `json:"permissionCode"`
}

// SysMenuGetOut is the output for retrieving a menu.
type SysMenuGetOut struct {
	*entity.SysMenu
}

// SysMenuUpdateIn is the input for updating a menu.
type SysMenuUpdateIn struct {
	ID             string `json:"id" v:"required#ID不能为空"`
	Name           string `json:"name" v:"required#名称不能为空"`
	Path           string `json:"path"`
	Component      string `json:"component"`
	Icon           string `json:"icon"`
	Type           string `json:"type" v:"required#类型不能为空"`
	ParentId       string `json:"parentId"`
	Meta           string `json:"meta"`
	Status         int    `json:"status"`
	Order          int    `json:"order"`
	PermissionCode string `json:"permissionCode"`
}

type SysMenuGetListIn struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SysMenuGetListOut struct {
	List  []*SysMenuGetOut `json:"list"`
	Total int              `json:"total"`
}
