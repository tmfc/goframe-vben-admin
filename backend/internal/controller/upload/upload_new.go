// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package upload

import (
	"context"

	"backend/api/upload"
	"backend/api/upload/v1"
)

type ControllerV1 struct{}

func NewV1() upload.IUploadV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) Upload(ctx context.Context, req *v1.UploadReq) (res *v1.UploadRes, err error) {
	if req == nil {
		req = &v1.UploadReq{}
	}
	path, err := handleUpload(ctx, req.File)
	if err != nil {
		return nil, err
	}
	return &v1.UploadRes{Path: path}, nil
}
