package user

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	v1 "backend/api/user/v1"
	"backend/internal/consts"
	"backend/internal/model"
	"backend/internal/service"
	"backend/internal/testutil"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

type apiEnvelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type mockUserService struct {
	InfoFunc   func(ctx context.Context, token string) (*v1.UserInfoRes, error)
	ListFunc   func(ctx context.Context, in model.UserListIn) (*model.UserListOut, error)
	GetFunc    func(ctx context.Context, id int64) (*model.UserListItem, error)
	CreateFunc func(ctx context.Context, in model.UserCreateIn) (int64, error)
	UpdateFunc func(ctx context.Context, in model.UserUpdateIn) error
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *mockUserService) Info(ctx context.Context, token string) (*v1.UserInfoRes, error) {
	return m.InfoFunc(ctx, token)
}

func (m *mockUserService) List(ctx context.Context, in model.UserListIn) (*model.UserListOut, error) {
	if m.ListFunc == nil {
		return &model.UserListOut{}, nil
	}
	return m.ListFunc(ctx, in)
}

func (m *mockUserService) Get(ctx context.Context, id int64) (*model.UserListItem, error) {
	if m.GetFunc == nil {
		return &model.UserListItem{Id: id}, nil
	}
	return m.GetFunc(ctx, id)
}

func (m *mockUserService) Create(ctx context.Context, in model.UserCreateIn) (int64, error) {
	if m.CreateFunc == nil {
		return 1, nil
	}
	return m.CreateFunc(ctx, in)
}

func (m *mockUserService) Update(ctx context.Context, in model.UserUpdateIn) error {
	if m.UpdateFunc == nil {
		return nil
	}
	return m.UpdateFunc(ctx, in)
}

func (m *mockUserService) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc == nil {
		return nil
	}
	return m.DeleteFunc(ctx, id)
}

// startUserAPIServer starts a temporary HTTP server binding the user controller.
func startUserAPIServer(t *testing.T) *ghttp.Server {
	t.Helper()

	s := g.Server(guid.S())
	s.SetAddr(ghttp.FreePortAddress)
	// Use default handler response envelope
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(NewV1())
	})

	s.SetDumpRouterMap(false)
	s.Start()
	t.Cleanup(func() { s.Shutdown() })
	return s
}

// decodeEnvelope is a small helper to decode the standard API envelope.
func decodeEnvelope(t *testing.T, content string) apiEnvelope {
	t.Helper()
	var env apiEnvelope
	if err := json.Unmarshal([]byte(content), &env); err != nil {
		t.Fatalf("failed to decode envelope: %v", err)
	}
	return env
}

// TestUserController_UserInfo_HTTP covers the /user/info endpoint with and without Authorization header.
func TestUserController_UserInfo_HTTP(t *testing.T) {
	testutil.RequireDatabase(t)

	// Mock user service so controller does not hit real DB logic.
	original := service.User()
	mockSvc := &mockUserService{
		InfoFunc: func(ctx context.Context, token string) (*v1.UserInfoRes, error) {
			return &v1.UserInfoRes{
				UserId:   1,
				Username: "test-user",
				RealName: "Test User",
				Roles:    []string{"admin"},
				HomePath: "/dashboard",
				Token:    token,
			}, nil
		},
	}
	service.RegisterUser(mockSvc)
	t.Cleanup(func() { service.RegisterUser(original) })

	s := startUserAPIServer(t)

	client := g.Client()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

	// Case 1: Missing Authorization header should return unauthorized error.
	gtest.C(t, func(t *gtest.T) {
		content := client.GetContent(context.Background(), "/user/info")
		env := decodeEnvelope(t, content)
		t.Assert(env.Code, consts.ErrorCodeUnauthorized.Code())
	})

	// Case 2: With Authorization header should succeed and return mocked user info.
	gtest.C(t, func(t *gtest.T) {
		client.SetHeader("Authorization", "Bearer mock-token")
		defer client.SetHeader("Authorization", "")

		content := client.GetContent(context.Background(), "/user/info")
		env := decodeEnvelope(t, content)
		t.Assert(env.Code, gcode.CodeOK.Code())

		var res v1.UserInfoRes
		t.AssertNil(json.Unmarshal(env.Data, &res))
		t.Assert(res.Username, "test-user")
		t.Assert(res.UserId, 1)
		t.AssertIN("mock-token", res.Token)
	})
}

// TestUserController_CRUD covers the basic CRUD methods by invoking controller directly with mock service.
func TestUserController_CRUD(t *testing.T) {
	original := service.User()
	mockSvc := &mockUserService{}
	service.RegisterUser(mockSvc)
	t.Cleanup(func() { service.RegisterUser(original) })

	ctrl := &ControllerV1{}
	ctx := context.TODO()

	gtest.C(t, func(t *gtest.T) {
		// List
		_, err := ctrl.UserList(ctx, &v1.UserListReq{Page: 1, PageSize: 10, Username: "foo"})
		t.AssertNil(err)

		// Get
		_, err = ctrl.GetUser(ctx, &v1.GetUserReq{ID: 1})
		t.AssertNil(err)

		// Create
		createReq := &v1.CreateUserReq{UserCreateIn: model.UserCreateIn{Username: "foo", Password: "bar"}}
		_, err = ctrl.CreateUser(ctx, createReq)
		t.AssertNil(err)

		// Update
		updateReq := &v1.UpdateUserReq{ID: 1, UserUpdateIn: model.UserUpdateIn{Username: "foo"}}
		_, err = ctrl.UpdateUser(ctx, updateReq)
		t.AssertNil(err)

		// Delete
		_, err = ctrl.DeleteUser(ctx, &v1.DeleteUserReq{ID: 1})
		t.AssertNil(err)
	})
}
