// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys_dict_type

import (
	"context"

	"backend/api/sys_dict_type"
	"backend/api/sys_dict_type/v1"
	"backend/internal/model"
	"backend/internal/service"
)

type ControllerV1 struct{}

func NewV1() sys_dict_type.ISysDictTypeV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) CreateDictType(ctx context.Context, req *v1.CreateDictTypeReq) (res *v1.CreateDictTypeRes, err error) {
	if req == nil {
		req = &v1.CreateDictTypeReq{}
	}
	id, err := service.SysDictType().CreateDictType(ctx, req.SysDictTypeCreateIn)
	if err != nil {
		return nil, err
	}
	res = &v1.CreateDictTypeRes{Id: id}
	return
}

func (c *ControllerV1) GetDictType(ctx context.Context, req *v1.GetDictTypeReq) (res *v1.GetDictTypeRes, err error) {
	if req == nil {
		req = &v1.GetDictTypeReq{}
	}
	out, err := service.SysDictType().GetDictType(ctx, model.SysDictTypeGetIn{Id: req.ID})
	if err != nil {
		return nil, err
	}
	res = &v1.GetDictTypeRes{SysDictTypeGetOut: out}
	return
}

func (c *ControllerV1) UpdateDictType(ctx context.Context, req *v1.UpdateDictTypeReq) (res *v1.UpdateDictTypeRes, err error) {
	if req == nil {
		req = &v1.UpdateDictTypeReq{}
	}
	req.SysDictTypeUpdateIn.Id = req.ID
	if err = service.SysDictType().UpdateDictType(ctx, req.SysDictTypeUpdateIn); err != nil {
		return nil, err
	}
	res = &v1.UpdateDictTypeRes{}
	return
}

func (c *ControllerV1) DeleteDictType(ctx context.Context, req *v1.DeleteDictTypeReq) (res *v1.DeleteDictTypeRes, err error) {
	if req == nil {
		req = &v1.DeleteDictTypeReq{}
	}
	if err = service.SysDictType().DeleteDictType(ctx, model.SysDictTypeDeleteIn{Id: req.ID}); err != nil {
		return nil, err
	}
	res = &v1.DeleteDictTypeRes{}
	return
}

func (c *ControllerV1) GetDictTypeList(ctx context.Context, req *v1.GetDictTypeListReq) (res *v1.GetDictTypeListRes, err error) {
	if req == nil {
		req = &v1.GetDictTypeListReq{}
	}
	out, err := service.SysDictType().GetDictTypeList(ctx, model.SysDictTypeListIn{
		Page:     req.Page,
		PageSize: req.PageSize,
		TypeCode: req.TypeCode,
		TypeName: req.TypeName,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.GetDictTypeListRes{Items: out.Items, Total: out.Total}
	return
}
