package service

import (

	"context"

	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/errors/gcode"

	"github.com/gogf/gf/v2/errors/gerror"

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
		"tenant_id":          tenantID,
		columns.Name:         in.Name,
		columns.Description:  in.Description,
		columns.ParentId:     in.ParentId,
		columns.Status:       in.Status,
		columns.CreatorId:    in.CreatorId,
		columns.ModifierId:   in.ModifierId,
		columns.DeptId:       in.DeptId,
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
