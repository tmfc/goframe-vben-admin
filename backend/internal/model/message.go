package model

import "github.com/gogf/gf/v2/os/gtime"

// MessageCreateInput Input for creating a message
type MessageCreateInput struct {
	Title      string                 `json:"title"`
	Content    string                 `json:"content"`
	Type       int                    `json:"type"`
	ReceiverID *int64                 `json:"receiverId"` // Nullable
	Path       string                 `json:"path"`
	Params     map[string]interface{} `json:"params"`
	TenantID   int64                  `json:"tenantId"`  // Context injected
	CreatorID  int64                  `json:"creatorId"` // Context injected
}

// MessageItem Output item for list
type MessageItem struct {
	ID         int64                  `json:"id"`
	Title      string                 `json:"title"`
	Content    string                 `json:"content"`
	Type       int                    `json:"type"`
	ReceiverID *int64                 `json:"receiverId"`
	Path       string                 `json:"path"`
	Params     map[string]interface{} `json:"params"`
	Status     int                    `json:"status"`
	CreatedAt  *gtime.Time            `json:"createdAt"`
	IsRead     bool                   `json:"isRead"` // Computed
}
