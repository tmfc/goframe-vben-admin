package model

import (
	"backend/internal/model/entity"
)

// SysDeptCreateIn is the input for creating a new department.
type SysDeptCreateIn struct {
	Name      string `json:"name" v:"required#名称不能为空"`
	ParentId  string `json:"parentId"`
	Status    int    `json:"status"`
	Order     int    `json:"order"`
	CreatorId int64  `json:"creatorId"`
}

// SysDeptGetOut is the output for retrieving a department.
type SysDeptGetOut struct {
	*entity.SysDept
}

// SysDeptUpdateIn is the input for updating a department.
type SysDeptUpdateIn struct {
	ID        string `json:"id" v:"required#ID不能为空"`
	Name      string `json:"name" v:"required#名称不能为空"`
	ParentId  string `json:"parentId"`
	Status    int    `json:"status"`
	Order     int    `json:"order"`
	ModifierId int64  `json:"modifierId"`
}

type SysDeptGetListIn struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SysDeptGetListOut struct {
	List  []*SysDeptGetOut `json:"list"`
	Total int              `json:"total"`
}

type SysDeptTreeOut struct {
	Id       string          `json:"id"`
	Name     string          `json:"name"`
	ParentId string          `json:"parentId"`
	Status   int             `json:"status"`
	Order    int             `json:"order"`
	Children []*SysDeptTreeOut `json:"children,omitempty"`
}