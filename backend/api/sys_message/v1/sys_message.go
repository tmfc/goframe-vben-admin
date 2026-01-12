package v1

import (
	"backend/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// GetListReq 获取我的消息列表
type GetListReq struct {
	g.Meta `path:"/message/list" method:"get" summary:"Get my message list" tags:"Message Center"`
	Page   int `json:"page" v:"min:1" d:"1"`
	Size   int `json:"size" v:"max:100" d:"10"`
}

// GetListRes 消息列表响应
type GetListRes struct {
	List  []*model.MessageItem `json:"list"`
	Total int                  `json:"total"`
}

// SetReadReq 标记已读
type SetReadReq struct {
	g.Meta `path:"/message/read" method:"post" summary:"Mark messages as read" tags:"Message Center"`
	IDs    []int64 `json:"ids" v:"required#Message IDs required"`
}

// SetReadRes 标记已读响应
type SetReadRes struct{}

// SetAllReadReq 全部标记已读
type SetAllReadReq struct {
	g.Meta `path:"/message/read-all" method:"post" summary:"Mark all messages as read" tags:"Message Center"`
}

// SetAllReadRes 全部标记已读响应
type SetAllReadRes struct{}
