// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"backend/internal/model"
	"context"
)

type (
	ISysMessage interface {
		SendMessage(ctx context.Context, in model.MessageCreateInput) (int64, error)
		GetList(ctx context.Context, userId int64, page int, size int) ([]*model.MessageItem, int, error)
		SetRead(ctx context.Context, userId int64, messageIds []int64) error
		SetAllRead(ctx context.Context, userId int64) error
	}
)

var (
	localSysMessage ISysMessage
)

func SysMessage() ISysMessage {
	if localSysMessage == nil {
		panic("implement not found for interface ISysMessage, forgot register?")
	}
	return localSysMessage
}

func RegisterSysMessage(i ISysMessage) {
	localSysMessage = i
}
