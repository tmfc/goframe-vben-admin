// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package auth

import (
	"context"

	"backend/api/auth"
	"backend/api/auth/v1"
	"backend/internal/service"
)

type ControllerV1 struct{}

func NewV1() auth.IAuthV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	if req == nil {
		req = &v1.LoginReq{}
	}
	return service.Auth().Login(ctx, *req)
}

func (c *ControllerV1) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error) {
	if req == nil {
		req = &v1.RefreshTokenReq{}
	}
	return service.Auth().RefreshToken(ctx, *req)
}

func (c *ControllerV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	if req == nil {
		req = &v1.LogoutReq{}
	}
	return service.Auth().Logout(ctx, *req)
}

func (c *ControllerV1) GetAccessCodes(ctx context.Context, req *v1.GetAccessCodesReq) (res *v1.GetAccessCodesRes, err error) {
	if req == nil {
		req = &v1.GetAccessCodesReq{}
	}
	return service.Auth().GetAccessCodes(ctx, *req)
}
