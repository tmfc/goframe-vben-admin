package v1

import (
	"backend/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// CreateMenuReq defines the request structure for creating a new menu.
type CreateMenuReq struct {
	g.Meta `path:"/sys-menu" method:"post" summary:"Create a new menu" tags:"System Menu"`
	model.SysMenuCreateIn
}

// CreateMenuRes defines the response structure for creating a new menu.
type CreateMenuRes struct {
	Id int64 `json:"id"`
}

// GenerateButtonsReq defines the request structure for generating default buttons under a menu.
type GenerateButtonsReq struct {
	g.Meta `path:"/sys-menu/{id}/generate-buttons" method:"post" summary:"Generate default buttons for a menu" tags:"System Menu"`
	Id     int64 `json:"id" in:"path" v:"min:1#菜单ID不能为空"`
}

// GenerateButtonsRes defines the response structure for generated buttons.
type GenerateButtonsRes struct {
	Created int `json:"created"`
	Skipped int `json:"skipped"`
}

// GetMenuReq defines the request structure for retrieving a menu.
type GetMenuReq struct {
	g.Meta `path:"/sys-menu/{id}" method:"get" summary:"Retrieve a menu by ID" tags:"System Menu"`
	ID     int64 `json:"id" v:"min:1#ID不能为空"`
}

// GetMenuRes defines the response structure for retrieving a menu.
type GetMenuRes struct {
	*model.SysMenuGetOut
}

// UpdateMenuReq defines the request structure for updating a menu.
type UpdateMenuReq struct {
	g.Meta `path:"/sys-menu/{id}" method:"put" summary:"Update a menu by ID" tags:"System Menu"`
	ID     int64 `json:"id" v:"min:1#ID不能为空"`
	model.SysMenuUpdateIn
}

// UpdateMenuRes defines the response structure for updating a menu.
type UpdateMenuRes struct{}

// DeleteMenuReq defines the request structure for deleting a menu.
type DeleteMenuReq struct {
	g.Meta `path:"/sys-menu/{id}" method:"delete" summary:"Delete a menu by ID" tags:"System Menu"`
	ID     int64 `json:"id" v:"min:1#ID不能为空"`
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
	Id        int64       `json:"id"`
	Pid       int64       `json:"pid,omitempty"`
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
	Order         int    `json:"order"`
	Badge         string `json:"badge,omitempty"`
	BadgeType     string `json:"badgeType,omitempty"`
	BadgeVariants string `json:"badgeVariants,omitempty"`
	IframeSrc     string `json:"iframeSrc,omitempty"`
	Link          string `json:"link,omitempty"`
}
