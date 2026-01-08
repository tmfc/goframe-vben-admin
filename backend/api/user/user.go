// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	v1 "backend/api/user/v1"
)

type IUserV1 interface {
	UserInfo(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error)
	List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error)
	Create(ctx context.Context, req *v1.CreateUserReq) (res *v1.CreateUserRes, err error)
	Update(ctx context.Context, req *v1.UpdateUserReq) (res *v1.UpdateUserRes, err error)
	Delete(ctx context.Context, req *v1.DeleteUserReq) (res *v1.DeleteUserRes, err error)
	Get(ctx context.Context, req *v1.GetUserReq) (res *v1.GetUserRes, err error)
}
