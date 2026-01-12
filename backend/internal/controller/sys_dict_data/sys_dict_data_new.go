// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys_dict_data

import (
	"context"

	"backend/api/sys_dict_data"
	"backend/api/sys_dict_data/v1"
	"backend/internal/model"
	"backend/internal/service"
)

type ControllerV1 struct{}

func NewV1() sys_dict_data.ISysDictDataV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) CreateDictData(ctx context.Context, req *v1.CreateDictDataReq) (res *v1.CreateDictDataRes, err error) {
	if req == nil {
		req = &v1.CreateDictDataReq{}
	}
	id, err := service.SysDictData().CreateDictData(ctx, req.SysDictDataCreateIn)
	if err != nil {
		return nil, err
	}
	res = &v1.CreateDictDataRes{Id: id}
	return
}

func (c *ControllerV1) GetDictData(ctx context.Context, req *v1.GetDictDataReq) (res *v1.GetDictDataRes, err error) {
	if req == nil {
		req = &v1.GetDictDataReq{}
	}
	out, err := service.SysDictData().GetDictData(ctx, model.SysDictDataGetIn{Id: req.ID})
	if err != nil {
		return nil, err
	}
	res = &v1.GetDictDataRes{SysDictDataGetOut: out}
	return
}

func (c *ControllerV1) UpdateDictData(ctx context.Context, req *v1.UpdateDictDataReq) (res *v1.UpdateDictDataRes, err error) {
	if req == nil {
		req = &v1.UpdateDictDataReq{}
	}
	req.SysDictDataUpdateIn.Id = req.ID
	if err = service.SysDictData().UpdateDictData(ctx, req.SysDictDataUpdateIn); err != nil {
		return nil, err
	}
	res = &v1.UpdateDictDataRes{}
	return
}

func (c *ControllerV1) DeleteDictData(ctx context.Context, req *v1.DeleteDictDataReq) (res *v1.DeleteDictDataRes, err error) {
	if req == nil {
		req = &v1.DeleteDictDataReq{}
	}
	if err = service.SysDictData().DeleteDictData(ctx, model.SysDictDataDeleteIn{Id: req.ID}); err != nil {
		return nil, err
	}
	res = &v1.DeleteDictDataRes{}
	return
}

func (c *ControllerV1) GetDictDataList(ctx context.Context, req *v1.GetDictDataListReq) (res *v1.GetDictDataListRes, err error) {
	if req == nil {
		req = &v1.GetDictDataListReq{}
	}
	out, err := service.SysDictData().GetDictDataList(ctx, model.SysDictDataListIn{
		Page:       req.Page,
		PageSize:   req.PageSize,
		DictTypeId: req.DictTypeId,
		TypeCode:   req.TypeCode,
		Label:      req.Label,
		Value:      req.Value,
		Status:     req.Status,
	})
	if err != nil {
		return nil, err
	}
	res = &v1.GetDictDataListRes{Items: out.Items, Total: out.Total}
	return
}
