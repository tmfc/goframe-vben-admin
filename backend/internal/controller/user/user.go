package user

import (
	"context"
	"strings"

	backendConst "backend/internal/consts"
	"backend/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"backend/api/user/v1"
)

// ControllerV1 handles user-related endpoints.
type ControllerV1 struct{}

// Info returns the profile of the authenticated user.
func (c *ControllerV1) Info(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error) {
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
