package service

import (
	"context"
	"encoding/json"
	"strings"

	"backend/api/menu/v1"

	"github.com/gogf/gf/v2/frame/g"
)

var localMenu IMenu

// Menu returns the menu service instance.
func Menu() IMenu {
	return localMenu
}

// RegisterMenu sets the instance used by menu related handlers.
func RegisterMenu(i IMenu) {
	localMenu = i
}

var _ IMenu = (*sMenu)(nil)

func init() {
	RegisterMenu(NewMenu())
}

// NewMenu creates a new menu service instance.
func NewMenu() *sMenu {
	return &sMenu{}
}

// IMenu defines the menu service interface.
type IMenu interface {
	All(ctx context.Context) (v1.MenuAllRes, error)
}

type sMenu struct{}

// All returns the menu tree for the current tenant.
func (s *sMenu) All(ctx context.Context) (v1.MenuAllRes, error) {
	menus, err := fetchMenuFromDB(ctx)
	if err != nil || len(menus) == 0 {
		return filterMenuRoutes(defaultMenuList()), nil
	}
	return filterMenuRoutes(menus), nil
}

func defaultMenuList() v1.MenuAllRes {
	return v1.MenuAllRes{
		{
			Id:        1,
			Name:      "Workspace",
			Status:    1,
			Type:      "menu",
			Icon:      "mdi:dashboard",
			Path:      "/workspace",
			Component: "/dashboard/workspace/index",
			Meta: &v1.MenuMeta{
				Icon:     "carbon:workspace",
				Title:    "page.dashboard.workspace",
				AffixTab: true,
				Order:    0,
			},
		},
		{
			Id:     2,
			Name:   "System",
			Status: 1,
			Type:   "catalog",
			Path:   "/system",
			Meta: &v1.MenuMeta{
				Icon:          "carbon:settings",
				Order:         9997,
				Title:         "system.title",
				Badge:         "new",
				BadgeType:     "normal",
				BadgeVariants: "primary",
			},
			Children: []*v1.MenuItem{
				{
					Id:       201,
					Pid:      2,
					Path:     "/system/menu",
					Name:     "SystemMenu",
					AuthCode: "System:Menu:List",
					Status:   1,
					Type:     "menu",
					Meta: &v1.MenuMeta{
						Icon:  "carbon:menu",
						Title: "system.menu.title",
					},
					Component: "/system/menu/list",
					Children: []*v1.MenuItem{
						{
							Id:       20101,
							Pid:      201,
							Name:     "SystemMenuCreate",
							Status:   1,
							Type:     "button",
							AuthCode: "System:Menu:Create",
							Meta: &v1.MenuMeta{
								Title: "common.create",
							},
						},
						{
							Id:       20102,
							Pid:      201,
							Name:     "SystemMenuEdit",
							Status:   1,
							Type:     "button",
							AuthCode: "System:Menu:Edit",
							Meta: &v1.MenuMeta{
								Title: "common.edit",
							},
						},
						{
							Id:       20103,
							Pid:      201,
							Name:     "SystemMenuDelete",
							Status:   1,
							Type:     "button",
							AuthCode: "System:Menu:Delete",
							Meta: &v1.MenuMeta{
								Title: "common.delete",
							},
						},
					},
				},
				{
					Id:       202,
					Pid:      2,
					Path:     "/system/dept",
					Name:     "SystemDept",
					Status:   1,
					Type:     "menu",
					AuthCode: "System:Dept:List",
					Meta: &v1.MenuMeta{
						Icon:  "carbon:container-services",
						Title: "system.dept.title",
					},
					Component: "/system/dept/list",
					Children: []*v1.MenuItem{
						{
							Id:       20401,
							Pid:      202,
							Name:     "SystemDeptCreate",
							Status:   1,
							Type:     "button",
							AuthCode: "System:Dept:Create",
							Meta: &v1.MenuMeta{
								Title: "common.create",
							},
						},
						{
							Id:       20402,
							Pid:      202,
							Name:     "SystemDeptEdit",
							Status:   1,
							Type:     "button",
							AuthCode: "System:Dept:Edit",
							Meta: &v1.MenuMeta{
								Title: "common.edit",
							},
						},
						{
							Id:       20403,
							Pid:      202,
							Name:     "SystemDeptDelete",
							Status:   1,
							Type:     "button",
							AuthCode: "System:Dept:Delete",
							Meta: &v1.MenuMeta{
								Title: "common.delete",
							},
						},
					},
				},
			},
		},
		{
			Id:     9,
			Name:   "Project",
			Path:   "/vben-admin",
			Type:   "catalog",
			Status: 1,
			Meta: &v1.MenuMeta{
				BadgeType: "dot",
				Order:     9998,
				Title:     "demos.vben.title",
				Icon:      "carbon:data-center",
			},
			Children: []*v1.MenuItem{
				{
					Id:        901,
					Pid:       9,
					Name:      "VbenDocument",
					Path:      "/vben-admin/document",
					Component: "IFrameView",
					Type:      "embedded",
					Status:    1,
					Meta: &v1.MenuMeta{
						Icon:      "carbon:book",
						IframeSrc: "https://doc.vben.pro",
						Title:     "demos.vben.document",
					},
				},
				{
					Id:        902,
					Pid:       9,
					Name:      "VbenGithub",
					Path:      "/vben-admin/github",
					Component: "IFrameView",
					Type:      "link",
					Status:    1,
					Meta: &v1.MenuMeta{
						Icon:  "carbon:logo-github",
						Link:  "https://github.com/vbenjs/vue-vben-admin",
						Title: "Github",
					},
				},
				{
					Id:        903,
					Pid:       9,
					Name:      "VbenAntdv",
					Path:      "/vben-admin/antdv",
					Component: "IFrameView",
					Type:      "link",
					Status:    0,
					Meta: &v1.MenuMeta{
						Icon:      "carbon:hexagon-vertical-solid",
						BadgeType: "dot",
						Link:      "https://ant.vben.pro",
						Title:     "demos.vben.antdv",
					},
				},
			},
		},
		{
			Id:        10,
			Component: "_core/about/index",
			Type:      "menu",
			Status:    1,
			Name:      "About",
			Path:      "/about",
			Meta: &v1.MenuMeta{
				Icon:  "lucide:copyright",
				Order: 9999,
				Title: "demos.vben.about",
			},
		},
	}
}

func filterMenuRoutes(items []*v1.MenuItem) v1.MenuAllRes {
	filtered := make([]*v1.MenuItem, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		if item.Type == "button" {
			continue
		}
		if item.Path == "" {
			continue
		}
		if len(item.Children) > 0 {
			item.Children = filterMenuRoutes(item.Children)
		}
		filtered = append(filtered, item)
	}
	return filtered
}

type menuRecord struct {
	Id             string `json:"id" orm:"id"`
	TenantId       string `json:"tenantId" orm:"tenant_id"`
	ParentId       string `json:"parentId" orm:"parent_id"`
	Name           string `json:"name" orm:"name"`
	Path           string `json:"path" orm:"path"`
	Component      string `json:"component" orm:"component"`
	Icon           string `json:"icon" orm:"icon"`
	Order          int    `json:"order" orm:"order"`
	Type           string `json:"type" orm:"type"`
	Visible        int    `json:"visible" orm:"visible"`
	Status         int    `json:"status" orm:"status"`
	PermissionCode string `json:"permissionCode" orm:"permission_code"`
	Meta           string `json:"meta" orm:"meta"`
}

func fetchMenuFromDB(ctx context.Context) (v1.MenuAllRes, error) {
	tenantID := resolveTenantID(ctx)
	var records []menuRecord
	err := g.DB().Ctx(ctx).Model("sys_menu").
		Where("tenant_id", tenantID).
		Where("status", 1).
		Where("deleted_at is null").
		Order("\"order\" asc").
		Scan(&records)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}

	itemsByID := make(map[string]*v1.MenuItem, len(records))
	for _, record := range records {
		item := &v1.MenuItem{
			Id:        0,
			Pid:       0,
			Name:      record.Name,
			Path:      record.Path,
			Component: record.Component,
			Type:      record.Type,
			Status:    record.Status,
			Icon:      record.Icon,
			AuthCode:  record.PermissionCode,
		}

		if record.Meta != "" {
			var meta v1.MenuMeta
			if err := json.Unmarshal([]byte(record.Meta), &meta); err == nil {
				item.Meta = &meta
			}
		}
		if item.Meta == nil && record.Order != 0 {
			item.Meta = &v1.MenuMeta{}
		}
		if item.Meta != nil {
			item.Meta.Order = record.Order
		}

		itemsByID[record.Id] = item
	}

	var roots []*v1.MenuItem
	for _, record := range records {
		item := itemsByID[record.Id]
		if record.ParentId == "" {
			roots = append(roots, item)
			continue
		}
		parent := itemsByID[record.ParentId]
		if parent == nil {
			roots = append(roots, item)
			continue
		}
		parent.Children = append(parent.Children, item)
	}

	return roots, nil
}

func resolveTenantID(ctx context.Context) string {
	const defaultTenantID = "00000000-0000-0000-0000-000000000000"
	token, err := resolveAccessToken(ctx, "")
	if err != nil {
		return defaultTenantID
	}
	claims, err := parseToken(token)
	if err != nil {
		return defaultTenantID
	}
	tenantID, _ := claims["tenantId"].(string)
	if strings.TrimSpace(tenantID) == "" {
		return defaultTenantID
	}
	return tenantID
}
