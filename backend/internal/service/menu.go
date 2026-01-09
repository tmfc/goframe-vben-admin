package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	v1 "backend/api/menu/v1"
	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
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
	CreateMenu(ctx context.Context, in model.SysMenuCreateIn) (id string, err error)
	GetMenu(ctx context.Context, id string) (out *model.SysMenuGetOut, err error)
	UpdateMenu(ctx context.Context, in model.SysMenuUpdateIn) (err error)
	DeleteMenu(ctx context.Context, id string) (err error)
	GetMenuList(ctx context.Context, in model.SysMenuGetListIn) (out *model.SysMenuGetListOut, err error)
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

// CreateMenu creates a new menu.
func (s *sMenu) CreateMenu(ctx context.Context, in model.SysMenuCreateIn) (id string, err error) {
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return "", err
	}

	parentId := in.ParentId
	if parentId == "0" {
		parentId = ""
	}

	// 使用事务确保 menu 和 permission 的创建是原子操作
	err = dao.SysMenu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		insertData := buildMenuCreateData(ctx, in, parentId)
		_, err = tx.Model("sys_menu").Data(insertData).Insert()
		if err != nil {
			return err
		}

		var insertedMenu entity.SysMenu
		query := tx.Model("sys_menu").
			Where(dao.SysMenu.Columns().Name, in.Name).
			OrderDesc(dao.SysMenu.Columns().CreatedAt).
			Limit(1)
		if parentId != "" {
			query = query.Where(dao.SysMenu.Columns().ParentId, parentId)
		} else {
			query = query.WhereNull(dao.SysMenu.Columns().ParentId)
		}
		if err = query.Scan(&insertedMenu); err != nil {
			return err
		}

		// 如果提供了 permission_code，则自动创建对应的 Permission
		if in.PermissionCode != "" {
			_, err = tx.Model("sys_permission").Data(g.Map{
				dao.SysPermission.Columns().Name:        in.Name,
				dao.SysPermission.Columns().Description: fmt.Sprintf("Permission for menu: %s", in.Name),
				dao.SysPermission.Columns().Status:      1,
			}).Insert()
			if err != nil {
				return err
			}
		}

		id = insertedMenu.Id
		return nil
	})

	return id, err
}

// GetMenu retrieves a menu by ID.
func (s *sMenu) GetMenu(ctx context.Context, id string) (out *model.SysMenuGetOut, err error) {
	out = &model.SysMenuGetOut{}
	err = dao.SysMenu.Ctx(ctx).Where(dao.SysMenu.Columns().Id, id).Scan(&out.SysMenu)
	if err != nil {
		return nil, err
	}
	if out.SysMenu == nil {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Menu with ID %s not found", id)
	}
	return out, nil
}

// UpdateMenu updates an existing menu.
func (s *sMenu) UpdateMenu(ctx context.Context, in model.SysMenuUpdateIn) (err error) {
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}

	// 使用事务确保 menu 和 permission 的更新是原子操作
	err = dao.SysMenu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 获取原始 menu 信息以检查 permission_code 是否变化
		var originalMenu entity.SysMenu
		err = tx.Model("sys_menu").Where(dao.SysMenu.Columns().Id, in.ID).Scan(&originalMenu)
		if err != nil {
			return err
		}

		parentId := in.ParentId
		if parentId == "0" {
			parentId = ""
		}

		_, err = tx.Model("sys_menu").Data(buildMenuUpdateData(in, parentId)).
			Where(dao.SysMenu.Columns().Id, in.ID).Update()
		if err != nil {
			return err
		}

		// 如果 permission_code 发生变化，则同步更新 Permission
		if originalMenu.PermissionCode != in.PermissionCode {
			// 如果旧的 permission_code 存在，删除对应的 Permission
			if originalMenu.PermissionCode != "" {
				// 查找对应的 permission 记录
				var oldPerms []*entity.SysPermission
				err = tx.Model("sys_permission").
					Where(dao.SysPermission.Columns().Name, originalMenu.Name).
					Scan(&oldPerms)
				if err == nil && len(oldPerms) > 0 {
					for _, perm := range oldPerms {
						_, err = tx.Model("sys_permission").
							Where(dao.SysPermission.Columns().Id, perm.Id).
							Delete()
						if err != nil {
							return err
						}
					}
				}
			}

			// 如果新的 permission_code 存在，创建对应的 Permission
			if in.PermissionCode != "" {
				_, err = tx.Model("sys_permission").Data(g.Map{
					dao.SysPermission.Columns().Name:        in.Name,
					dao.SysPermission.Columns().Description: fmt.Sprintf("Permission for menu: %s", in.Name),
					dao.SysPermission.Columns().Status:      1,
				}).Insert()
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

// DeleteMenu deletes a menu by ID.
func (s *sMenu) DeleteMenu(ctx context.Context, id string) (err error) {
	// 使用事务确保 menu 和 permission 的删除是原子操作
	err = dao.SysMenu.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 获取要删除的 menu 信息
		var menuToDelete entity.SysMenu
		err = tx.Model("sys_menu").Where(dao.SysMenu.Columns().Id, id).Scan(&menuToDelete)
		if err != nil {
			return err
		}

		// 删除 menu
		_, err = tx.Model("sys_menu").Where(dao.SysMenu.Columns().Id, id).Delete()
		if err != nil {
			return err
		}

		// 如果 menu 有 permission_code，则级联删除对应的 Permission
		if menuToDelete.PermissionCode != "" {
			// 查找对应的 permission 记录
			var perms []*entity.SysPermission
			err = tx.Model("sys_permission").
				Where(dao.SysPermission.Columns().Name, menuToDelete.Name).
				Scan(&perms)
			if err == nil && len(perms) > 0 {
				for _, perm := range perms {
					_, err = tx.Model("sys_permission").
						Where(dao.SysPermission.Columns().Id, perm.Id).
						Delete()
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	return err
}

// GetMenuList retrieves menu list.
func (s *sMenu) GetMenuList(ctx context.Context, in model.SysMenuGetListIn) (out *model.SysMenuGetListOut, err error) {
	out = &model.SysMenuGetListOut{}
	query := dao.SysMenu.Ctx(ctx).OmitEmpty()

	if in.Name != "" {
		// Use case-insensitive matching for PostgreSQL so that queries like
		// "manage" can match names such as "Manage".
		query = query.Where(dao.SysMenu.Columns().Name+" ILIKE ?", "%"+in.Name+"%")
	}
	if in.Status != "" {
		query = query.Where(dao.SysMenu.Columns().Status, in.Status)
	}

	out.Total, err = query.Count()
	if err != nil {
		return nil, err
	}

	err = query.Scan(&out.List)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func buildMenuCreateData(ctx context.Context, in model.SysMenuCreateIn, parentId string) g.Map {
	columns := dao.SysMenu.Columns()
	data := g.Map{
		columns.Name:      in.Name,
		columns.Path:      in.Path,
		columns.Component: in.Component,
		columns.Icon:      in.Icon,
		columns.Type:      in.Type,
		columns.Status:    in.Status,
		columns.Order:     in.Order,
		columns.TenantId:  resolveTenantID(ctx),
	}
	if in.Meta != "" {
		data[columns.Meta] = in.Meta
	}
	if in.PermissionCode != "" {
		data[columns.PermissionCode] = in.PermissionCode
	}
	if parentId != "" {
		data[columns.ParentId] = parentId
	}
	return data
}

func buildMenuUpdateData(in model.SysMenuUpdateIn, parentId string) g.Map {
	columns := dao.SysMenu.Columns()
	data := g.Map{
		columns.Name:      in.Name,
		columns.Path:      in.Path,
		columns.Component: in.Component,
		columns.Icon:      in.Icon,
		columns.Type:      in.Type,
		columns.Status:    in.Status,
		columns.Order:     in.Order,
	}
	if in.Meta != "" {
		data[columns.Meta] = in.Meta
	}
	if in.PermissionCode != "" {
		data[columns.PermissionCode] = in.PermissionCode
	}
	if parentId != "" {
		data[columns.ParentId] = parentId
	}
	return data
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

	var records []menuRecord
	err := dao.SysMenu.Ctx(ctx).
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
		if item.Meta == nil {
			item.Meta = &v1.MenuMeta{}
		}
		item.Meta.Order = record.Order

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

	return sortMenuItems(roots), nil
}

func sortMenuItems(items []*v1.MenuItem) []*v1.MenuItem {
	if len(items) == 0 {
		return items
	}
	sort.SliceStable(items, func(i, j int) bool {
		return menuOrder(items[i]) < menuOrder(items[j])
	})
	for _, item := range items {
		item.Children = sortMenuItems(item.Children)
	}
	return items
}

func menuOrder(item *v1.MenuItem) int {
	if item == nil || item.Meta == nil {
		return 9999
	}
	return item.Meta.Order
}
