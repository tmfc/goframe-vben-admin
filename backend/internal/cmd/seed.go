package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"backend/internal/consts"
	"backend/internal/dao"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"golang.org/x/crypto/bcrypt"
)

const defaultSeedFile = "resource/seed/initial.json"

var Seed = gcmd.Command{
	Name:  "seed",
	Usage: "seed [OPTIONS]",
	Brief: "import initial seed data (idempotent)",
	Arguments: []gcmd.Argument{
		{Name: "file", Short: "f", Brief: "seed file path (JSON)"},
		{Name: "tenant", Short: "t", Brief: "default tenant id"},
		{Name: "dry-run", Short: "n", Brief: "print actions without inserting", Orphan: true},
	},
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		filePath := strings.TrimSpace(parser.GetOpt("file").String())
		tenantID := strings.TrimSpace(parser.GetOpt("tenant").String())
		if tenantID == "" {
			tenantID = consts.DefaultTenantID
		}
		dryRun := parser.GetOpt("dry-run").Bool()

		seedPath, seedData, err := loadSeedData(filePath)
		if err != nil {
			return err
		}

		seedCtx := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)
		result, err := seedInitialData(seedCtx, seedData, tenantID, dryRun)
		if err != nil {
			return err
		}

		fmt.Printf("Seed file: %s\n", seedPath)
		fmt.Printf("Tenants: inserted=%d skipped=%d\n", result.Tenants.Inserted, result.Tenants.Skipped)
		fmt.Printf("Depts: inserted=%d skipped=%d\n", result.Depts.Inserted, result.Depts.Skipped)
		fmt.Printf("Roles: inserted=%d skipped=%d\n", result.Roles.Inserted, result.Roles.Skipped)
		fmt.Printf("Permissions: inserted=%d skipped=%d\n", result.Permissions.Inserted, result.Permissions.Skipped)
		fmt.Printf("Menus: inserted=%d skipped=%d\n", result.Menus.Inserted, result.Menus.Skipped)
		fmt.Printf("Users: inserted=%d skipped=%d\n", result.Users.Inserted, result.Users.Skipped)
		fmt.Printf("Role-Permissions: inserted=%d skipped=%d\n", result.RolePermissions.Inserted, result.RolePermissions.Skipped)
		fmt.Printf("User-Roles: inserted=%d skipped=%d\n", result.UserRoles.Inserted, result.UserRoles.Skipped)
		fmt.Printf("Casbin Rules: inserted=%d skipped=%d\n", result.CasbinRules.Inserted, result.CasbinRules.Skipped)
		if dryRun {
			fmt.Println("Dry run enabled: no data was written.")
		}
		return nil
	},
}

func init() {
	if err := Main.AddCommand(&Seed); err != nil {
		panic(err)
	}
}

type SeedData struct {
	Tenants         []SeedTenant         `json:"tenants"`
	Users           []SeedUser           `json:"users"`
	Roles           []SeedRole           `json:"roles"`
	Depts           []SeedDept           `json:"depts"`
	Permissions     []SeedPermission     `json:"permissions"`
	Menus           []SeedMenu           `json:"menus"`
	UserRoles       []SeedUserRole       `json:"user_roles"`
	RolePermissions []SeedRolePermission `json:"role_permissions"`
	CasbinRules     []SeedCasbinRule     `json:"casbin_rules"`
}

type SeedTenant struct {
	ID     *int64 `json:"id,omitempty"`
	Name   string `json:"name"`
	Status *int   `json:"status,omitempty"`
}

type SeedUser struct {
	ID           *int64   `json:"id,omitempty"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	PasswordHash string   `json:"password_hash"`
	RealName     string   `json:"real_name"`
	Avatar       string   `json:"avatar"`
	HomePath     string   `json:"home_path"`
	Status       *int     `json:"status,omitempty"`
	Roles        []string `json:"roles"`
	DeptName     string   `json:"dept_name"`
	DeptID       *int64   `json:"dept_id,omitempty"`
	TenantID     *int64   `json:"tenant_id,omitempty"`
}

type SeedMenu struct {
	ID             *int64          `json:"id,omitempty"`
	TenantID       *int64          `json:"tenant_id,omitempty"`
	ParentID       *int64          `json:"parent_id,omitempty"`
	ParentName     string          `json:"parent_name"`
	Name           string          `json:"name"`
	Path           string          `json:"path"`
	Component      *string         `json:"component"`
	Icon           *string         `json:"icon"`
	Order          *int            `json:"order,omitempty"`
	Type           string          `json:"type"`
	Visible        *int            `json:"visible,omitempty"`
	Status         *int            `json:"status,omitempty"`
	PermissionCode *string         `json:"permission_code"`
	Meta           json.RawMessage `json:"meta"`
}

type SeedRole struct {
	ID          *int64 `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *int64 `json:"parent_id,omitempty"`
	ParentName  string `json:"parent_name"`
	Status      *int   `json:"status,omitempty"`
	TenantID    *int64 `json:"tenant_id,omitempty"`
}

type SeedPermission struct {
	ID          *int64 `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *int64 `json:"parent_id,omitempty"`
	ParentName  string `json:"parent_name"`
	Status      *int   `json:"status,omitempty"`
	TenantID    *int64 `json:"tenant_id,omitempty"`
}

type SeedDept struct {
	ID         *int64 `json:"id,omitempty"`
	Name       string `json:"name"`
	ParentID   *int64 `json:"parent_id,omitempty"`
	ParentName string `json:"parent_name"`
	Order      *int   `json:"order,omitempty"`
	Status     *int   `json:"status,omitempty"`
	TenantID   *int64 `json:"tenant_id,omitempty"`
}

type SeedUserRole struct {
	UserID   *int64 `json:"user_id,omitempty"`
	Username string `json:"username"`
	RoleID   *int64 `json:"role_id,omitempty"`
	RoleName string `json:"role"`
	TenantID *int64 `json:"tenant_id,omitempty"`
}

type SeedRolePermission struct {
	RoleID        *int64 `json:"role_id,omitempty"`
	RoleName      string `json:"role"`
	PermissionID  *int64 `json:"permission_id,omitempty"`
	PermissionKey string `json:"permission"`
	TenantID      *int64 `json:"tenant_id,omitempty"`
}

type SeedCasbinRule struct {
	Ptype string `json:"ptype"`
	V0    string `json:"v0"`
	V1    string `json:"v1"`
	V2    string `json:"v2"`
	V3    string `json:"v3"`
	V4    string `json:"v4"`
	V5    string `json:"v5"`
}

type seedStats struct {
	Inserted int
	Skipped  int
}

type seedResult struct {
	Tenants         seedStats
	Depts           seedStats
	Roles           seedStats
	Permissions     seedStats
	Menus           seedStats
	Users           seedStats
	UserRoles       seedStats
	RolePermissions seedStats
	CasbinRules     seedStats
}

func loadSeedData(path string) (string, *SeedData, error) {
	if strings.TrimSpace(path) == "" {
		path = defaultSeedFile
	}
	if !gfile.Exists(path) {
		if searched, _ := gfile.Search(path); searched != "" {
			path = searched
		}
	}
	if !gfile.Exists(path) {
		return "", nil, fmt.Errorf("seed file not found: %s", path)
	}
	content := gfile.GetContents(path)
	if strings.TrimSpace(content) == "" {
		return path, &SeedData{}, nil
	}
	var data SeedData
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return "", nil, fmt.Errorf("failed to parse seed file: %w", err)
	}
	return path, &data, nil
}

func seedInitialData(ctx context.Context, data *SeedData, defaultTenant string, dryRun bool) (*seedResult, error) {
	result := &seedResult{}
	roleIDs := make(map[string]uint)
	permIDs := make(map[string]uint)
	userIDs := make(map[string]int64)
	deptIDs := make(map[string]int64)
	menuIDs := make(map[string]int64)

	if err := seedTenants(ctx, data.Tenants, dryRun, &result.Tenants); err != nil {
		return nil, err
	}
	if err := seedDepts(ctx, data.Depts, defaultTenant, dryRun, deptIDs, &result.Depts); err != nil {
		return nil, err
	}
	if err := seedRoles(ctx, data.Roles, defaultTenant, dryRun, roleIDs, &result.Roles); err != nil {
		return nil, err
	}
	if err := seedPermissions(ctx, data.Permissions, defaultTenant, dryRun, permIDs, &result.Permissions); err != nil {
		return nil, err
	}
	if err := seedMenus(ctx, data.Menus, defaultTenant, dryRun, menuIDs, &result.Menus); err != nil {
		return nil, err
	}
	if err := seedUsers(ctx, data.Users, defaultTenant, dryRun, deptIDs, userIDs, &result.Users); err != nil {
		return nil, err
	}
	if err := seedRolePermissions(ctx, data.RolePermissions, defaultTenant, dryRun, roleIDs, permIDs, &result.RolePermissions); err != nil {
		return nil, err
	}
	if err := seedUserRoles(ctx, data.UserRoles, defaultTenant, dryRun, userIDs, roleIDs, &result.UserRoles); err != nil {
		return nil, err
	}
	if err := seedCasbinRules(ctx, data.CasbinRules, dryRun, &result.CasbinRules); err != nil {
		return nil, err
	}
	return result, nil
}

func seedTenants(ctx context.Context, tenants []SeedTenant, dryRun bool, stats *seedStats) error {
	var fakeID int64
	for _, tenant := range tenants {
		name := strings.TrimSpace(tenant.Name)
		if name == "" {
			return fmt.Errorf("tenant name is required")
		}
		query := dao.SysTenant.Ctx(ctx)
		if tenant.ID != nil && *tenant.ID > 0 {
			query = query.Where(dao.SysTenant.Columns().Id, *tenant.ID)
		} else {
			query = query.Where(dao.SysTenant.Columns().Name, name)
		}
		count, err := query.Count()
		if err != nil {
			return err
		}
		if count > 0 {
			stats.Skipped++
			continue
		}

		insertData := map[string]any{
			dao.SysTenant.Columns().Name:   name,
			dao.SysTenant.Columns().Status: intValueOrDefault(tenant.Status, 1),
		}
		if tenant.ID != nil && *tenant.ID > 0 {
			insertData[dao.SysTenant.Columns().Id] = *tenant.ID
		}

		if dryRun {
			fakeID++
			stats.Inserted++
			continue
		}
		_, err = dao.SysTenant.Ctx(ctx).Data(insertData).Insert()
		if err != nil {
			return err
		}
		stats.Inserted++
	}
	return nil
}

func seedDepts(ctx context.Context, depts []SeedDept, defaultTenant string, dryRun bool, deptIDs map[string]int64, stats *seedStats) error {
	var fakeID int64
	for _, dept := range depts {
		name := strings.TrimSpace(dept.Name)
		if name == "" {
			return fmt.Errorf("dept name is required")
		}
		tenantID := resolveTenantID(defaultTenant, dept.TenantID)
		tenantIDValue := tenantIDToInt64(tenantID)
		if tenantIDValue == 0 {
			return fmt.Errorf("invalid tenant id for dept %s: %s", name, tenantID)
		}
		tenantCtx := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)
		parentID := resolveParentDeptID(tenantCtx, dept.ParentID, dept.ParentName, deptIDs)
		if parentID == -1 {
			return fmt.Errorf("dept parent not found: %s", dept.ParentName)
		}

		query := dao.SysDept.Ctx(tenantCtx).Where(dao.SysDept.Columns().Name, name)
		if parentID > 0 {
			query = query.Where(dao.SysDept.Columns().ParentId, parentID)
		} else {
			query = query.WhereNull(dao.SysDept.Columns().ParentId)
		}
		var existing struct {
			Id int64
		}
		if err := query.Fields(dao.SysDept.Columns().Id).Scan(&existing); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
		if existing.Id > 0 {
			stats.Skipped++
			deptIDs[mapKey(tenantID, name)] = existing.Id
			continue
		}

		insertData := map[string]any{
			dao.SysDept.Columns().Name:     name,
			dao.SysDept.Columns().Order:    intValueOrDefault(dept.Order, 0),
			dao.SysDept.Columns().Status:   intValueOrDefault(dept.Status, 1),
			dao.SysDept.Columns().TenantId: tenantIDValue,
		}
		if parentID > 0 {
			insertData[dao.SysDept.Columns().ParentId] = parentID
		}
		if dept.ID != nil && *dept.ID > 0 {
			insertData[dao.SysDept.Columns().Id] = *dept.ID
		}

		if dryRun {
			fakeID++
			stats.Inserted++
			deptIDs[mapKey(tenantID, name)] = fakeID
			continue
		}
		result, err := dao.SysDept.CtxNoTenant(tenantCtx).Data(insertData).Insert()
		if err != nil {
			return err
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		stats.Inserted++
		deptIDs[mapKey(tenantID, name)] = lastID
	}
	return nil
}

func seedMenus(ctx context.Context, menus []SeedMenu, defaultTenant string, dryRun bool, menuIDs map[string]int64, stats *seedStats) error {
	var fakeID int64
	for _, menu := range menus {
		name := strings.TrimSpace(menu.Name)
		path := strings.TrimSpace(menu.Path)
		menuType := strings.TrimSpace(menu.Type)
		if name == "" || path == "" || menuType == "" {
			return fmt.Errorf("menu name, path, and type are required")
		}
		tenantID := resolveTenantID(defaultTenant, menu.TenantID)
		tenantIDValue := tenantIDToInt64(tenantID)
		if tenantIDValue == 0 {
			return fmt.Errorf("invalid tenant id for menu %s: %s", name, tenantID)
		}
		tenantCtx := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)
		parentID := resolveParentMenuID(tenantCtx, menu.ParentID, menu.ParentName, menuIDs)
		if parentID == -1 {
			return fmt.Errorf("menu parent not found: %s", menu.ParentName)
		}

		query := dao.SysMenu.Ctx(tenantCtx)
		if menu.ID != nil && *menu.ID > 0 {
			query = query.Where(dao.SysMenu.Columns().Id, *menu.ID)
		} else {
			query = query.Where(dao.SysMenu.Columns().Path, path)
		}
		var existing struct {
			Id int64
		}
		if err := query.Fields(dao.SysMenu.Columns().Id).Scan(&existing); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
		if existing.Id > 0 {
			stats.Skipped++
			menuIDs[mapKey(tenantID, name)] = existing.Id
			continue
		}

		insertData := map[string]any{
			dao.SysMenu.Columns().Name:     name,
			dao.SysMenu.Columns().Path:     path,
			dao.SysMenu.Columns().Type:     menuType,
			dao.SysMenu.Columns().Order:    intValueOrDefault(menu.Order, 0),
			dao.SysMenu.Columns().Visible:  intValueOrDefault(menu.Visible, 1),
			dao.SysMenu.Columns().Status:   intValueOrDefault(menu.Status, 1),
			dao.SysMenu.Columns().TenantId: tenantIDValue,
		}
		if parentID > 0 {
			insertData[dao.SysMenu.Columns().ParentId] = parentID
		}
		if menu.Component != nil {
			insertData[dao.SysMenu.Columns().Component] = strings.TrimSpace(*menu.Component)
		}
		if menu.Icon != nil {
			insertData[dao.SysMenu.Columns().Icon] = strings.TrimSpace(*menu.Icon)
		}
		if menu.PermissionCode != nil {
			insertData[dao.SysMenu.Columns().PermissionCode] = strings.TrimSpace(*menu.PermissionCode)
		}
		if len(menu.Meta) > 0 {
			insertData[dao.SysMenu.Columns().Meta] = string(menu.Meta)
		}
		if menu.ID != nil && *menu.ID > 0 {
			insertData[dao.SysMenu.Columns().Id] = *menu.ID
		}

		if dryRun {
			fakeID++
			stats.Inserted++
			menuIDs[mapKey(tenantID, name)] = fakeID
			continue
		}
		result, err := dao.SysMenu.CtxNoTenant(tenantCtx).Data(insertData).Insert()
		if err != nil {
			return err
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		stats.Inserted++
		menuIDs[mapKey(tenantID, name)] = lastID
	}
	return nil
}

func seedRoles(ctx context.Context, roles []SeedRole, defaultTenant string, dryRun bool, roleIDs map[string]uint, stats *seedStats) error {
	var fakeID uint
	for _, role := range roles {
		name := strings.TrimSpace(role.Name)
		if name == "" {
			return fmt.Errorf("role name is required")
		}
		tenantID := resolveTenantID(defaultTenant, role.TenantID)
		tenantIDValue := tenantIDToInt64(tenantID)
		if tenantIDValue == 0 {
			return fmt.Errorf("invalid tenant id for role %s: %s", name, tenantID)
		}
		tenantCtx := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)
		parentID := resolveParentRoleID(tenantCtx, role.ParentID, role.ParentName, roleIDs)
		if parentID == -1 {
			return fmt.Errorf("role parent not found: %s", role.ParentName)
		}

		query := dao.SysRole.Ctx(tenantCtx).Where(dao.SysRole.Columns().Name, name)
		var existing struct {
			Id int64
		}
		if err := query.Fields(dao.SysRole.Columns().Id).Scan(&existing); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
		if existing.Id > 0 {
			stats.Skipped++
			roleIDs[mapKey(tenantID, name)] = uint(existing.Id)
			continue
		}

		insertData := map[string]any{
			dao.SysRole.Columns().Name:        name,
			dao.SysRole.Columns().Description: strings.TrimSpace(role.Description),
			dao.SysRole.Columns().Status:      intValueOrDefault(role.Status, 1),
			dao.SysRole.Columns().TenantId:    tenantIDValue,
		}
		if parentID > 0 {
			insertData[dao.SysRole.Columns().ParentId] = parentID
		}
		if role.ID != nil && *role.ID > 0 {
			insertData[dao.SysRole.Columns().Id] = *role.ID
		}

		if dryRun {
			fakeID++
			stats.Inserted++
			roleIDs[mapKey(tenantID, name)] = fakeID
			continue
		}
		result, err := dao.SysRole.CtxNoTenant(tenantCtx).Data(insertData).Insert()
		if err != nil {
			return err
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		stats.Inserted++
		roleIDs[mapKey(tenantID, name)] = uint(lastID)
	}
	return nil
}

func seedPermissions(ctx context.Context, permissions []SeedPermission, defaultTenant string, dryRun bool, permIDs map[string]uint, stats *seedStats) error {
	var fakeID uint
	for _, perm := range permissions {
		name := strings.TrimSpace(perm.Name)
		if name == "" {
			return fmt.Errorf("permission name is required")
		}
		tenantID := resolveTenantID(defaultTenant, perm.TenantID)
		tenantIDValue := tenantIDToInt64(tenantID)
		if tenantIDValue == 0 {
			return fmt.Errorf("invalid tenant id for permission %s: %s", name, tenantID)
		}
		tenantCtx := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)
		parentID := resolveParentPermissionID(tenantCtx, perm.ParentID, perm.ParentName, permIDs)
		if parentID == -1 {
			return fmt.Errorf("permission parent not found: %s", perm.ParentName)
		}

		query := dao.SysPermission.Ctx(tenantCtx).Where(dao.SysPermission.Columns().Name, name)
		var existing struct {
			Id int64
		}
		if err := query.Fields(dao.SysPermission.Columns().Id).Scan(&existing); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
		if existing.Id > 0 {
			stats.Skipped++
			permIDs[mapKey(tenantID, name)] = uint(existing.Id)
			continue
		}

		insertData := map[string]any{
			dao.SysPermission.Columns().Name:        name,
			dao.SysPermission.Columns().Description: strings.TrimSpace(perm.Description),
			dao.SysPermission.Columns().Status:      intValueOrDefault(perm.Status, 1),
			dao.SysPermission.Columns().TenantId:    tenantIDValue,
		}
		if parentID > 0 {
			insertData[dao.SysPermission.Columns().ParentId] = parentID
		}
		if perm.ID != nil && *perm.ID > 0 {
			insertData[dao.SysPermission.Columns().Id] = *perm.ID
		}

		if dryRun {
			fakeID++
			stats.Inserted++
			permIDs[mapKey(tenantID, name)] = fakeID
			continue
		}
		result, err := dao.SysPermission.CtxNoTenant(tenantCtx).Data(insertData).Insert()
		if err != nil {
			return err
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		stats.Inserted++
		permIDs[mapKey(tenantID, name)] = uint(lastID)
	}
	return nil
}

func seedUsers(ctx context.Context, users []SeedUser, defaultTenant string, dryRun bool, deptIDs map[string]int64, userIDs map[string]int64, stats *seedStats) error {
	var fakeID int64
	for _, user := range users {
		username := strings.TrimSpace(user.Username)
		if username == "" {
			return fmt.Errorf("username is required")
		}
		tenantID := resolveTenantID(defaultTenant, user.TenantID)
		tenantIDValue := tenantIDToInt64(tenantID)
		if tenantIDValue == 0 {
			return fmt.Errorf("invalid tenant id for user %s: %s", username, tenantID)
		}
		tenantCtx := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)

		query := dao.SysUser.Ctx(tenantCtx).Where(dao.SysUser.Columns().Username, username)
		var existing struct {
			Id int64
		}
		if err := query.Fields(dao.SysUser.Columns().Id).Scan(&existing); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
		if existing.Id > 0 {
			stats.Skipped++
			userIDs[mapKey(tenantID, username)] = existing.Id
			continue
		}

		passwordHash := strings.TrimSpace(user.PasswordHash)
		if passwordHash == "" {
			password := user.Password
			if strings.TrimSpace(password) == "" {
				return fmt.Errorf("password or password_hash is required for user %s", username)
			}
			hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			passwordHash = string(hashed)
		}

		var deptID int64
		if user.DeptID != nil && *user.DeptID > 0 {
			deptID = *user.DeptID
		} else if strings.TrimSpace(user.DeptName) != "" {
			deptID = resolveDeptIDByName(tenantCtx, strings.TrimSpace(user.DeptName), deptIDs)
			if deptID == 0 {
				return fmt.Errorf("dept not found for user %s: %s", username, user.DeptName)
			}
		}

		rolesJSON := "[]"
		if len(user.Roles) > 0 {
			if b, err := json.Marshal(user.Roles); err == nil {
				rolesJSON = string(b)
			}
		}

		insertData := map[string]any{
			dao.SysUser.Columns().Username:  username,
			dao.SysUser.Columns().Password:  passwordHash,
			dao.SysUser.Columns().RealName:  strings.TrimSpace(user.RealName),
			dao.SysUser.Columns().Avatar:    strings.TrimSpace(user.Avatar),
			dao.SysUser.Columns().HomePath:  strings.TrimSpace(user.HomePath),
			dao.SysUser.Columns().Status:    intValueOrDefault(user.Status, 1),
			dao.SysUser.Columns().Roles:     rolesJSON,
			dao.SysUser.Columns().ExtInfo:   nil,
			dao.SysUser.Columns().DeletedAt: nil,
			dao.SysUser.Columns().TenantId:  tenantIDValue,
		}
		if deptID > 0 {
			insertData[dao.SysUser.Columns().DeptId] = deptID
		}
		if user.ID != nil && *user.ID > 0 {
			insertData[dao.SysUser.Columns().Id] = *user.ID
		}

		if dryRun {
			fakeID++
			stats.Inserted++
			userIDs[mapKey(tenantID, username)] = fakeID
			continue
		}
		result, err := dao.SysUser.CtxNoTenant(tenantCtx).Data(insertData).Insert()
		if err != nil {
			return err
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		stats.Inserted++
		userIDs[mapKey(tenantID, username)] = lastID
	}
	return nil
}

func seedRolePermissions(ctx context.Context, items []SeedRolePermission, defaultTenant string, dryRun bool, roleIDs map[string]uint, permIDs map[string]uint, stats *seedStats) error {
	for _, item := range items {
		tenantID := resolveTenantID(defaultTenant, item.TenantID)
		tenantIDValue := tenantIDToInt64(tenantID)
		if tenantIDValue == 0 {
			return fmt.Errorf("invalid tenant id for role permission: %s", tenantID)
		}
		tenantCtx := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)
		roleID, err := resolveRoleID(tenantCtx, tenantID, item.RoleID, item.RoleName, roleIDs)
		if err != nil {
			return err
		}
		permID, err := resolvePermissionID(tenantCtx, tenantID, item.PermissionID, item.PermissionKey, permIDs)
		if err != nil {
			return err
		}

		query := dao.SysRolePermission.Ctx(tenantCtx).
			Where(dao.SysRolePermission.Columns().RoleId, roleID).
			Where(dao.SysRolePermission.Columns().PermissionId, permID)
		count, err := query.Count()
		if err != nil {
			return err
		}
		if count > 0 {
			stats.Skipped++
			continue
		}

		if dryRun {
			stats.Inserted++
			continue
		}
		_, err = dao.SysRolePermission.CtxNoTenant(tenantCtx).Data(map[string]any{
			dao.SysRolePermission.Columns().RoleId:       roleID,
			dao.SysRolePermission.Columns().PermissionId: permID,
			dao.SysRolePermission.Columns().TenantId:     tenantIDValue,
		}).Insert()
		if err != nil {
			return err
		}
		stats.Inserted++
	}
	return nil
}

func seedUserRoles(ctx context.Context, items []SeedUserRole, defaultTenant string, dryRun bool, userIDs map[string]int64, roleIDs map[string]uint, stats *seedStats) error {
	for _, item := range items {
		tenantID := resolveTenantID(defaultTenant, item.TenantID)
		tenantIDValue := tenantIDToInt64(tenantID)
		if tenantIDValue == 0 {
			return fmt.Errorf("invalid tenant id for user role: %s", tenantID)
		}
		tenantCtx := context.WithValue(ctx, consts.CtxKeyTenantID, tenantID)
		userID, err := resolveUserID(tenantCtx, tenantID, item.UserID, item.Username, userIDs)
		if err != nil {
			return err
		}
		roleID, err := resolveRoleID(tenantCtx, tenantID, item.RoleID, item.RoleName, roleIDs)
		if err != nil {
			return err
		}

		query := dao.SysUserRole.Ctx(tenantCtx).
			Where(dao.SysUserRole.Columns().UserId, userID).
			Where(dao.SysUserRole.Columns().RoleId, roleID)
		count, err := query.Count()
		if err != nil {
			return err
		}
		if count > 0 {
			stats.Skipped++
			continue
		}

		if dryRun {
			stats.Inserted++
			continue
		}
		_, err = dao.SysUserRole.CtxNoTenant(tenantCtx).Data(map[string]any{
			dao.SysUserRole.Columns().UserId:   userID,
			dao.SysUserRole.Columns().RoleId:   roleID,
			dao.SysUserRole.Columns().TenantId: tenantIDValue,
		}).Insert()
		if err != nil {
			return err
		}
		stats.Inserted++
	}
	return nil
}

func seedCasbinRules(ctx context.Context, rules []SeedCasbinRule, dryRun bool, stats *seedStats) error {
	for _, rule := range rules {
		ptype := strings.TrimSpace(rule.Ptype)
		if ptype == "" {
			return fmt.Errorf("casbin rule ptype is required")
		}
		query := dao.CasbinRule.Ctx(ctx).
			Where(dao.CasbinRule.Columns().Ptype, ptype).
			Where(dao.CasbinRule.Columns().V0, strings.TrimSpace(rule.V0)).
			Where(dao.CasbinRule.Columns().V1, strings.TrimSpace(rule.V1)).
			Where(dao.CasbinRule.Columns().V2, strings.TrimSpace(rule.V2)).
			Where(dao.CasbinRule.Columns().V3, strings.TrimSpace(rule.V3)).
			Where(dao.CasbinRule.Columns().V4, strings.TrimSpace(rule.V4)).
			Where(dao.CasbinRule.Columns().V5, strings.TrimSpace(rule.V5))
		count, err := query.Count()
		if err != nil {
			return err
		}
		if count > 0 {
			stats.Skipped++
			continue
		}

		if dryRun {
			stats.Inserted++
			continue
		}
		_, err = dao.CasbinRule.Ctx(ctx).Data(map[string]any{
			dao.CasbinRule.Columns().Ptype: ptype,
			dao.CasbinRule.Columns().V0:    strings.TrimSpace(rule.V0),
			dao.CasbinRule.Columns().V1:    strings.TrimSpace(rule.V1),
			dao.CasbinRule.Columns().V2:    strings.TrimSpace(rule.V2),
			dao.CasbinRule.Columns().V3:    strings.TrimSpace(rule.V3),
			dao.CasbinRule.Columns().V4:    strings.TrimSpace(rule.V4),
			dao.CasbinRule.Columns().V5:    strings.TrimSpace(rule.V5),
		}).Insert()
		if err != nil {
			return err
		}
		stats.Inserted++
	}
	return nil
}

func resolveTenantID(defaultTenant string, tenantID *int64) string {
	if tenantID == nil || *tenantID <= 0 {
		return defaultTenant
	}
	return strconv.FormatInt(*tenantID, 10)
}

func tenantIDToInt64(tenantID string) int64 {
	id, err := strconv.ParseInt(strings.TrimSpace(tenantID), 10, 64)
	if err != nil || id <= 0 {
		return 0
	}
	return id
}

func mapKey(tenantID, name string) string {
	return tenantID + ":" + strings.ToLower(strings.TrimSpace(name))
}

func intValueOrDefault(value *int, defaultValue int) int {
	if value != nil {
		return *value
	}
	return defaultValue
}

func resolveParentDeptID(ctx context.Context, parentID *int64, parentName string, deptIDs map[string]int64) int64 {
	if parentID != nil && *parentID > 0 {
		return *parentID
	}
	if strings.TrimSpace(parentName) == "" {
		return 0
	}
	deptID := resolveDeptIDByName(ctx, strings.TrimSpace(parentName), deptIDs)
	if deptID == 0 {
		return -1
	}
	return deptID
}

func resolveParentMenuID(ctx context.Context, parentID *int64, parentName string, menuIDs map[string]int64) int64 {
	if parentID != nil && *parentID > 0 {
		return *parentID
	}
	if strings.TrimSpace(parentName) == "" {
		return 0
	}
	menuID := resolveMenuIDByName(ctx, strings.TrimSpace(parentName), menuIDs)
	if menuID == 0 {
		return -1
	}
	return menuID
}

func resolveParentRoleID(ctx context.Context, parentID *int64, parentName string, roleIDs map[string]uint) int64 {
	if parentID != nil && *parentID > 0 {
		return *parentID
	}
	if strings.TrimSpace(parentName) == "" {
		return 0
	}
	roleID := resolveRoleIDByName(ctx, strings.TrimSpace(parentName), roleIDs)
	if roleID == 0 {
		return -1
	}
	return int64(roleID)
}

func resolveParentPermissionID(ctx context.Context, parentID *int64, parentName string, permIDs map[string]uint) int64 {
	if parentID != nil && *parentID > 0 {
		return *parentID
	}
	if strings.TrimSpace(parentName) == "" {
		return 0
	}
	permID := resolvePermissionIDByName(ctx, strings.TrimSpace(parentName), permIDs)
	if permID == 0 {
		return -1
	}
	return int64(permID)
}

func resolveDeptIDByName(ctx context.Context, name string, deptIDs map[string]int64) int64 {
	tenantID := tenantIDFromCtx(ctx)
	key := mapKey(tenantID, name)
	if id, ok := deptIDs[key]; ok && id > 0 {
		return id
	}
	var existing struct {
		Id int64
	}
	err := dao.SysDept.Ctx(ctx).
		Where(dao.SysDept.Columns().Name, name).
		Fields(dao.SysDept.Columns().Id).
		Scan(&existing)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0
		}
		return 0
	}
	if existing.Id == 0 {
		return 0
	}
	deptIDs[key] = existing.Id
	return existing.Id
}

func resolveMenuIDByName(ctx context.Context, name string, menuIDs map[string]int64) int64 {
	tenantID := tenantIDFromCtx(ctx)
	key := mapKey(tenantID, name)
	if id, ok := menuIDs[key]; ok && id > 0 {
		return id
	}
	var existing struct {
		Id int64
	}
	err := dao.SysMenu.Ctx(ctx).
		Where(dao.SysMenu.Columns().Name, name).
		Fields(dao.SysMenu.Columns().Id).
		Scan(&existing)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0
		}
		return 0
	}
	if existing.Id == 0 {
		return 0
	}
	menuIDs[key] = existing.Id
	return existing.Id
}

func resolveRoleIDByName(ctx context.Context, name string, roleIDs map[string]uint) uint {
	tenantID := tenantIDFromCtx(ctx)
	key := mapKey(tenantID, name)
	if id, ok := roleIDs[key]; ok && id > 0 {
		return id
	}
	var existing struct {
		Id int64
	}
	err := dao.SysRole.Ctx(ctx).
		Where(dao.SysRole.Columns().Name, name).
		Fields(dao.SysRole.Columns().Id).
		Scan(&existing)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0
		}
		return 0
	}
	if existing.Id == 0 {
		return 0
	}
	roleIDs[key] = uint(existing.Id)
	return uint(existing.Id)
}

func resolvePermissionIDByName(ctx context.Context, name string, permIDs map[string]uint) uint {
	tenantID := tenantIDFromCtx(ctx)
	key := mapKey(tenantID, name)
	if id, ok := permIDs[key]; ok && id > 0 {
		return id
	}
	var existing struct {
		Id int64
	}
	err := dao.SysPermission.Ctx(ctx).
		Where(dao.SysPermission.Columns().Name, name).
		Fields(dao.SysPermission.Columns().Id).
		Scan(&existing)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0
		}
		return 0
	}
	if existing.Id == 0 {
		return 0
	}
	permIDs[key] = uint(existing.Id)
	return uint(existing.Id)
}

func resolveUserID(ctx context.Context, tenantID string, userID *int64, username string, userIDs map[string]int64) (int64, error) {
	if userID != nil && *userID > 0 {
		return *userID, nil
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return 0, fmt.Errorf("username is required for user role mapping")
	}
	key := mapKey(tenantID, username)
	if id, ok := userIDs[key]; ok && id > 0 {
		return id, nil
	}
	var existing struct {
		Id int64
	}
	err := dao.SysUser.Ctx(ctx).
		Where(dao.SysUser.Columns().Username, username).
		Fields(dao.SysUser.Columns().Id).
		Scan(&existing)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found: %s", username)
		}
		return 0, err
	}
	if existing.Id == 0 {
		return 0, fmt.Errorf("user not found: %s", username)
	}
	userIDs[key] = existing.Id
	return existing.Id, nil
}

func resolveRoleID(ctx context.Context, tenantID string, roleID *int64, roleName string, roleIDs map[string]uint) (int64, error) {
	if roleID != nil && *roleID > 0 {
		return *roleID, nil
	}
	roleName = strings.TrimSpace(roleName)
	if roleName == "" {
		return 0, fmt.Errorf("role name is required for mapping")
	}
	id := resolveRoleIDByName(ctx, roleName, roleIDs)
	if id == 0 {
		return 0, fmt.Errorf("role not found: %s", roleName)
	}
	return int64(id), nil
}

func resolvePermissionID(ctx context.Context, tenantID string, permID *int64, permName string, permIDs map[string]uint) (int64, error) {
	if permID != nil && *permID > 0 {
		return *permID, nil
	}
	permName = strings.TrimSpace(permName)
	if permName == "" {
		return 0, fmt.Errorf("permission name is required for mapping")
	}
	id := resolvePermissionIDByName(ctx, permName, permIDs)
	if id == 0 {
		return 0, fmt.Errorf("permission not found: %s", permName)
	}
	return int64(id), nil
}

func tenantIDFromCtx(ctx context.Context) string {
	if v := ctx.Value(consts.CtxKeyTenantID); v != nil {
		if tenantID, ok := v.(string); ok && strings.TrimSpace(tenantID) != "" {
			return tenantID
		}
	}
	return consts.DefaultTenantID
}
