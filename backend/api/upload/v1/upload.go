package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// UploadReq defines the request structure for file uploads.
type UploadReq struct {
	g.Meta `path:"/api/v1/upload" method:"post" summary:"Upload file" tags:"Upload" mime:"multipart/form-data"`
	File   *ghttp.UploadFile `json:"file" type:"file" v:"required#Upload file is required"`
}

// UploadRes defines the response structure for file uploads.
type UploadRes struct {
	Path string `json:"path"`
}
