package user

import (
	"context"
	"strings"

	backendConst "backend/internal/consts"
	"backend/internal/model"
	"backend/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "backend/api/user/v1"
)

// ControllerV1 handles user-related endpoints.
type ControllerV1 struct{}

// Info returns the profile of the authenticated user.
func (c *ControllerV1) UserInfo(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error) {
	if req == nil {
		req = &v1.UserInfoReq{}
	}
	gReq := g.RequestFromCtx(ctx)
	authHeader := gReq.Header.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	if token == "" {
		return nil, gerror.NewCode(backendConst.ErrorCodeUnauthorized, "missing authorization token")
	}
	return service.User().Info(ctx, token)
}

// List returns paginated users.
func (c *ControllerV1) List(ctx context.Context, req *v1.UserListReq) (res *v1.UserListRes, err error) {
	out, err := service.User().List(ctx, model.UserListIn{
		Page:     req.Page,
		PageSize: req.PageSize,
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UserListRes{
		Items: out.Items,
		Total: out.Total,
	}, nil
}

// Get returns a single user by id.
func (c *ControllerV1) Get(ctx context.Context, req *v1.GetUserReq) (res *v1.GetUserRes, err error) {
	out, err := service.User().Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserRes{UserListItem: out}, nil
}

// Create creates a new user.
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateUserReq) (res *v1.CreateUserRes, err error) {
	id, err := service.User().Create(ctx, req.UserCreateIn)
	if err != nil {
		return nil, err
	}
	return &v1.CreateUserRes{Id: id}, nil
}

// Update updates an existing user.
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateUserReq) (res *v1.UpdateUserRes, err error) {
	in := req.UserUpdateIn
	in.ID = req.ID
	if err := service.User().Update(ctx, in); err != nil {
		return nil, err
	}
	return &v1.UpdateUserRes{}, nil
}

// Delete removes a user by id.
func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteUserReq) (res *v1.DeleteUserRes, err error) {
	if err := service.User().Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return &v1.DeleteUserRes{}, nil
}
