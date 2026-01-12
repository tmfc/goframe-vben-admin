package sys_dict_type

import (
	"context"

	"backend/api/sys_dict_type/v1"
)

type ISysDictTypeV1 interface {
	CreateDictType(ctx context.Context, req *v1.CreateDictTypeReq) (res *v1.CreateDictTypeRes, err error)
	GetDictType(ctx context.Context, req *v1.GetDictTypeReq) (res *v1.GetDictTypeRes, err error)
	GetDictTypeList(ctx context.Context, req *v1.GetDictTypeListReq) (res *v1.GetDictTypeListRes, err error)
	UpdateDictType(ctx context.Context, req *v1.UpdateDictTypeReq) (res *v1.UpdateDictTypeRes, err error)
	DeleteDictType(ctx context.Context, req *v1.DeleteDictTypeReq) (res *v1.DeleteDictTypeRes, err error)
}
