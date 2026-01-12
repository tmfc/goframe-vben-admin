package tests

import (
	"testing"

	"backend/internal/model"
	"backend/internal/model/entity"
	"backend/internal/logic/sys_message"
	"backend/internal/service"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Message_E2E(t *testing.T) {
	testutil.RequireDatabase(t)
	ctx := gctx.New()

	gtest.C(t, func(t *gtest.T) {
		// 1. Prepare user and token
		// Assuming user ID 1 exists from seed
		userId := int64(1)
		token, err := service.NewJWTTokenService().GenerateAccessToken(&entity.SysUser{
			Id:       userId,
			Username: "vben",
			TenantId: 1,
		})
		t.AssertNil(err)
		_ = token // Token unused in direct service call test

		s := sys_message.New()
		
		// 2. Clear old messages for this user to have a clean state
		_, err = g.DB().Exec(ctx, "DELETE FROM sys_message WHERE receiver_id = ?", userId)
		t.AssertNil(err)

		// 3. Send a test message
		msgId, err := s.SendMessage(ctx, model.MessageCreateInput{
			Title:      "E2E Test Message",
			ReceiverID: &userId,
			TenantID:   1,
			CreatorID:  1,
		})
		t.AssertNil(err)

		// 4. Test API: List
		// We use g.Client() to simulate HTTP request if server is running, 
		// but here we can call service directly or use ghttp.Test
		// For simplicity in this environment, let's verify via service but through the controller's logic flow
		
		list, total, err := s.GetList(ctx, userId, 1, 10)
		t.AssertNil(err)
		t.Assert(total >= 1, true)
		t.Assert(list[0].Title, "E2E Test Message")
		t.Assert(list[0].IsRead, false)

		// 5. Test API: SetRead
		err = s.SetRead(ctx, userId, []int64{msgId})
		t.AssertNil(err)

		// 6. Verify Read
		list, _, err = s.GetList(ctx, userId, 1, 10)
		t.AssertNil(err)
		t.Assert(list[0].IsRead, true)
	})
}
