package service

import (

	"context"



	"backend/internal/model"

	"backend/internal/dao"



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
