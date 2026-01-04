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

type sSysRole struct{}

var localSysRole ISysRole

func SysRole() ISysRole {
	return localSysRole
}

func RegisterSysRole(i ISysRole) {
	localSysRole = i
}

func init() {
	RegisterSysRole(NewSysRole())
}

func NewSysRole() *sSysRole {
	return &sSysRole{}
}

// ISysRole defines the interface for sys_role service.
type ISysRole interface {
	CreateRole(ctx context.Context, in model.SysRoleCreateIn) (id uint, err error)
	GetRole(ctx context.Context, in model.SysRoleGetIn) (out *model.SysRoleGetOut, err error)
	UpdateRole(ctx context.Context, in model.SysRoleUpdateIn) (err error)
	DeleteRole(ctx context.Context, in model.SysRoleDeleteIn) (err error)
	SyncRoleToCasbin(ctx context.Context, roleID uint) (err error)
	RemoveRoleFromCasbin(ctx context.Context, roleID uint) (err error)
	SyncAllRolesToCasbin(ctx context.Context) (err error)
}

// CreateRole creates a new role.
func (s *sSysRole) CreateRole(ctx context.Context, in model.SysRoleCreateIn) (id uint, err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return 0, err
	}
	tenantID := resolveTenantID(ctx)
	// Check if role name already exists
	count, err := dao.SysRole.Ctx(ctx).
		Where("tenant_id", tenantID).
		Where(dao.SysRole.Columns().Name, in.Name).
		Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, gerror.NewCodef(gcode.CodeValidationFailed, "Role name '%s' already exists", in.Name)
	}

	columns := dao.SysRole.Columns()
	result, err := dao.SysRole.Ctx(ctx).Data(g.Map{
		"tenant_id":         tenantID,
		columns.Name:        in.Name,
		columns.Description: in.Description,
		columns.ParentId:    in.ParentId,
		columns.Status:      in.Status,
		columns.CreatorId:   in.CreatorId,
		columns.ModifierId:  in.ModifierId,
		columns.DeptId:      in.DeptId,
	}).Insert()
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	createdID := uint(lastInsertId)
	if err := s.SyncRoleToCasbin(ctx, createdID); err != nil {
		return createdID, err
	}
	return createdID, nil
}

// GetRole retrieves a role by ID.
func (s *sSysRole) GetRole(ctx context.Context, in model.SysRoleGetIn) (out *model.SysRoleGetOut, err error) {
	out = &model.SysRoleGetOut{}
	err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, in.Id).Scan(&out.SysRole)
	if err != nil {
		return nil, err
	}
	if out.SysRole == nil || out.SysRole.Id == 0 {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.Id)
	}
	return out, nil
}

// UpdateRole updates an existing role.
func (s *sSysRole) UpdateRole(ctx context.Context, in model.SysRoleUpdateIn) (err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}

	tenantID := resolveTenantID(ctx)
	if in.ParentId != 0 {
		if in.ParentId == in.Id {
			return gerror.NewCodef(gcode.CodeValidationFailed, "Role parent cannot be itself")
		}
		parentCount, err := dao.SysRole.Ctx(ctx).
			Where("tenant_id", tenantID).
			Where(dao.SysRole.Columns().Id, in.ParentId).
			Count()
		if err != nil {
			return err
		}
		if parentCount == 0 {
			return gerror.NewCodef(gcode.CodeValidationFailed, "Parent role with ID %d not found", in.ParentId)
		}
	}
	// Check if role exists
	var existingRole entity.SysRole
	err = dao.SysRole.Ctx(ctx).
		Where("tenant_id", tenantID).
		Where(dao.SysRole.Columns().Id, in.Id).
		Scan(&existingRole)
	if err != nil {
		return gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.Id)
	}
	if existingRole.Id == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.Id)
	}

	columns := dao.SysRole.Columns()
	updateData := g.Map{
		columns.Name:        in.Name,
		columns.Description: in.Description,
		columns.ParentId:    in.ParentId,
		columns.Status:      in.Status,
		columns.ModifierId:  in.ModifierId,
		columns.DeptId:      in.DeptId,
	}
	if in.UpdatedAt != nil {
		updateData[columns.UpdatedAt] = in.UpdatedAt
	}
	_, err = dao.SysRole.Ctx(ctx).
		Data(updateData).
		Where("tenant_id", tenantID).
		Where(dao.SysRole.Columns().Id, in.Id).
		Update()
	if err != nil {
		return err
	}
	if strings.TrimSpace(existingRole.Name) != strings.TrimSpace(in.Name) {
		if err := removeRolePolicies(ctx, existingRole.Name, existingRole.TenantId); err != nil {
			return err
		}
	}
	return s.SyncRoleToCasbin(ctx, in.Id)
}

// DeleteRole deletes a role by ID.
func (s *sSysRole) DeleteRole(ctx context.Context, in model.SysRoleDeleteIn) (err error) {
	tenantID := resolveTenantID(ctx)
	// Check if role exists
	var existingRole entity.SysRole
	err = dao.SysRole.Ctx(ctx).
		Where("tenant_id", tenantID).
		Where(dao.SysRole.Columns().Id, in.Id).
		Scan(&existingRole)
	if err != nil {
		return gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.Id)
	}
	if existingRole.Id == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.Id)
	}

	childCount, err := dao.SysRole.Ctx(ctx).
		Where("tenant_id", tenantID).
		Where(dao.SysRole.Columns().ParentId, in.Id).
		Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Role with ID %d has child roles", in.Id)
	}

	bindingCount, err := dao.SysRolePermission.Ctx(ctx).
		Where(dao.SysRolePermission.Columns().RoleId, in.Id).
		Count()
	if err != nil {
		return err
	}
	if bindingCount > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Role with ID %d has permission bindings", in.Id)
	}

	_, err = dao.SysRole.Ctx(ctx).
		Where("tenant_id", tenantID).
		Where(dao.SysRole.Columns().Id, in.Id).
		Delete()
	if err != nil {
		return err
	}
	return removeRolePolicies(ctx, existingRole.Name, existingRole.TenantId)
}

type rolePermissionRecord struct {
	Code  string `json:"code" orm:"code"`
	Scope string `json:"scope" orm:"scope"`
}

// SyncRoleToCasbin syncs a role's permissions to Casbin.
func (s *sSysRole) SyncRoleToCasbin(ctx context.Context, roleID uint) (err error) {
	role, err := fetchRole(ctx, roleID)
	if err != nil {
		return err
	}

	roleName := strings.TrimSpace(role.Name)
	if roleName == "" {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Role name cannot be empty")
	}
	domain := NormalizeDomain(role.TenantId)

	enforcer, err := Casbin(ctx)
	if err != nil {
		return err
	}
	if _, err := enforcer.RemoveFilteredNamedPolicy("p", 0, roleName, domain); err != nil {
		return err
	}

	records, err := fetchRolePermissions(ctx, roleID)
	if err != nil || len(records) == 0 {
		return err
	}

	policies := make([][]string, 0, len(records))
	seen := make(map[string]struct{})
	for _, record := range records {
		code := strings.TrimSpace(record.Code)
		if code == "" {
			continue
		}
		scope := strings.TrimSpace(record.Scope)
		if scope == "" {
			scope = "*"
		}
		key := code + "\x00" + scope
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		policies = append(policies, []string{roleName, domain, code, scope})
	}
	if len(policies) == 0 {
		return nil
	}
	for _, policy := range policies {
		if _, err := enforcer.AddPolicy(policy); err != nil {
			return err
		}
	}
	return nil
}

// RemoveRoleFromCasbin removes a role's policies from Casbin.
func (s *sSysRole) RemoveRoleFromCasbin(ctx context.Context, roleID uint) (err error) {
	role, err := fetchRole(ctx, roleID)
	if err != nil {
		return err
	}
	return removeRolePolicies(ctx, role.Name, role.TenantId)
}

// SyncAllRolesToCasbin syncs all roles to Casbin.
func (s *sSysRole) SyncAllRolesToCasbin(ctx context.Context) (err error) {
	var roles []entity.SysRole
	if err := dao.SysRole.Ctx(ctx).Scan(&roles); err != nil {
		return err
	}
	for _, role := range roles {
		if err := s.SyncRoleToCasbin(ctx, uint(role.Id)); err != nil {
			return err
		}
	}
	return nil
}

func fetchRole(ctx context.Context, roleID uint) (*entity.SysRole, error) {
	var role entity.SysRole
	if err := dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, roleID).Scan(&role); err != nil {
		return nil, err
	}
	if role.Id == 0 {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", roleID)
	}
	return &role, nil
}

func fetchRolePermissions(ctx context.Context, roleID uint) ([]rolePermissionRecord, error) {
	var records []rolePermissionRecord
	err := dao.SysRolePermission.Ctx(ctx).
		As("rp").
		LeftJoin("sys_permission p", "rp.permission_id = p.id").
		Fields("p.name as code", "rp.scope as scope").
		Where("rp.role_id", roleID).
		Scan(&records)
	return records, err
}

func removeRolePolicies(ctx context.Context, roleName, tenantID string) error {
	roleName = strings.TrimSpace(roleName)
	if roleName == "" {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Role name cannot be empty")
	}
	domain := NormalizeDomain(tenantID)
	enforcer, err := Casbin(ctx)
	if err != nil {
		return err
	}
	_, err = enforcer.RemoveFilteredNamedPolicy("p", 0, roleName, domain)
	return err
}
