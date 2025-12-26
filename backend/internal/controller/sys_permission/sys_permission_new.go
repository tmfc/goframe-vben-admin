// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys_permission

import (
	"context"

	"backend/api/sys_permission"
	"backend/api/sys_permission/v1"
	"backend/internal/model"
	"backend/internal/service"
)

type ControllerV1 struct{}

func NewV1() sys_permission.ISysPermissionV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) CreatePermission(ctx context.Context, req *v1.CreatePermissionReq) (res *v1.CreatePermissionRes, err error) {
	if req == nil {
		req = &v1.CreatePermissionReq{}
	}
	id, err := service.SysPermission().CreatePermission(ctx, req.SysPermissionCreateIn)
	if err != nil {
		return nil, err
	}
	res = &v1.CreatePermissionRes{Id: id}
	return
}

func (c *ControllerV1) GetPermission(ctx context.Context, req *v1.GetPermissionReq) (res *v1.GetPermissionRes, err error) {
	if req == nil {
		req = &v1.GetPermissionReq{}
	}
	out, err := service.SysPermission().GetPermission(ctx, model.SysPermissionGetIn{Id: req.ID})
	if err != nil {
		return nil, err
	}
	res = &v1.GetPermissionRes{SysPermissionGetOut: out}
	return
}

func (c *ControllerV1) UpdatePermission(ctx context.Context, req *v1.UpdatePermissionReq) (res *v1.UpdatePermissionRes, err error) {
	if req == nil {
		req = &v1.UpdatePermissionReq{}
	}
	req.SysPermissionUpdateIn.Id = req.ID
	err = service.SysPermission().UpdatePermission(ctx, req.SysPermissionUpdateIn)
	if err != nil {
		return nil, err
	}
	res = &v1.UpdatePermissionRes{}
	return
}

func (c *ControllerV1) DeletePermission(ctx context.Context, req *v1.DeletePermissionReq) (res *v1.DeletePermissionRes, err error) {
	if req == nil {
		req = &v1.DeletePermissionReq{}
	}
	err = service.SysPermission().DeletePermission(ctx, model.SysPermissionDeleteIn{Id: req.ID})
	if err != nil {
		return nil, err
	}
	res = &v1.DeletePermissionRes{}
	return
}
