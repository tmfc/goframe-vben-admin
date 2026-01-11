package dept

import (
	"context"

	"backend/api/dept"
	"backend/api/dept/v1"
	"backend/internal/model"
	"backend/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type ControllerV1 struct{}

func NewV1() dept.IDeptV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) CreateDept(ctx context.Context, req *v1.CreateDeptReq) (res *v1.CreateDeptRes, err error) {
	res = &v1.CreateDeptRes{}
	id, err := service.Dept().CreateDept(ctx, req.SysDeptCreateIn)
	if err != nil {
		return nil, err
	}
	res.Id = id
	return res, nil
}

func (c *ControllerV1) GetDept(ctx context.Context, req *v1.GetDeptReq) (res *v1.GetDeptRes, err error) {
	if req.ID == 0 {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "ID不能为空")
	}
	res = &v1.GetDeptRes{}
	res.SysDeptGetOut, err = service.Dept().GetDept(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ControllerV1) UpdateDept(ctx context.Context, req *v1.UpdateDeptReq) (res *v1.UpdateDeptRes, err error) {
	res = &v1.UpdateDeptRes{}
	if req.ID == 0 {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "ID不能为空")
	}
	req.SysDeptUpdateIn.ID = req.ID
	if err = service.Dept().UpdateDept(ctx, req.SysDeptUpdateIn); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ControllerV1) DeleteDept(ctx context.Context, req *v1.DeleteDeptReq) (res *v1.DeleteDeptRes, err error) {
	res = &v1.DeleteDeptRes{}
	if req.ID == 0 {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "ID不能为空")
	}
	if err = service.Dept().DeleteDept(ctx, req.ID); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ControllerV1) GetDeptList(ctx context.Context, req *v1.GetDeptListReq) (res *v1.GetDeptListRes, err error) {
	res = &v1.GetDeptListRes{}
	out, err := service.Dept().GetDeptList(ctx, model.SysDeptGetListIn{
		Name:   req.Name,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	res.List = out.List
	res.Total = out.Total
	return res, nil
}

func (c *ControllerV1) GetDeptTree(ctx context.Context, req *v1.GetDeptTreeReq) (res *v1.GetDeptTreeRes, err error) {
	res = &v1.GetDeptTreeRes{}
	list, err := service.Dept().GetDeptTree(ctx)
	if err != nil {
		return nil, err
	}
	res.List = list
	return res, nil
}
