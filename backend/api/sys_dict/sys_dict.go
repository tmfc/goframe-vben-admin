package sys_dict

import (
	"context"

	"backend/api/sys_dict/v1"
)

type ISysDictV1 interface {
	GetDictOptions(ctx context.Context, req *v1.GetDictOptionsReq) (res *v1.GetDictOptionsRes, err error)
	GetDictLabel(ctx context.Context, req *v1.GetDictLabelReq) (res *v1.GetDictLabelRes, err error)
}
