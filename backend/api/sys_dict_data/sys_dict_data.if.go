// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys_dict_data

import (
	"context"

	"backend/api/sys_dict_data/v1"
)

type ISysDictDataV1 interface {
	CreateDictData(ctx context.Context, req *v1.CreateDictDataReq) (res *v1.CreateDictDataRes, err error)
	GetDictData(ctx context.Context, req *v1.GetDictDataReq) (res *v1.GetDictDataRes, err error)
	UpdateDictData(ctx context.Context, req *v1.UpdateDictDataReq) (res *v1.UpdateDictDataRes, err error)
	DeleteDictData(ctx context.Context, req *v1.DeleteDictDataReq) (res *v1.DeleteDictDataRes, err error)
	GetDictDataList(ctx context.Context, req *v1.GetDictDataListReq) (res *v1.GetDictDataListRes, err error)
}
