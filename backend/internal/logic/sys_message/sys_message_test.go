package sys_message_test

import (
	"context"
	"testing"

	"backend/internal/model"
	"backend/internal/logic/sys_message"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_SendMessage(t *testing.T) {
	testutil.RequireDatabase(t)

	gtest.C(t, func(t *gtest.T) {
		s := sys_message.New()
		
		input := model.MessageCreateInput{
			Title:     "Test Message",
			Content:   "Hello",
			Type:      1,
			TenantID:  1,
			CreatorID: 1,
		}

		id, err := s.SendMessage(context.Background(), input)
		t.AssertNil(err)
		t.AssertGT(id, 0)
	})
}

func Test_GetList_And_SetRead(t *testing.T) {
	testutil.RequireDatabase(t)

	gtest.C(t, func(t *gtest.T) {
		s := sys_message.New()
		ctx := context.Background()
		userId := int64(1)

		// 1. Send a message to user
		input := model.MessageCreateInput{
			Title:      "Private Message",
			Type:       1,
			ReceiverID: &userId,
			TenantID:   1,
			CreatorID:  1,
		}
		id, err := s.SendMessage(ctx, input)
		t.AssertNil(err)

		// 2. Get List (Should be unread)
		// Assuming GetList signature: GetList(ctx context.Context, userId int64, page, size int) ([]*model.MessageItem, int, error)
		list, total, err := s.GetList(ctx, userId, 1, 10)
		t.AssertNil(err)
		t.AssertGT(total, 0)
		
		var foundItem *model.MessageItem
		for _, item := range list {
			if item.ID == id {
				foundItem = item
				break
			}
		}
		t.AssertNE(foundItem, nil)
		t.Assert(foundItem.IsRead, false)

		// 3. Set Read
		err = s.SetRead(ctx, userId, []int64{id})
		t.AssertNil(err)

		// 4. Get List (Should be read)
		list, _, err = s.GetList(ctx, userId, 1, 10)
		t.AssertNil(err)
		
		foundItem = nil
		for _, item := range list {
			if item.ID == id {
				foundItem = item
				break
			}
		}
		t.AssertNE(foundItem, nil)
		t.Assert(foundItem.IsRead, true)
	})
}
