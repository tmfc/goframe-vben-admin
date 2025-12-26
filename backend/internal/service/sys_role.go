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
	CreateRole(ctx context.Context, in model.SysRoleCreateIn) (err error)
	GetRole(ctx context.Context, in model.SysRoleGetIn) (out *model.SysRoleGetOut, err error)
	UpdateRole(ctx context.Context, in model.SysRoleUpdateIn) (err error)
	DeleteRole(ctx context.Context, in model.SysRoleDeleteIn) (err error)
}

// CreateRole creates a new role.
func (s *sSysRole) CreateRole(ctx context.Context, in model.SysRoleCreateIn) (err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}
	// Check if role name already exists
	count, err := dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Name, in.Name).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Role name '%s' already exists", in.Name)
	}

	_, err = dao.SysRole.Ctx(ctx).Data(in).Insert()
	return err
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

	// Check if role exists
	count, err := dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, in.Id).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.Id)
	}

	_, err = dao.SysRole.Ctx(ctx).Data(in).Where(dao.SysRole.Columns().Id, in.Id).Update()
	return err
}

// DeleteRole deletes a role by ID.
func (s *sSysRole) DeleteRole(ctx context.Context, in model.SysRoleDeleteIn) (err error) {
	// Check if role exists
	count, err := dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, in.Id).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Role with ID %d not found", in.Id)
	}

	_, err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns().Id, in.Id).Delete()
	return err
}
