package model

import "github.com/gogf/gf/v2/os/gtime"

// UserCreateIn is the input for creating a user.
type UserCreateIn struct {
	Username string         `json:"username" v:"required"`
	Password string         `json:"password" v:"required"`
	RealName string         `json:"realName"`
	Status   int            `json:"status"`
	Roles    string         `json:"roles"`
	HomePath string         `json:"homePath"`
	Avatar   string         `json:"avatar"`
	DeptId   int64          `json:"deptId"`
	ExtInfo  map[string]any `json:"extInfo"`
}

// UserUpdateIn is the input for updating a user.
type UserUpdateIn struct {
	ID       int64          `json:"id" v:"required"`
	Username string         `json:"username" v:"required"`
	Password string         `json:"password"`
	RealName string         `json:"realName"`
	Status   int            `json:"status"`
	Roles    string         `json:"roles"`
	HomePath string         `json:"homePath"`
	Avatar   string         `json:"avatar"`
	DeptId   int64          `json:"deptId"`
	ExtInfo  map[string]any `json:"extInfo"`
}

// UserListIn is the input for listing users.
type UserListIn struct {
	Page     int    `json:"page"  d:"1"`
	PageSize int    `json:"pageSize" d:"10"`
	Username string `json:"username"`
}

// UserListItem is a lightweight user representation for list responses.
type UserListItem struct {
	Id        int64          `json:"id"`
	Username  string         `json:"username"`
	RealName  string         `json:"realName"`
	Status    int            `json:"status"`
	Roles     string         `json:"roles"`
	HomePath  string         `json:"homePath"`
	Avatar    string         `json:"avatar"`
	CreatedAt *gtime.Time    `json:"createdAt"`
	ExtInfo   map[string]any `json:"extInfo"`
	DeptId    int64          `json:"deptId"`
}

// UserListOut is the output for listing users.
type UserListOut struct {
	Items []*UserListItem `json:"items"`
	Total int             `json:"total"`
}
