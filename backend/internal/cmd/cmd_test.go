package cmd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"

	"backend/internal/controller/auth"
	"backend/internal/controller/dept"
	"backend/internal/controller/hello"
	"backend/internal/controller/menu"
	"backend/internal/controller/sys_permission"
	"backend/internal/controller/sys_role"
	"backend/internal/controller/user"
	"backend/internal/middleware"
)

func TestMenuRoutesRegistered(t *testing.T) {
	s := g.Server(guid.S())
	s.SetAddr(ghttp.FreePortAddress)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse, middleware.CasbinAuthz())
		group.Bind(
			hello.NewV1(),
			auth.NewV1(),
			dept.NewV1(),
			sys_permission.NewV1(),
			sys_role.NewV1(),
			menu.NewV1(),
			user.NewV1(),
		)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	time.Sleep(100 * time.Millisecond)
	defer s.Shutdown()

	gtest.C(t, func(t *gtest.T) {
		// Check for a few key menu routes
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		
		resp, err := client.Get(context.Background(), "/menu/all")
		t.AssertNil(err)
		defer resp.Close()
		t.AssertNE(resp.StatusCode, 404)

		resp, err = client.Post(context.Background(), "/sys-menu", g.Map{})
		t.AssertNil(err)
		defer resp.Close()
		t.AssertNE(resp.StatusCode, 404)
		
		resp, err = client.Get(context.Background(), "/sys-menu/1")
		t.AssertNil(err)
		defer resp.Close()
		t.AssertNE(resp.StatusCode, 404)
	})
}
