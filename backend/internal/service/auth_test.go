package service

import (
	"context"
	"testing"

	"backend/api/auth/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestAuth_Login(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	testUsername := "testuser"
	testPassword := "testpassword"
	testTenantId := "1"

	// Clean up before test (hard delete to avoid soft-delete key conflicts).
	dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()

	// Create a tenant for testing
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	// Create a user for testing
	if err := Auth().CreateUserForTest(ctx, testUsername, testPassword); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()
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

func TestAuth_RefreshToken(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	testUsername := "testuser_refresh"
	testPassword := "testpassword_refresh"
	testTenantId := "1"

	// Clean up before test
	dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()

	// Create a tenant for testing
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	// Create a user for testing
	if err := Auth().CreateUserForTest(ctx, testUsername, testPassword); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// First login to get tokens
		loginRes, err := Auth().Login(ctx, v1.LoginReq{
			Username: testUsername,
			Password: testPassword,
		})
		t.AssertNil(err)
		t.AssertNE(loginRes.RefreshToken, "")

		// Test successful refresh
		refreshRes, err := Auth().RefreshToken(ctx, v1.RefreshTokenReq{
			RefreshToken: loginRes.RefreshToken,
		})
		t.AssertNil(err)
		t.AssertNE(refreshRes.AccessToken, "")
		t.AssertNE(refreshRes.RefreshToken, "")

		// Test refresh with invalid token
		_, err = Auth().RefreshToken(ctx, v1.RefreshTokenReq{
			RefreshToken: "invalid_token",
		})
		t.AssertNE(err, nil)

		// Test refresh with empty token
		_, err = Auth().RefreshToken(ctx, v1.RefreshTokenReq{
			RefreshToken: "",
		})
		t.AssertNE(err, nil)
	})
}

func TestAuth_Logout(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	testUsername := "testuser_logout"
	testPassword := "testpassword_logout"
	testTenantId := "1"

	// Clean up before test
	dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()

	// Create a tenant for testing
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	// Create a user for testing
	if err := Auth().CreateUserForTest(ctx, testUsername, testPassword); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// First login to get tokens
		loginRes, err := Auth().Login(ctx, v1.LoginReq{
			Username: testUsername,
			Password: testPassword,
		})
		t.AssertNil(err)
		t.AssertNE(loginRes.RefreshToken, "")

		// Test successful logout
		logoutRes, err := Auth().Logout(ctx, v1.LogoutReq{
			RefreshToken: loginRes.RefreshToken,
		})
		t.AssertNil(err)
		t.AssertNE(logoutRes, nil)

		// Test logout with empty token (should still succeed)
		logoutRes, err = Auth().Logout(ctx, v1.LogoutReq{
			RefreshToken: "",
		})
		t.AssertNil(err)
		t.AssertNE(logoutRes, nil)
	})
}

func TestAuth_GetAccessCodes(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	testUsername := "testuser_codes"
	testPassword := "testpassword_codes"
	testTenantId := "1"

	// Clean up before test
	dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()

	// Create a tenant for testing
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	// Create a user for testing
	if err := Auth().CreateUserForTest(ctx, testUsername, testPassword); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// First login to get access token
		loginRes, err := Auth().Login(ctx, v1.LoginReq{
			Username: testUsername,
			Password: testPassword,
		})
		t.AssertNil(err)
		t.AssertNE(loginRes.AccessToken, "")

		// Test successful get access codes
		codesRes, err := Auth().GetAccessCodes(ctx, v1.GetAccessCodesReq{
			Token: loginRes.AccessToken,
		})
		t.AssertNil(err)
		t.AssertNE(codesRes, nil)
		t.AssertGT(len(codesRes.Codes), 0)

		// Test with invalid token
		_, err = Auth().GetAccessCodes(ctx, v1.GetAccessCodesReq{
			Token: "invalid_token",
		})
		t.AssertNE(err, nil)

		// Test with empty token
		_, err = Auth().GetAccessCodes(ctx, v1.GetAccessCodesReq{
			Token: "",
		})
		t.AssertNE(err, nil)
	})
}

func TestAuth_CreateUserForTest(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	testUsername := "testuser_create"
	testPassword := "testpassword_create"
	testTenantId := "1"

	// Clean up before test
	dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()

	// Create a tenant for testing
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, testTenantId).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   testTenantId,
			dao.SysTenant.Columns().Name: "Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}

	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Test successful user creation
		err := Auth().CreateUserForTest(ctx, testUsername, testPassword)
		t.AssertNil(err)

		// Verify user was created
		var user *entity.SysUser
		err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Username, testUsername).Scan(&user)
		t.AssertNil(err)
		t.AssertNE(user, nil)
		t.Assert(user.Username, testUsername)
	})
}