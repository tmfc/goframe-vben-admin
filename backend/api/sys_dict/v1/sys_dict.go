package v1

import (
	"backend/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// GetDictOptionsReq defines the request structure for retrieving dict options.
type GetDictOptionsReq struct {
	g.Meta   `path:"/sys-dict/options/{typeCode}" method:"get" summary:"Get dict options by type code" tags:"System Dict"`
	TypeCode string `json:"typeCode" v:"required#类型编码不能为空"`
}

// GetDictOptionsRes defines the response structure for retrieving dict options.
type GetDictOptionsRes struct {
	List []model.SysDictOption `json:"list"`
}

// GetDictLabelReq defines the request structure for retrieving dict label.
type GetDictLabelReq struct {
	g.Meta   `path:"/sys-dict/label/{typeCode}/{value}" method:"get" summary:"Get dict label by type code and value" tags:"System Dict"`
	TypeCode string `json:"typeCode" v:"required#类型编码不能为空"`
	Value    string `json:"value" v:"required#值不能为空"`
}

// GetDictLabelRes defines the response structure for retrieving dict label.
type GetDictLabelRes struct {
	Label string `json:"label"`
}
