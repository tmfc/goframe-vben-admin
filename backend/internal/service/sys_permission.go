package service

import (

	"context"



	"backend/internal/model"

	"backend/internal/dao"



	"github.com/gogf/gf/v2/errors/gcode"

	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/v2/frame/g"

)

var localSysPermission ISysPermission

func SysPermission() ISysPermission {
	return localSysPermission
}

func RegisterSysPermission(i ISysPermission) {
	localSysPermission = i
}

type sSysPermission struct{}

func init() {
	RegisterSysPermission(NewSysPermission())
}

func NewSysPermission() *sSysPermission {
	return &sSysPermission{}
}

// ISysPermission defines the interface for sys_permission service.
type ISysPermission interface {
	CreatePermission(ctx context.Context, in model.SysPermissionCreateIn) (err error)
	GetPermission(ctx context.Context, in model.SysPermissionGetIn) (out *model.SysPermissionGetOut, err error)
	UpdatePermission(ctx context.Context, in model.SysPermissionUpdateIn) (err error)
	DeletePermission(ctx context.Context, in model.SysPermissionDeleteIn) (err error)
}

// CreatePermission creates a new permission.
func (s *sSysPermission) CreatePermission(ctx context.Context, in model.SysPermissionCreateIn) (err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}
	// Check if permission name already exists
	count, err := dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Name, in.Name).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Permission name '%s' already exists", in.Name)
	}

	_, err = dao.SysPermission.Ctx(ctx).Data(in).Insert()
	return err
}

// GetPermission retrieves a permission by ID.
func (s *sSysPermission) GetPermission(ctx context.Context, in model.SysPermissionGetIn) (out *model.SysPermissionGetOut, err error) {
	out = &model.SysPermissionGetOut{}
	err = dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Id, in.Id).Scan(&out.SysPermission)
	if err != nil {
		return nil, err
	}
	if out.SysPermission == nil {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Permission with ID %d not found", in.Id)
	}
	return out, nil
}

// UpdatePermission updates an existing permission.
func (s *sSysPermission) UpdatePermission(ctx context.Context, in model.SysPermissionUpdateIn) (err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}

	// Check if permission exists
	count, err := dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Id, in.Id).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Permission with ID %d not found", in.Id)
	}

	_, err = dao.SysPermission.Ctx(ctx).Data(in).Where(dao.SysPermission.Columns().Id, in.Id).Update()
	return err
}

// DeletePermission deletes a permission by ID.
func (s *sSysPermission) DeletePermission(ctx context.Context, in model.SysPermissionDeleteIn) (err error) {
	// Check if permission exists
	count, err := dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Id, in.Id).Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Permission with ID %d not found", in.Id)
	}

	_, err = dao.SysPermission.Ctx(ctx).Where(dao.SysPermission.Columns().Id, in.Id).Delete()
	return err
}
