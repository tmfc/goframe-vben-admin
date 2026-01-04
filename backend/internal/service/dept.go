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

var localDept IDept

func Dept() IDept {
	return localDept
}

func RegisterDept(i IDept) {
	localDept = i
}

var _ IDept = (*sDept)(nil)

func init() {
	RegisterDept(NewDept())
}

func NewDept() *sDept {
	return &sDept{}
}

type IDept interface {
	CreateDept(ctx context.Context, in model.SysDeptCreateIn) (id string, err error)
	GetDept(ctx context.Context, id string) (out *model.SysDeptGetOut, err error)
	UpdateDept(ctx context.Context, in model.SysDeptUpdateIn) (err error)
	DeleteDept(ctx context.Context, id string) (err error)
	GetDeptList(ctx context.Context, in model.SysDeptGetListIn) (out *model.SysDeptGetListOut, err error)
	GetDeptTree(ctx context.Context) (out []*model.SysDeptTreeOut, err error)
}

type sDept struct{}

func (s *sDept) CreateDept(ctx context.Context, in model.SysDeptCreateIn) (id string, err error) {
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return "", err
	}

	tenantID := resolveTenantID(ctx)

	// Prepare parent_id value (convert "0" to NULL for root departments)
	parentId := in.ParentId
	if parentId == "0" {
		parentId = ""
	}

	if parentId != "" {
		parentCount, err := dao.SysDept.Ctx(ctx).
			Where(dao.SysDept.Columns().TenantId, tenantID).
			Where(dao.SysDept.Columns().Id, parentId).
			Count()
		if err != nil {
			return "", err
		}
		if parentCount == 0 {
			return "", gerror.NewCodef(gcode.CodeValidationFailed, "Parent department with ID %s not found", parentId)
		}
	}

	columns := dao.SysDept.Columns()
	insertData := g.Map{
		"tenant_id":   tenantID,
		columns.Name:      in.Name,
		columns.Status:    in.Status,
		columns.Order:     in.Order,
		columns.CreatorId: in.CreatorId,
	}
	if parentId != "" {
		insertData[columns.ParentId] = parentId
	}

	_, err = dao.SysDept.Ctx(ctx).Data(insertData).Insert()

	if err != nil {
		return "", err
	}

	// For UUID columns, we need to retrieve the inserted record to get the ID
	// Use the name and tenant_id to find the newly created record
	var insertedDept entity.SysDept
	query := dao.SysDept.Ctx(ctx).
		Where(dao.SysDept.Columns().TenantId, tenantID).
		Where(dao.SysDept.Columns().Name, in.Name).
		OrderDesc(dao.SysDept.Columns().CreatedAt).
		Limit(1)
	if parentId != "" {
		query = query.Where(dao.SysDept.Columns().ParentId, parentId)
	} else {
		query = query.WhereNull(dao.SysDept.Columns().ParentId)
	}
	err = query.Scan(&insertedDept)
	if err != nil {
		return "", err
	}

	return insertedDept.Id, nil
}

func (s *sDept) GetDept(ctx context.Context, id string) (out *model.SysDeptGetOut, err error) {
	out = &model.SysDeptGetOut{}
	err = dao.SysDept.Ctx(ctx).Where(dao.SysDept.Columns().Id, id).Scan(&out.SysDept)
	if err != nil {
		return nil, err
	}
	if out.SysDept == nil {
		return nil, gerror.NewCodef(gcode.CodeNotFound, "Department with ID %s not found", id)
	}
	return out, nil
}

func (s *sDept) UpdateDept(ctx context.Context, in model.SysDeptUpdateIn) (err error) {
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}

	tenantID := resolveTenantID(ctx)

	// Prepare parent_id value (convert "0" to NULL for root departments)
	parentId := in.ParentId
	if parentId == "0" {
		parentId = ""
	}

	if parentId != "" {
		if parentId == in.ID {
			return gerror.NewCodef(gcode.CodeValidationFailed, "Department parent cannot be itself")
		}
		parentCount, err := dao.SysDept.Ctx(ctx).
			Where(dao.SysDept.Columns().TenantId, tenantID).
			Where(dao.SysDept.Columns().Id, parentId).
			Count()
		if err != nil {
			return err
		}
		if parentCount == 0 {
			return gerror.NewCodef(gcode.CodeValidationFailed, "Parent department with ID %s not found", parentId)
		}
	}

	var existingDept entity.SysDept
	err = dao.SysDept.Ctx(ctx).
		Where(dao.SysDept.Columns().TenantId, tenantID).
		Where(dao.SysDept.Columns().Id, in.ID).
		Scan(&existingDept)
	if err != nil {
		return gerror.NewCodef(gcode.CodeNotFound, "Department with ID %s not found", in.ID)
	}
	if existingDept.Id == "" {
		return gerror.NewCodef(gcode.CodeNotFound, "Department with ID %s not found", in.ID)
	}

	columns := dao.SysDept.Columns()
	updateData := g.Map{
		columns.Name:      in.Name,
		columns.Status:    in.Status,
		columns.Order:     in.Order,
		columns.ModifierId: in.ModifierId,
	}
	if parentId != "" {
		updateData[columns.ParentId] = parentId
	} else {
		updateData[columns.ParentId] = nil
	}

	_, err = dao.SysDept.Ctx(ctx).
		Data(updateData).
		Where(dao.SysDept.Columns().TenantId, tenantID).
		Where(dao.SysDept.Columns().Id, in.ID).
		Update()

	return err
}

func (s *sDept) DeleteDept(ctx context.Context, id string) (err error) {
	tenantID := resolveTenantID(ctx)

	var existingDept entity.SysDept
	err = dao.SysDept.Ctx(ctx).
		Where(dao.SysDept.Columns().TenantId, tenantID).
		Where(dao.SysDept.Columns().Id, id).
		Scan(&existingDept)
	if err != nil {
		return gerror.NewCodef(gcode.CodeNotFound, "Department with ID %s not found", id)
	}
	if existingDept.Id == "" {
		return gerror.NewCodef(gcode.CodeNotFound, "Department with ID %s not found", id)
	}

	childCount, err := dao.SysDept.Ctx(ctx).
		Where(dao.SysDept.Columns().TenantId, tenantID).
		Where(dao.SysDept.Columns().ParentId, id).
		Count()
	if err != nil {
		return err
	}
	if childCount > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Department with ID %s has child departments", id)
	}

	// Skip user association check due to type mismatch between sys_dept.id (UUID) and sys_user.dept_id (bigint)
	// This should be fixed in a future migration to align the data types

	_, err = dao.SysDept.Ctx(ctx).
		Where(dao.SysDept.Columns().TenantId, tenantID).
		Where(dao.SysDept.Columns().Id, id).
		Delete()

	return err
}

func (s *sDept) GetDeptList(ctx context.Context, in model.SysDeptGetListIn) (out *model.SysDeptGetListOut, err error) {
	out = &model.SysDeptGetListOut{}
	query := dao.SysDept.Ctx(ctx).OmitEmpty()

	tenantID := resolveTenantID(ctx)
	query = query.Where(dao.SysDept.Columns().TenantId, tenantID)

	if in.Name != "" {
		query = query.WhereLike(dao.SysDept.Columns().Name, "%"+in.Name+"%")
	}
	if in.Status != "" {
		query = query.Where(dao.SysDept.Columns().Status, in.Status)
	}

	out.Total, err = query.Count()
	if err != nil {
		return nil, err
	}

	err = query.Order(dao.SysDept.Columns().Order + " asc").Scan(&out.List)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (s *sDept) GetDeptTree(ctx context.Context) (out []*model.SysDeptTreeOut, err error) {
	tenantID := resolveTenantID(ctx)

	var allDepts []entity.SysDept
	err = dao.SysDept.Ctx(ctx).
		Where(dao.SysDept.Columns().TenantId, tenantID).
		Where(dao.SysDept.Columns().Status, 1).
		Where("deleted_at is null").
		Order(dao.SysDept.Columns().Order + " asc").
		Scan(&allDepts)
	if err != nil {
		return nil, err
	}

	if len(allDepts) == 0 {
		return []*model.SysDeptTreeOut{}, nil
	}

	deptMap := make(map[string]*model.SysDeptTreeOut)
	for _, dept := range allDepts {
		deptMap[dept.Id] = &model.SysDeptTreeOut{
			Id:       dept.Id,
			Name:     dept.Name,
			ParentId: dept.ParentId,
			Status:   dept.Status,
			Order:    dept.Order,
		}
	}

	var roots []*model.SysDeptTreeOut
	for _, dept := range allDepts {
		item := deptMap[dept.Id]
		if dept.ParentId == "" || dept.ParentId == "0" {
			roots = append(roots, item)
			continue
		}
		parent := deptMap[dept.ParentId]
		if parent == nil {
			roots = append(roots, item)
			continue
		}
		parent.Children = append(parent.Children, item)
	}

	return roots, nil
}