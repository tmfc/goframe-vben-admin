package v1

import (
	"backend/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// UserInfoReq defines the request structure for fetching the current user data.
type UserInfoReq struct {
	g.Meta `path:"/user/info" method:"get" summary:"Fetch logged in user info" tags:"User"`
}

// UserInfoRes defines the response structure for the authenticated user's profile.
type UserInfoRes struct {
	UserId   int64    `json:"userId"`
	Username string   `json:"username"`
	RealName string   `json:"realName"`
	Avatar   string   `json:"avatar"`
	Roles    []string `json:"roles"`
	Desc     string   `json:"desc"`
	HomePath string   `json:"homePath"`
	Token    string   `json:"token"`
}

// UserListReq defines the request structure for listing users.
type UserListReq struct {
	g.Meta   `path:"/users" method:"get" summary:"List users" tags:"User"`
	Page     int    `json:"page"     d:"1"`
	PageSize int    `json:"pageSize" d:"10"`
	Username string `json:"username"`
}

// UserListRes defines the response structure for listing users.
type UserListRes struct {
	Items []*model.UserListItem `json:"items"`
	Total int                   `json:"total"`
}

// CreateUserReq defines the request structure for creating a user.
type CreateUserReq struct {
	g.Meta `path:"/users" method:"post" summary:"Create user" tags:"User"`
	model.UserCreateIn
}

// CreateUserRes defines the response structure for creating a user.
type CreateUserRes struct {
	Id int64 `json:"id"`
}

// UpdateUserReq defines the request structure for updating a user.
type UpdateUserReq struct {
	g.Meta `path:"/users/{id}" method:"put" summary:"Update user" tags:"User"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
	model.UserUpdateIn
}

// UpdateUserRes defines the response structure for updating a user.
type UpdateUserRes struct{}

// DeleteUserReq defines the request structure for deleting a user.
type DeleteUserReq struct {
	g.Meta `path:"/users/{id}" method:"delete" summary:"Delete user" tags:"User"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
}

// DeleteUserRes defines the response structure for deleting a user.
type DeleteUserRes struct{}

// GetUserReq defines the request structure for retrieving a user.
type GetUserReq struct {
	g.Meta `path:"/users/{id}" method:"get" summary:"Get user" tags:"User"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
}

// GetUserRes defines the response structure for retrieving a user.
type GetUserRes struct {
	*model.UserListItem
}
