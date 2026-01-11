package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// LoginReq defines the request structure for user login.
type LoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" summary:"User login" tags:"Authentication"`
	Username string `json:"username" v:"required#Username is required"`
	Password string `json:"password" v:"required#Password is required"`
}

// LoginRes defines the response structure for user login.
type LoginRes struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	UserInfo     struct { // Based on frontend's expected UserInfo
		ID       int64    `json:"id"`
		Username string   `json:"username"`
		RealName string   `json:"realName"`
		Avatar   string   `json:"avatar"`
		HomePath string   `json:"homePath"`
		Roles    []string `json:"roles"`
	} `json:"userInfo"`
}

// RefreshTokenReq defines the request structure for refreshing access token.
type RefreshTokenReq struct {
	g.Meta       `path:"/auth/refresh" method:"post" summary:"Refresh access token" tags:"Authentication"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshTokenRes defines the response structure for refreshing access token.
type RefreshTokenRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LogoutReq defines the request structure for user logout.
type LogoutReq struct {
	g.Meta       `path:"/auth/logout" method:"post" summary:"User logout" tags:"Authentication"`
	RefreshToken string `json:"refreshToken"`
}

// LogoutRes defines the response structure for user logout.
type LogoutRes struct{}

// GetAccessCodesReq defines the request structure for getting user access codes.
type GetAccessCodesReq struct {
	g.Meta `path:"/auth/codes" method:"get" summary:"Get user access codes" tags:"Authentication"`
	Token  string `json:"-"`
}

// GetAccessCodesRes defines the response structure for getting user access codes.
type GetAccessCodesRes struct {
	Codes []string `json:"codes"`
}
