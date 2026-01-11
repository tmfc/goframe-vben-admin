package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"backend/internal/model"
)

type CreateDeptReq struct {
	g.Meta `path:"/sys-dept" method:"post" summary:"Create a new department" tags:"System Department"`
	model.SysDeptCreateIn
}

type CreateDeptRes struct {
	Id int64 `json:"id"`
}

type GetDeptReq struct {
	g.Meta `path:"/sys-dept/{id}" method:"get" summary:"Retrieve a department by ID" tags:"System Department"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
}

type GetDeptRes struct {
	*model.SysDeptGetOut
}

type UpdateDeptReq struct {
	g.Meta `path:"/sys-dept/{id}" method:"put" summary:"Update a department by ID" tags:"System Department"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
	model.SysDeptUpdateIn
}

type UpdateDeptRes struct{}

type DeleteDeptReq struct {
	g.Meta `path:"/sys-dept/{id}" method:"delete" summary:"Delete a department by ID" tags:"System Department"`
	ID     int64 `json:"id" v:"required#ID不能为空"`
}

type DeleteDeptRes struct{}

type GetDeptListReq struct {
	g.Meta `path:"/sys-dept/list" method:"get" summary:"Retrieve department list" tags:"System Department"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type GetDeptListRes struct {
	List  []*model.SysDeptGetOut `json:"list"`
	Total int                    `json:"total"`
}

type GetDeptTreeReq struct {
	g.Meta `path:"/sys-dept/tree" method:"get" summary:"Retrieve department tree" tags:"System Department"`
}

type GetDeptTreeRes struct {
	List []*model.SysDeptTreeOut `json:"list"`
}
