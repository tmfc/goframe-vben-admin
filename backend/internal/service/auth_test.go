package service

import (
	"context"
	"testing"

	"backend/api/auth/v1"
	"backend/internal/consts"
	"backend/internal/dao"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestAuth_Login(t *testing.T) {
	ctx := context.TODO()
	testUsername := "testuser"
	testPassword := "testpassword"
	testTenantId := "00000000-0000-0000-0000-000000000000"

	// Clean up before test (hard delete to avoid soft-delete key conflicts).
	dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()
	dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().TenantId, testTenantId).Delete()
	dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Id, testTenantId).Delete()

	// Create a tenant for testing
	_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
		dao.SysTenant.Columns().Id:   testTenantId,
		dao.SysTenant.Columns().Name: "Test Tenant",
	}).Insert()
	if err != nil {
		t.Fatalf("failed to create tenant: %v", err)
	}

	// Create a user for testing
	if err := Auth().CreateUserForTest(ctx, testUsername, testPassword); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().TenantId, testTenantId).Delete()
		dao.SysTenant.Ctx(ctx).Unscoped().Where(dao.SysTenant.Columns().Id, testTenantId).Delete()
	})

	testCases := []struct {
		name        string
		req         v1.LoginReq
		wantErrCode gcode.Code
		wantUser    string
	}{
		{
			name: "Successful Login",
			req: v1.LoginReq{
				Username: testUsername,
				Password: testPassword,
			},
			wantUser: testUsername,
		},
		{
			name: "Incorrect Password",
			req: v1.LoginReq{
				Username: testUsername,
				Password: "wrongpassword",
			},
			wantErrCode: consts.ErrorCodeIncorrectPassword,
		},
		{
			name: "User Not Found",
			req: v1.LoginReq{
				Username: "nonexistentuser",
				Password: testPassword,
			},
			wantErrCode: consts.ErrorCodeUserNotFound,
		},
		{
			name: "Empty Username",
			req: v1.LoginReq{
				Username: "",
				Password: testPassword,
			},
			wantErrCode: consts.ErrorCodeUserNotFound,
		},
		{
			name: "Empty Password",
			req: v1.LoginReq{
				Username: testUsername,
				Password: "",
			},
			wantErrCode: consts.ErrorCodeIncorrectPassword,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				res, err := Auth().Login(ctx, tc.req)
				if tc.wantErrCode == nil {
					t.AssertNil(err)
					t.AssertNE(res.AccessToken, "")
					t.AssertNE(res.RefreshToken, "")
					t.Assert(res.UserInfo.Username, tc.wantUser)

					refreshRes, refreshErr := Auth().RefreshToken(ctx, v1.RefreshTokenReq{
						RefreshToken: res.RefreshToken,
					})
					t.AssertNil(refreshErr)
					t.AssertNE(refreshRes.AccessToken, "")
					t.AssertNE(refreshRes.RefreshToken, "")
					codesRes, codesErr := Auth().GetAccessCodes(ctx, v1.GetAccessCodesReq{
						Token: refreshRes.AccessToken,
					})
					t.AssertNil(codesErr)
					t.AssertGT(len(codesRes.Codes), 0)

					_, logoutErr := Auth().Logout(ctx, v1.LogoutReq{
						RefreshToken: refreshRes.RefreshToken,
					})
					t.AssertNil(logoutErr)
					return
				}

				t.AssertNE(err, nil)
				t.Assert(res, nil)
				t.Assert(gerror.Code(err), tc.wantErrCode)
			})
		})
	}
}
