// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"backend/api/user/v1"
)

// IUserV1 defines the user controller interface.
type IUserV1 interface {
	Info(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error)
}
