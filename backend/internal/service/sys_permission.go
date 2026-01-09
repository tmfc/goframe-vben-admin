package service

import (
	"context"

	"backend/internal/dao"
	"backend/internal/model"
	"backend/internal/model/entity"

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
	CreatePermission(ctx context.Context, in model.SysPermissionCreateIn) (id uint, err error)
	GetPermission(ctx context.Context, in model.SysPermissionGetIn) (out *model.SysPermissionGetOut, err error)
	UpdatePermission(ctx context.Context, in model.SysPermissionUpdateIn) (err error)
	DeletePermission(ctx context.Context, in model.SysPermissionDeleteIn) (err error)
	GetPermissionList(ctx context.Context, in model.SysPermissionListIn) (out *model.SysPermissionListOut, err error)
	GetPermissionTree(ctx context.Context) (out []*model.SysPermissionTreeOut, err error)
}

// CreatePermission creates a new permission.
func (s *sSysPermission) CreatePermission(ctx context.Context, in model.SysPermissionCreateIn) (id uint, err error) {
	// Validate input
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return 0, err
	}
	// Check if permission name already exists
	count, err := dao.SysPermission.Ctx(ctx).
		Where(dao.SysPermission.Columns().Name, in.Name).
		Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, gerror.NewCodef(gcode.CodeValidationFailed, "Permission name '%s' already exists", in.Name)
	}

	columns := dao.SysPermission.Columns()
	result, err := dao.SysPermission.Ctx(ctx).Data(g.Map{
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

// GetPermission retrieves a permission by ID.
func (s *sSysPermission) GetPermission(ctx context.Context, in model.SysPermissionGetIn) (out *model.SysPermissionGetOut, err error) {
	out = &model.SysPermissionGetOut{}
	err = dao.SysPermission.Ctx(ctx).
		Where(dao.SysPermission.Columns().Id, in.Id).
		Scan(&out.SysPermission)
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
	count, err := dao.SysPermission.Ctx(ctx).
		Where(dao.SysPermission.Columns().Id, in.Id).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Permission with ID %d not found", in.Id)
	}

	columns := dao.SysPermission.Columns()
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
	_, err = dao.SysPermission.Ctx(ctx).
		Data(updateData).
		Where(dao.SysPermission.Columns().Id, in.Id).
		Update()
	return err
}

// DeletePermission deletes a permission by ID.
func (s *sSysPermission) DeletePermission(ctx context.Context, in model.SysPermissionDeleteIn) (err error) {
	// Check if permission exists
	count, err := dao.SysPermission.Ctx(ctx).
		Where(dao.SysPermission.Columns().Id, in.Id).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return gerror.NewCodef(gcode.CodeNotFound, "Permission with ID %d not found", in.Id)
	}

	var permission entity.SysPermission
	if err = dao.SysPermission.Ctx(ctx).
		Where(dao.SysPermission.Columns().Id, in.Id).
		Scan(&permission); err != nil {
		return err
	}

	menuCount, err := dao.SysMenu.Ctx(ctx).
		Where(dao.SysMenu.Columns().Name, permission.Name).
		Where(dao.SysMenu.Columns().PermissionCode+" <> ?", "").
		Count()
	if err != nil {
		return err
	}
	if menuCount > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Permission with ID %d is managed by menu", in.Id)
	}

	childCount, err := dao.SysPermission.Ctx(ctx).
		Where(dao.SysPermission.Columns().ParentId, in.Id).
		Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Permission with ID %d has child permissions", in.Id)
	}

	roleBindingCount, err := dao.SysRolePermission.Ctx(ctx).
		Where(dao.SysRolePermission.Columns().PermissionId, in.Id).
		Count()
	if err != nil {
		return err
	}
	if roleBindingCount > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Permission with ID %d is assigned to roles", in.Id)
	}

	_, err = dao.SysPermission.Ctx(ctx).
		Where(dao.SysPermission.Columns().Id, in.Id).
		Delete()
	return err
}

// GetPermissionList lists permissions with pagination and filters.
func (s *sSysPermission) GetPermissionList(ctx context.Context, in model.SysPermissionListIn) (out *model.SysPermissionListOut, err error) {
	out = &model.SysPermissionListOut{}

	m := dao.SysPermission.Ctx(ctx)
	if in.Name != "" {
		m = m.WhereLike(dao.SysPermission.Columns().Name, "%"+in.Name+"%")
	}
	if in.Status != "" {
		m = m.Where(dao.SysPermission.Columns().Status, in.Status)
	}

	out.Total, err = m.Count()
	if err != nil {
		return nil, err
	}

	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	err = m.Page(in.Page, in.PageSize).Scan(&out.Items)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetPermissionTree retrieves the permission tree structure.
func (s *sSysPermission) GetPermissionTree(ctx context.Context) (out []*model.SysPermissionTreeOut, err error) {
	var permissions []*entity.SysPermission
	err = dao.SysPermission.Ctx(ctx).
		Where(dao.SysPermission.Columns().Status, 1).
		OrderAsc(dao.SysPermission.Columns().Id).
		Scan(&permissions)
	if err != nil {
		return nil, err
	}

	if len(permissions) == 0 {
		return []*model.SysPermissionTreeOut{}, nil
	}

	permissionMap := make(map[int64]*model.SysPermissionTreeOut)
	for _, perm := range permissions {
		permissionMap[perm.Id] = &model.SysPermissionTreeOut{
			Id:          perm.Id,
			Name:        perm.Name,
			Description: perm.Description,
			ParentId:    perm.ParentId,
			Status:      perm.Status,
			CreatedAt:   perm.CreatedAt.Format("Y-m-d H:i:s"),
			UpdatedAt:   perm.UpdatedAt.Format("Y-m-d H:i:s"),
		}
	}

	var roots []*model.SysPermissionTreeOut
	for _, perm := range permissions {
		item := permissionMap[perm.Id]
		if perm.ParentId == 0 {
			roots = append(roots, item)
		} else {
			if parent, exists := permissionMap[perm.ParentId]; exists {
				parent.Children = append(parent.Children, item)
			}
		}
	}

	return roots, nil
}
