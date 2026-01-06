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
	AssignRoleToUser(ctx context.Context, in model.AssignRoleToUserIn) (out *model.AssignRoleToUserOut, err error)
	RemoveRoleFromUser(ctx context.Context, in model.RemoveRoleFromUserIn) (out *model.RemoveRoleFromUserOut, err error)
	GetUserRoles(ctx context.Context, in model.GetUserRolesIn) (out *model.GetUserRolesOut, err error)
	AssignRolesToUser(ctx context.Context, in model.AssignRolesToUserIn) (out *model.AssignRolesToUserOut, err error)
	GetUsersByRole(ctx context.Context, in model.GetUsersByRoleIn) (out *model.GetUsersByRoleOut, err error)
	SyncRoleToCasbin(ctx context.Context, roleID uint) (err error)
	RemoveRoleFromCasbin(ctx context.Context, roleID uint) (err error)
	SyncAllRolesToCasbin(ctx context.Context) (err error)
	AssignPermissionToRole(ctx context.Context, in model.SysRolePermissionIn) (err error)
	RemovePermissionFromRole(ctx context.Context, in model.SysRolePermissionIn) (err error)
	GetRolePermissions(ctx context.Context, in model.SysRoleGetIn) (out *model.SysRolePermissionOut, err error)
	AssignPermissionsToRole(ctx context.Context, in model.SysRolePermissionsIn) (err error)
	GetPermissionsByUser(ctx context.Context, in model.UserPermissionsIn) (out *model.UserPermissionsOut, err error)
}

// CreateRole creates a new role.
func (s *sSysRole) CreateRole(ctx context.Context, in model.SysRoleCreateIn) (id uint, err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return 0, err
	}
	// Check if role name already exists
	count, err := dao.SysRole.Ctx(ctx).
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

	if in.ParentId != 0 {
		if in.ParentId == in.Id {
			return gerror.NewCodef(gcode.CodeValidationFailed, "Role parent cannot be itself")
		}
		parentCount, err := dao.SysRole.Ctx(ctx).
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
	// Check if role exists
	var existingRole entity.SysRole
	err = dao.SysRole.Ctx(ctx).
		Where(dao.SysRole.Columns().Id, in.Id).
		Scan(&existingRole)
	if err != nil {
		return gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.Id)
	}
	if existingRole.Id == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.Id)
	}

	childCount, err := dao.SysRole.Ctx(ctx).
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
	err := dao.SysRolePermission.CtxNoTenant(ctx).
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

// AssignRoleToUser assigns a role to a user.
func (s *sSysRole) AssignRoleToUser(ctx context.Context, in model.AssignRoleToUserIn) (out *model.AssignRoleToUserOut, err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return nil, err
	}

	tenantID := resolveTenantID(ctx)

	// Check if user exists
	userCount, err := dao.SysUser.Ctx(ctx).
		Where(dao.SysUser.Columns().Id, in.UserId).
		Where(dao.SysUser.Columns().TenantId, tenantID).
		Count()
	if err != nil {
		return nil, err
	}
	if userCount == 0 {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "User with ID %s not found", in.UserId)
	}

	// Check if role exists
	roleCount, err := dao.SysRole.Ctx(ctx).
		Where(dao.SysRole.Columns().Id, in.RoleId).
		Where(dao.SysRole.Columns().TenantId, tenantID).
		Count()
	if err != nil {
		return nil, err
	}
	if roleCount == 0 {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.RoleId)
	}

	// Check if assignment already exists
	existingCount, err := dao.SysUserRole.Ctx(ctx).
		Where(dao.SysUserRole.Columns().TenantId, tenantID).
		Where(dao.SysUserRole.Columns().UserId, in.UserId).
		Where(dao.SysUserRole.Columns().RoleId, in.RoleId).
		Count()
	if err != nil {
		return nil, err
	}
	if existingCount > 0 {
		return &model.AssignRoleToUserOut{Success: true}, nil
	}

	// Create assignment
	_, err = dao.SysUserRole.Ctx(ctx).Data(g.Map{
		dao.SysUserRole.Columns().TenantId:  tenantID,
		dao.SysUserRole.Columns().UserId:    in.UserId,
		dao.SysUserRole.Columns().RoleId:    in.RoleId,
		dao.SysUserRole.Columns().CreatedBy: in.CreatedBy,
	}).Insert()
	if err != nil {
		return nil, err
	}

	return &model.AssignRoleToUserOut{Success: true}, nil
}

// RemoveRoleFromUser removes a role from a user.
func (s *sSysRole) RemoveRoleFromUser(ctx context.Context, in model.RemoveRoleFromUserIn) (out *model.RemoveRoleFromUserOut, err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return nil, err
	}

	tenantID := resolveTenantID(ctx)

	// Delete the assignment
	_, err = dao.SysUserRole.Ctx(ctx).
		Where(dao.SysUserRole.Columns().TenantId, tenantID).
		Where(dao.SysUserRole.Columns().UserId, in.UserId).
		Where(dao.SysUserRole.Columns().RoleId, in.RoleId).
		Delete()
	if err != nil {
		return nil, err
	}

	return &model.RemoveRoleFromUserOut{Success: true}, nil
}

// GetUserRoles retrieves all roles assigned to a user.
func (s *sSysRole) GetUserRoles(ctx context.Context, in model.GetUserRolesIn) (out *model.GetUserRolesOut, err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return nil, err
	}

	tenantID := resolveTenantID(ctx)

	// Get role IDs from user_role table
	var userRoles []struct {
		RoleId uint
	}
	err = dao.SysUserRole.Ctx(ctx).
		Fields(dao.SysUserRole.Columns().RoleId).
		Where(dao.SysUserRole.Columns().TenantId, tenantID).
		Where(dao.SysUserRole.Columns().UserId, in.UserId).
		Scan(&userRoles)
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return &model.GetUserRolesOut{Roles: []*entity.SysRole{}}, nil
	}

	// Get role IDs
	roleIds := make([]uint, len(userRoles))
	for i, ur := range userRoles {
		roleIds[i] = ur.RoleId
	}

	// Get role details
	var roles []*entity.SysRole
	err = dao.SysRole.Ctx(ctx).
		WhereIn(dao.SysRole.Columns().Id, roleIds).
		Where(dao.SysRole.Columns().TenantId, tenantID).
		Scan(&roles)
	if err != nil {
		return nil, err
	}

	return &model.GetUserRolesOut{Roles: roles}, nil
}

// AssignRolesToUser assigns multiple roles to a user.
func (s *sSysRole) AssignRolesToUser(ctx context.Context, in model.AssignRolesToUserIn) (out *model.AssignRolesToUserOut, err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return nil, err
	}

	if len(in.RoleIds) == 0 {
		return nil, gerror.NewCodef(gcode.CodeValidationFailed, "Role IDs cannot be empty")
	}

	tenantID := resolveTenantID(ctx)

	// Check if user exists
	userCount, err := dao.SysUser.Ctx(ctx).
		Where(dao.SysUser.Columns().Id, in.UserId).
		Where(dao.SysUser.Columns().TenantId, tenantID).
		Count()
	if err != nil {
		return nil, err
	}
	if userCount == 0 {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "User with ID %s not found", in.UserId)
	}

	// Check if all roles exist
	roleCount, err := dao.SysRole.Ctx(ctx).
		WhereIn(dao.SysRole.Columns().Id, in.RoleIds).
		Where(dao.SysRole.Columns().TenantId, tenantID).
		Count()
	if err != nil {
		return nil, err
	}
	if int(roleCount) != len(in.RoleIds) {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "One or more roles not found")
	}

	// Remove existing role assignments for this user
	_, err = dao.SysUserRole.Ctx(ctx).
		Where(dao.SysUserRole.Columns().TenantId, tenantID).
		Where(dao.SysUserRole.Columns().UserId, in.UserId).
		Delete()
	if err != nil {
		return nil, err
	}

	// Create new assignments
	if len(in.RoleIds) > 0 {
		insertData := make([]g.Map, len(in.RoleIds))
		for i, roleId := range in.RoleIds {
			insertData[i] = g.Map{
				dao.SysUserRole.Columns().TenantId:  tenantID,
				dao.SysUserRole.Columns().UserId:    in.UserId,
				dao.SysUserRole.Columns().RoleId:    roleId,
				dao.SysUserRole.Columns().CreatedBy: in.CreatedBy,
			}
		}
		_, err = dao.SysUserRole.Ctx(ctx).Data(insertData).Insert()
		if err != nil {
			return nil, err
		}
	}

	return &model.AssignRolesToUserOut{Success: true, Count: len(in.RoleIds)}, nil
}

// GetUsersByRole retrieves all users assigned to a role.
func (s *sSysRole) GetUsersByRole(ctx context.Context, in model.GetUsersByRoleIn) (out *model.GetUsersByRoleOut, err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return nil, err
	}

	tenantID := resolveTenantID(ctx)

	// Check if role exists
	roleCount, err := dao.SysRole.Ctx(ctx).
		Where(dao.SysRole.Columns().Id, in.RoleId).
		Where(dao.SysRole.Columns().TenantId, tenantID).
		Count()
	if err != nil {
		return nil, err
	}
	if roleCount == 0 {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.RoleId)
	}

	// Get user IDs from user_role table
	var userRoles []struct {
		UserId string
	}
	err = dao.SysUserRole.Ctx(ctx).
		Fields(dao.SysUserRole.Columns().UserId).
		Where(dao.SysUserRole.Columns().TenantId, tenantID).
		Where(dao.SysUserRole.Columns().RoleId, in.RoleId).
		Scan(&userRoles)
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return &model.GetUsersByRoleOut{Users: []*entity.SysUser{}}, nil
	}

	// Get user IDs
	userIds := make([]string, len(userRoles))
	for i, ur := range userRoles {
		userIds[i] = ur.UserId
	}

	// Get user details
	var users []*entity.SysUser
	err = dao.SysUser.Ctx(ctx).
		WhereIn(dao.SysUser.Columns().Id, userIds).
		Where(dao.SysUser.Columns().TenantId, tenantID).
		Scan(&users)
	if err != nil {
		return nil, err
	}

	return &model.GetUsersByRoleOut{Users: users}, nil
}

// AssignPermissionToRole assigns a permission to a role.
func (s *sSysRole) AssignPermissionToRole(ctx context.Context, in model.SysRolePermissionIn) (err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}

	// Check if role exists
	roleCount, err := dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, in.RoleID).Count()
	if err != nil {
		return err
	}
	if roleCount == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.RoleID)
	}

	// Check if permission exists
	permissionCount, err := dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Id, in.PermissionID).Count()
	if err != nil {
		return err
	}
	if permissionCount == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Permission with ID %d not found", in.PermissionID)
	}

	// Check if association already exists
	count, err := dao.SysRolePermission.Ctx(ctx).
		Where(dao.SysRolePermission.Columns().RoleId, in.RoleID).
		Where(dao.SysRolePermission.Columns().PermissionId, in.PermissionID).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Permission with ID %d is already assigned to role with ID %d", in.PermissionID, in.RoleID)
	}

	_, err = dao.SysRolePermission.Ctx(ctx).Data(g.Map{
		dao.SysRolePermission.Columns().RoleId:       in.RoleID,
		dao.SysRolePermission.Columns().PermissionId: in.PermissionID,
	}).Insert()
	return err
}

// RemovePermissionFromRole removes a permission from a role.
func (s *sSysRole) RemovePermissionFromRole(ctx context.Context, in model.SysRolePermissionIn) (err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}

	// Check if association exists
	count, err := dao.SysRolePermission.Ctx(ctx).
		Where(dao.SysRolePermission.Columns().RoleId, in.RoleID).
		Where(dao.SysRolePermission.Columns().PermissionId, in.PermissionID).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Permission with ID %d is not assigned to role with ID %d", in.PermissionID, in.RoleID)
	}

	_, err = dao.SysRolePermission.Ctx(ctx).
		Where(dao.SysRolePermission.Columns().RoleId, in.RoleID).
		Where(dao.SysRolePermission.Columns().PermissionId, in.PermissionID).
		Delete()
	return err
}

// GetRolePermissions retrieves all permissions for a role.
func (s *sSysRole) GetRolePermissions(ctx context.Context, in model.SysRoleGetIn) (out *model.SysRolePermissionOut, err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return nil, err
	}

	// Check if role exists
	roleCount, err := dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, in.Id).Count()
	if err != nil {
		return nil, err
	}
	if roleCount == 0 {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.Id)
	}

	var permissionIDs []uint
	result, err := dao.SysRolePermission.Ctx(ctx).
		Where(dao.SysRolePermission.Columns().RoleId, in.Id).
		Fields(dao.SysRolePermission.Columns().PermissionId).
		All()
	if err != nil {
		return nil, err
	}
	for _, row := range result {
		permissionID := row[dao.SysRolePermission.Columns().PermissionId].Int64()
		permissionIDs = append(permissionIDs, uint(permissionID))
	}

	out = &model.SysRolePermissionOut{
		PermissionIDs: permissionIDs,
	}
	return out, nil
}

// AssignPermissionsToRole assigns multiple permissions to a role.
func (s *sSysRole) AssignPermissionsToRole(ctx context.Context, in model.SysRolePermissionsIn) (err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}

	// Check if role exists
	roleCount, err := dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, in.RoleID).Count()
	if err != nil {
		return err
	}
	if roleCount == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.RoleID)
	}

	// Remove existing permissions
	_, err = dao.SysRolePermission.Ctx(ctx).
		Where(dao.SysRolePermission.Columns().RoleId, in.RoleID).
		Delete()
	if err != nil {
		return err
	}

	// Add new permissions
	if len(in.PermissionIDs) > 0 {
		data := g.List{}
		for _, permissionID := range in.PermissionIDs {
			data = append(data, g.Map{
				dao.SysRolePermission.Columns().RoleId:       in.RoleID,
				dao.SysRolePermission.Columns().PermissionId: permissionID,
			})
		}
		_, err = dao.SysRolePermission.Ctx(ctx).Data(data).Insert()
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPermissionsByUser retrieves all permissions for a user.
func (s *sSysRole) GetPermissionsByUser(ctx context.Context, in model.UserPermissionsIn) (out *model.UserPermissionsOut, err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return nil, err
	}

	// 1) 从 casbin_rule 的 g 规则获取用户绑定的角色 ID（v1 存 role_id）
	var roleIDs []uint
	casbinResult, err := dao.CasbinRule.Ctx(ctx).
		Where(dao.CasbinRule.Columns().Ptype, "g").
		Where(dao.CasbinRule.Columns().V0, "u_"+in.UserID).
		Fields(dao.CasbinRule.Columns().V1).
		All()
	if err != nil {
		return nil, err
	}
	for _, row := range casbinResult {
		roleID := row[dao.CasbinRule.Columns().V1].Int64()
		if roleID == 0 {
			continue
		}
		roleIDs = append(roleIDs, uint(roleID))
	}
	if len(roleIDs) == 0 {
		return &model.UserPermissionsOut{PermissionIDs: []uint{}}, nil
	}

	// 2) 根据 role_id 从 sys_role_permission 直接查出权限 ID
	var permissionIDs []uint
	permResult, err := dao.SysRolePermission.Ctx(ctx).
		WhereIn(dao.SysRolePermission.Columns().RoleId, roleIDs).
		Fields(dao.SysRolePermission.Columns().PermissionId).
		All()
	if err != nil {
		return nil, err
	}
	for _, row := range permResult {
		permissionID := row[dao.SysRolePermission.Columns().PermissionId].Int64()
		if permissionID == 0 {
			continue
		}
		permissionIDs = append(permissionIDs, uint(permissionID))
	}

	out = &model.UserPermissionsOut{PermissionIDs: permissionIDs}
	return out, nil
}
