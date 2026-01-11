package service

import (
	"context"
	"testing"

	"backend/api/auth/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"golang.org/x/crypto/bcrypt"
)

func TestUser_Info(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	testUsername := "testuser_info"
	testPassword := "testpassword"
	testTenantId := "1"

	// Clean up before test
	dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()

	// Ensure tenant exists
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
		// Login to get access token
		loginRes, err := Auth().Login(ctx, v1.LoginReq{
			Username: testUsername,
			Password: testPassword,
		})
		t.AssertNil(err)
		t.AssertNE(loginRes.AccessToken, "")

		// Test case 1: Get user info with valid token
		userInfo, err := User().Info(ctx, loginRes.AccessToken)
		t.AssertNil(err)
		t.AssertNE(userInfo, nil)
		t.Assert(userInfo.Username, testUsername)
		t.AssertNE(userInfo.UserId, "")
		t.AssertNE(userInfo.Roles, nil)
		t.AssertGT(len(userInfo.Roles), 0)
		t.Assert(userInfo.HomePath, "/dashboard")

		// Test case 2: Get user info with invalid token
		_, err = User().Info(ctx, "invalid.token.here")
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 3: Get user info with empty token
		_, err := User().Info(ctx, "")
		t.AssertNE(err, nil)
	})
}

func TestUser_Info_WithDefaultRole(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	testUsername := "testuser_default_role"
	testPassword := "testpassword"
	testTenantId := "1"

	// Clean up before test
	dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()

	// Ensure tenant exists
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

	// Create a user with empty roles to test default role assignment
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
		dao.SysUser.Columns().Username: testUsername,
		dao.SysUser.Columns().Password: string(hashedPassword),
		dao.SysUser.Columns().TenantId: testTenantId,
		dao.SysUser.Columns().Roles:    "[]", // Empty roles array
	}).Insert()
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Login to get access token
		loginRes, err := Auth().Login(ctx, v1.LoginReq{
			Username: testUsername,
			Password: testPassword,
		})
		t.AssertNil(err)
		t.AssertNE(loginRes.AccessToken, "")

		// Get user info and verify default role is assigned
		userInfo, err := User().Info(ctx, loginRes.AccessToken)
		t.AssertNil(err)
		t.AssertNE(userInfo, nil)
		t.AssertNE(userInfo.Roles, nil)
		t.AssertGT(len(userInfo.Roles), 0)
		// Verify default role is consts.RoleSuper
		t.Assert(userInfo.Roles[0], consts.DefaultRole())
	})
}

func TestUser_Info_WithCustomHomePath(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	testUsername := "testuser_homepath"
	testPassword := "testpassword"
	testTenantId := "1"
	customHomePath := "/custom/dashboard"

	// Clean up before test
	dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()

	// Ensure tenant exists
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

	// Create a user with custom home path
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
		dao.SysUser.Columns().Username: testUsername,
		dao.SysUser.Columns().Password: string(hashedPassword),
		dao.SysUser.Columns().TenantId: testTenantId,
		dao.SysUser.Columns().HomePath: customHomePath,
	}).Insert()
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	t.Cleanup(func() {
		dao.SysUser.Ctx(ctx).Unscoped().Where(dao.SysUser.Columns().Username, testUsername).Delete()
	})

	gtest.C(t, func(t *gtest.T) {
		// Login to get access token
		loginRes, err := Auth().Login(ctx, v1.LoginReq{
			Username: testUsername,
			Password: testPassword,
		})
		t.AssertNil(err)
		t.AssertNE(loginRes.AccessToken, "")

		// Get user info and verify custom home path is preserved
		userInfo, err := User().Info(ctx, loginRes.AccessToken)
		t.AssertNil(err)
		t.AssertNE(userInfo, nil)
		t.Assert(userInfo.HomePath, customHomePath)
	})
}