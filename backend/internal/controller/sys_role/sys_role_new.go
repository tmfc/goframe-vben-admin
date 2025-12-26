// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys_role

import (
	"context"

	"backend/api/sys_role"
	"backend/api/sys_role/v1"
	"backend/internal/service"
	"backend/internal/model"
)

type ControllerV1 struct{}

func NewV1() sys_role.ISysRoleV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) CreateRole(ctx context.Context, req *v1.CreateRoleReq) (res *v1.CreateRoleRes, err error) {
	if req == nil {
		req = &v1.CreateRoleReq{}
	}
	id, err := service.SysRole().CreateRole(ctx, req.SysRoleCreateIn)
	if err != nil {
		return nil, err
	}
	res = &v1.CreateRoleRes{Id: id}
	return
}

func (c *ControllerV1) GetRole(ctx context.Context, req *v1.GetRoleReq) (res *v1.GetRoleRes, err error) {
	if req == nil {
		req = &v1.GetRoleReq{}
	}
	out, err := service.SysRole().GetRole(ctx, model.SysRoleGetIn{Id: req.ID})
	if err != nil {
		return nil, err
	}
	res = &v1.GetRoleRes{SysRoleGetOut: out}
	return
}

func (c *ControllerV1) UpdateRole(ctx context.Context, req *v1.UpdateRoleReq) (res *v1.UpdateRoleRes, err error) {
	if req == nil {
		req = &v1.UpdateRoleReq{}
	}
	req.SysRoleUpdateIn.Id = req.ID
	err = service.SysRole().UpdateRole(ctx, req.SysRoleUpdateIn)
	if err != nil {
		return nil, err
	}
	res = &v1.UpdateRoleRes{}
	return
}

func (c *ControllerV1) DeleteRole(ctx context.Context, req *v1.DeleteRoleReq) (res *v1.DeleteRoleRes, err error) {
	if req == nil {
		req = &v1.DeleteRoleReq{}
	}
	err = service.SysRole().DeleteRole(ctx, model.SysRoleDeleteIn{Id: req.ID})
	if err != nil {
		return nil, err
	}
	res = &v1.DeleteRoleRes{}
	return
}