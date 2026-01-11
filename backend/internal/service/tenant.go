package service

import (
	"context"

	"backend/internal/dao"
	"backend/internal/model"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var (
	localTenant ITenant
)

// Tenant returns the tenant service instance.
func Tenant() ITenant {
	return localTenant
}

// RegisterTenant sets the instance used by tenant related handlers.
func RegisterTenant(i ITenant) {
	localTenant = i
}

var _ ITenant = (*sTenant)(nil)

func init() {
	RegisterTenant(NewTenant())
}

// NewTenant creates a new tenant service instance.
func NewTenant() *sTenant {
	return &sTenant{}
}

// ITenant defines the tenant service interface.
type ITenant interface {
	List(ctx context.Context, in model.TenantListIn) (out *model.TenantListOut, err error)
	Get(ctx context.Context, id int64) (out *model.TenantListItem, err error)
	Create(ctx context.Context, in model.TenantCreateIn) (id int64, err error)
	Update(ctx context.Context, in model.TenantUpdateIn) error
	Delete(ctx context.Context, id int64) error
}

type sTenant struct{}

// List returns paginated tenants.
func (s *sTenant) List(ctx context.Context, in model.TenantListIn) (out *model.TenantListOut, err error) {
	out = &model.TenantListOut{}
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	m := dao.SysTenant.Ctx(ctx)

	if in.Name != "" {
		m = m.WhereLike(dao.SysTenant.Columns().Name, "%"+in.Name+"%")
	}
	if in.Code != "" {
		m = m.WhereLike(dao.SysTenant.Columns().Code, "%"+in.Code+"%")
	}

	out.Total, err = m.Count()
	if err != nil {
		return nil, err
	}

	if out.Total == 0 {
		return out, nil
	}

	err = m.Page(in.Page, in.PageSize).
		OrderDesc(dao.SysTenant.Columns().CreatedAt).
		Scan(&out.Items)
	
	return out, err
}

// Get returns a single tenant by id.
func (s *sTenant) Get(ctx context.Context, id int64) (out *model.TenantListItem, err error) {
	out = &model.TenantListItem{}
	err = dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, id).Scan(out)
	if err != nil {
		return nil, err
	}
	if out.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "tenant not found")
	}
	return out, nil
}

// Create creates a new tenant.
func (s *sTenant) Create(ctx context.Context, in model.TenantCreateIn) (id int64, err error) {
	if err = g.Validator().Data(in).Run(ctx); err != nil {
		return 0, err
	}

	// Check for unique code
	count, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Code, in.Code).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, gerror.NewCodef(gcode.CodeValidationFailed, "Tenant Code '%s' already exists", in.Code)
	}

	data := g.Map{
		dao.SysTenant.Columns().Name:           in.Name,
		dao.SysTenant.Columns().Code:           in.Code,
		dao.SysTenant.Columns().ContactName:    in.ContactName,
		dao.SysTenant.Columns().ContactMobile:  in.ContactMobile,
		dao.SysTenant.Columns().StartAt:        in.StartAt,
		dao.SysTenant.Columns().ExpireAt:       in.ExpireAt,
		dao.SysTenant.Columns().PackageVersion: in.PackageVersion,
		dao.SysTenant.Columns().Domain:         in.Domain,
		dao.SysTenant.Columns().Status:         in.Status,
	}

	// Insert
	result, err := dao.SysTenant.Ctx(ctx).Data(data).Insert()
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	return id, err
}

// Update updates an existing tenant.
func (s *sTenant) Update(ctx context.Context, in model.TenantUpdateIn) error {
	if err := g.Validator().Data(in).Run(ctx); err != nil {
		return err
	}

	// Check for unique code (excluding self)
	count, err := dao.SysTenant.Ctx(ctx).
		Where(dao.SysTenant.Columns().Code, in.Code).
		WhereNot(dao.SysTenant.Columns().Id, in.ID).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return gerror.NewCodef(gcode.CodeValidationFailed, "Tenant Code '%s' already exists", in.Code)
	}

	data := g.Map{
		dao.SysTenant.Columns().Name:           in.Name,
		dao.SysTenant.Columns().Code:           in.Code,
		dao.SysTenant.Columns().ContactName:    in.ContactName,
		dao.SysTenant.Columns().ContactMobile:  in.ContactMobile,
		dao.SysTenant.Columns().StartAt:        in.StartAt,
		dao.SysTenant.Columns().ExpireAt:       in.ExpireAt,
		dao.SysTenant.Columns().PackageVersion: in.PackageVersion,
		dao.SysTenant.Columns().Domain:         in.Domain,
		dao.SysTenant.Columns().Status:         in.Status,
	}

	_, err = dao.SysTenant.Ctx(ctx).Data(data).
		Where(dao.SysTenant.Columns().Id, in.ID).
		Update()
	return err
}

// Delete removes a tenant by id.
func (s *sTenant) Delete(ctx context.Context, id int64) error {
	_, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, id).Delete()
	return err
}
