// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys_dict

import (
	"context"

	"backend/api/sys_dict"
	"backend/api/sys_dict/v1"
	"backend/internal/service"
)

type ControllerV1 struct{}

func NewV1() sys_dict.ISysDictV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) GetDictOptions(ctx context.Context, req *v1.GetDictOptionsReq) (res *v1.GetDictOptionsRes, err error) {
	if req == nil {
		req = &v1.GetDictOptionsReq{}
	}
	list, err := service.SysDict().GetDictOptions(ctx, req.TypeCode)
	if err != nil {
		return nil, err
	}
	res = &v1.GetDictOptionsRes{List: list}
	return
}

func (c *ControllerV1) GetDictLabel(ctx context.Context, req *v1.GetDictLabelReq) (res *v1.GetDictLabelRes, err error) {
	if req == nil {
		req = &v1.GetDictLabelReq{}
	}
	label, err := service.SysDict().GetDictLabel(ctx, req.TypeCode, req.Value)
	if err != nil {
		return nil, err
	}
	res = &v1.GetDictLabelRes{Label: label}
	return
}
