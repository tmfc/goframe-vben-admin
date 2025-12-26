package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"backend/internal/model"
)

// CreateMenuReq defines the request structure for creating a new menu.
type CreateMenuReq struct {
	g.Meta `path:"/sys-menu" method:"post" summary:"Create a new menu" tags:"System Menu"`
	model.SysMenuCreateIn
}

// CreateMenuRes defines the response structure for creating a new menu.
type CreateMenuRes struct {
	Id uint `json:"id"`
}

// GetMenuReq defines the request structure for retrieving a menu.
type GetMenuReq struct {
	g.Meta `path:"/sys-menu/{id}" method:"get" summary:"Retrieve a menu by ID" tags:"System Menu"`
	ID     uint `json:"id" v:"required#ID不能为空"`
}

// GetMenuRes defines the response structure for retrieving a menu.
type GetMenuRes struct {
	*model.SysMenuGetOut
}

// UpdateMenuReq defines the request structure for updating a menu.
type UpdateMenuReq struct {
	g.Meta `path:"/sys-menu/{id}" method:"put" summary:"Update a menu by ID" tags:"System Menu"`
	ID     uint `json:"id" v:"required#ID不能为空"`
	model.SysMenuUpdateIn
}

// UpdateMenuRes defines the response structure for updating a menu.
type UpdateMenuRes struct{}

// DeleteMenuReq defines the request structure for deleting a menu.
type DeleteMenuReq struct {
	g.Meta `path:"/sys-menu/{id}" method:"delete" summary:"Delete a menu by ID" tags:"System Menu"`
	ID     uint `json:"id" v:"required#ID不能为空"`
}

// DeleteMenuRes defines the response structure for deleting a menu.
type DeleteMenuRes struct{}

// GetMenuListReq defines the request structure for retrieving menu list.
type GetMenuListReq struct {
	g.Meta `path:"/sys-menu/list" method:"get" summary:"Retrieve menu list" tags:"System Menu"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// GetMenuListRes defines the response structure for retrieving menu list.
type GetMenuListRes struct {
	List  []*model.SysMenuGetOut `json:"list"`
	Total int                    `json:"total"`
}

// MenuAllReq defines the request structure for retrieving all menus.
type MenuAllReq struct {
	g.Meta `path:"/menu/all" method:"get" summary:"Get all menus" tags:"System Menu"`
}

// MenuAllRes defines the response structure for retrieving all menus.
type MenuAllRes []*MenuItem

type MenuItem struct {
	Id        int         `json:"id"`
	Pid       int         `json:"pid,omitempty"`
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Component string      `json:"component,omitempty"`
	Type      string      `json:"type"`
	Status    int         `json:"status"`
	Icon      string      `json:"icon,omitempty"`
	AuthCode  string      `json:"authCode,omitempty"`
	Meta      *MenuMeta   `json:"meta,omitempty"`
	Children  []*MenuItem `json:"children,omitempty"`
}

type MenuMeta struct {
	Icon          string `json:"icon,omitempty"`
	Title         string `json:"title,omitempty"`
	AffixTab      bool   `json:"affixTab,omitempty"`
	Order         int    `json:"order,omitempty"`
	Badge         string `json:"badge,omitempty"`
	BadgeType     string `json:"badgeType,omitempty"`
	BadgeVariants string `json:"badgeVariants,omitempty"`
	IframeSrc     string `json:"iframeSrc,omitempty"`
	Link          string `json:"link,omitempty"`
}