package service

import (
	"context"

	"backend/internal/dao"
	"backend/internal/model"

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
	return uint(lastInsertId), nil
}

// GetRole retrieves a role by ID.
func (s *sSysRole) GetRole(ctx context.Context, in model.SysRoleGetIn) (out *model.SysRoleGetOut, err error) {
	out = &model.SysRoleGetOut{}
	err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, in.Id).Scan(&out.SysRole)
	if err != nil {
		return nil, err
	}
	if out.SysRole == nil {
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
	count, err := dao.SysRole.Ctx(ctx).
		Where("tenant_id", tenantID).
		Where(dao.SysRole.Columns().Id, in.Id).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
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
	return err
}

// DeleteRole deletes a role by ID.
func (s *sSysRole) DeleteRole(ctx context.Context, in model.SysRoleDeleteIn) (err error) {
	tenantID := resolveTenantID(ctx)
	// Check if role exists
	count, err := dao.SysRole.Ctx(ctx).
		Where("tenant_id", tenantID).
		Where(dao.SysRole.Columns().Id, in.Id).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
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
	return err
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

	// Get user roles
	var roleIDs []uint
	// The role id is stored in v1, and the user id is stored in v0.
	// The user id is prefixed with "u_".
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
		roleIDs = append(roleIDs, uint(roleID))
	}

	if len(roleIDs) == 0 {
		return &model.UserPermissionsOut{PermissionIDs: []uint{}}, nil
	}

	// Get role permissions
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
		permissionIDs = append(permissionIDs, uint(permissionID))
	}

	out = &model.UserPermissionsOut{
		PermissionIDs: permissionIDs,
	}
	return out, nil
}