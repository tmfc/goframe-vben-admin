// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package dept

import (
	"context"

	"backend/api/dept/v1"
)

type IDeptV1 interface {
	CreateDept(ctx context.Context, req *v1.CreateDeptReq) (res *v1.CreateDeptRes, err error)
	GetDept(ctx context.Context, req *v1.GetDeptReq) (res *v1.GetDeptRes, err error)
	UpdateDept(ctx context.Context, req *v1.UpdateDeptReq) (res *v1.UpdateDeptRes, err error)
	DeleteDept(ctx context.Context, req *v1.DeleteDeptReq) (res *v1.DeleteDeptRes, err error)
	GetDeptList(ctx context.Context, req *v1.GetDeptListReq) (res *v1.GetDeptListRes, err error)
	GetDeptTree(ctx context.Context, req *v1.GetDeptTreeReq) (res *v1.GetDeptTreeRes, err error)
}