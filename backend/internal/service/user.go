package service

import (
	"context"

	"backend/api/user/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	localUser IUser
)

// User returns the user service instance.
func User() IUser {
	return localUser
}

// RegisterUser sets the instance used by user related handlers.
func RegisterUser(i IUser) {
	localUser = i
}

var _ IUser = (*sUser)(nil)

func init() {
	RegisterUser(NewUser())
}

// NewUser creates a new user service instance.
func NewUser() *sUser {
	return &sUser{}
}

// IUser defines the user service interface.
type IUser interface {
	Info(ctx context.Context, token string) (res *v1.UserInfoRes, err error)
}

type sUser struct{}

// Info returns the current authenticated user's profile by validating the JWT.
func (s *sUser) Info(ctx context.Context, token string) (res *v1.UserInfoRes, err error) {
	claims, err := parseToken(token)
	if err != nil {
		return nil, err
	}

	id, _ := claims["id"].(string)
	if id == "" {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "user id not found in token")
	}

	var user entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Id, id).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user.Id == "" {
		return nil, gerror.NewCode(consts.ErrorCodeUserNotFound, "user not found")
	}

	roles := parseRoles(user.Roles)
	if len(roles) == 0 {
		roles = []string{"super"}
	}

	homePath := user.HomePath
	if homePath == "" {
		homePath = "/dashboard"
	}

	res = &v1.UserInfoRes{
		UserId:   user.Id,
		Username: user.Username,
		RealName: user.RealName,
		Avatar:   user.Avatar,
		Roles:    roles,
		Desc:     user.RealName,
		HomePath: homePath,
		Token:    token,
	}
	return
}
