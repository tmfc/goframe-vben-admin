package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserInfoReq defines the request structure for fetching the current user data.
type UserInfoReq struct {
	g.Meta `path:"/user/info" method:"get" summary:"Fetch logged in user info" tags:"User"`
}

// UserInfoRes defines the response structure for the authenticated user's profile.
type UserInfoRes struct {
	UserId   string   `json:"userId"`
	Username string   `json:"username"`
	RealName string   `json:"realName"`
	Avatar   string   `json:"avatar"`
	Roles    []string `json:"roles"`
	Desc     string   `json:"desc"`
	HomePath string   `json:"homePath"`
	Token    string   `json:"token"`
}
