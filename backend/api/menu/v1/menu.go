package v1

import "github.com/gogf/gf/v2/frame/g"

// MenuMeta describes menu metadata expected by the frontend.
type MenuMeta struct {
	Icon                     string `json:"icon,omitempty"`
	Title                    string `json:"title,omitempty"`
	Order                    int    `json:"order,omitempty"`
	AffixTab                 bool   `json:"affixTab,omitempty"`
	Badge                    string `json:"badge,omitempty"`
	BadgeType                string `json:"badgeType,omitempty"`
	BadgeVariants            string `json:"badgeVariants,omitempty"`
	IframeSrc                string `json:"iframeSrc,omitempty"`
	Link                     string `json:"link,omitempty"`
	MenuVisibleWithForbidden bool   `json:"menuVisibleWithForbidden,omitempty"`
}

// MenuItem represents a backend menu item for dynamic routing.
type MenuItem struct {
	Id        int         `json:"id,omitempty"`
	Pid       int         `json:"pid,omitempty"`
	Name      string      `json:"name,omitempty"`
	Path      string      `json:"path,omitempty"`
	Component string      `json:"component,omitempty"`
	Type      string      `json:"type,omitempty"`
	Status    int         `json:"status,omitempty"`
	Icon      string      `json:"icon,omitempty"`
	AuthCode  string      `json:"authCode,omitempty"`
	Meta      *MenuMeta   `json:"meta,omitempty"`
	Children  []*MenuItem `json:"children,omitempty"`
}

// MenuAllReq defines the request structure for fetching menu list.
type MenuAllReq struct {
	g.Meta `path:"/menu/all" method:"get" summary:"Fetch all menus" tags:"Menu"`
}

// MenuAllRes defines the response structure for fetching menu list.
type MenuAllRes []*MenuItem
