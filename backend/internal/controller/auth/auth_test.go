package auth

import (
	"context"
	"testing"

	"backend/api/auth/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/service"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestAuthController_Login(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	testUsername := "ctrltestuser"
	testPassword := "ctrltestpassword"
	testTenantId := "1"

	// Clean up before test.
	dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()

	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Controller Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	if err := service.Auth().CreateUserForTest(ctx, testUsername, testPassword); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()
	})

	ctrl := NewV1()

	gtest.C(t, func(t *gtest.T) {
		res, err := ctrl.Login(ctx, &v1.LoginReq{
			Username: testUsername,
			Password: testPassword,
		})
		t.AssertNil(err)
		t.AssertNE(res.AccessToken, "")
		t.Assert(res.UserInfo.Username, testUsername)
	})

	gtest.C(t, func(t *gtest.T) {
		res, err := ctrl.Login(ctx, &v1.LoginReq{
			Username: testUsername,
			Password: "wrongpassword",
		})
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), consts.ErrorCodeIncorrectPassword)
	})

	gtest.C(t, func(t *gtest.T) {
		res, err := ctrl.Login(ctx, nil)
		t.AssertNil(res)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), consts.ErrorCodeUserNotFound)
	})
}
