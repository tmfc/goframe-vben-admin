package sys_message

import (
	"context"
	"fmt"

	"backend/internal/dao"
	"backend/internal/logic/sys_message/mq"
	"backend/internal/model"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type sSysMessage struct{}

func New() *sSysMessage {
	return &sSysMessage{}
}

func (s *sSysMessage) SendMessage(ctx context.Context, in model.MessageCreateInput) (int64, error) {
	id, err := dao.SysMessage.Ctx(ctx).Data(in).InsertAndGetId()
	if err != nil {
		return 0, err
	}

	// Async push to Redis topic "sys_message"
	// We use a dedicated package 'mq' for this
	err = mq.Publish(ctx, "sys_message", in)
	if err != nil {
		// Log error but don't fail the database transaction
		g.Log().Error(ctx, "Failed to publish message to Redis", err)
	}

	return id, nil
}

func (s *sSysMessage) GetList(ctx context.Context, userId int64, page, size int) ([]*model.MessageItem, int, error) {
	var items []*model.MessageItem

	md := dao.SysMessage.Ctx(ctx).As("m")

	// Condition: receiver_id = userId OR receiver_id IS NULL
	md = md.Where("m.receiver_id = ? OR m.receiver_id IS NULL", userId)

	// Join with sys_message_read
	joinOn := fmt.Sprintf("m.id=r.message_id AND r.user_id=%d", userId)
	md = md.LeftJoin("sys_message_read r", joinOn)

	fields := "m.*, (CASE WHEN r.id IS NOT NULL THEN true ELSE false END) as is_read"

	total, err := md.Count()
	if err != nil {
		return nil, 0, err
	}

	err = md.Fields(fields).Page(page, size).Order("m.created_at DESC").Scan(&items)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (s *sSysMessage) SetRead(ctx context.Context, userId int64, messageIds []int64) error {
	if len(messageIds) == 0 {
		return nil
	}

	return dao.SysMessageRead.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, mid := range messageIds {
			count, err := dao.SysMessageRead.Ctx(ctx).
				Where("message_id=? AND user_id=?", mid, userId).Count()
			if err != nil {
				return err
			}
			if count == 0 {
				_, err := dao.SysMessageRead.Ctx(ctx).
					Data(g.Map{
						"message_id": mid,
						"user_id":    userId,
					}).Insert()
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *sSysMessage) SetAllRead(ctx context.Context, userId int64) error {
	// Find all message IDs that this user should see but hasn't read yet
	// INSERT INTO sys_message_read (message_id, user_id)
	// SELECT m.id, ? FROM sys_message m
	// LEFT JOIN sys_message_read r ON m.id = r.message_id AND r.user_id = ?
	// WHERE (m.receiver_id = ? OR m.receiver_id IS NULL) AND r.id IS NULL

	sql := `
		INSERT INTO sys_message_read (message_id, user_id)
		SELECT m.id, ? FROM sys_message m
		LEFT JOIN sys_message_read r ON m.id = r.message_id AND r.user_id = ?
		WHERE (m.receiver_id = ? OR m.receiver_id IS NULL) AND r.id IS NULL
	`
	_, err := dao.SysMessageRead.DB().Exec(ctx, sql, userId, userId, userId)
	return err
}