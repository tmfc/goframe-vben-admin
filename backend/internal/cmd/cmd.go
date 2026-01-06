package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"backend/internal/controller/auth"
	"backend/internal/controller/dept"
	"backend/internal/controller/hello"
	"backend/internal/controller/menu"
	"backend/internal/controller/sys_permission"
	"backend/internal/controller/sys_role"
	"backend/internal/controller/user"
	"backend/internal/middleware"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			cfg := g.Cfg()
			port := 5666
			if portVar, err := cfg.Get(ctx, "server.port"); err == nil && portVar != nil {
				if p := portVar.Int(); p > 0 {
					port = p
				}
			}
			s.SetPort(port)
			if addrVar, err := cfg.Get(ctx, "server.address"); err == nil && addrVar != nil {
				if addr := strings.TrimSpace(addrVar.String()); addr != "" {
					if strings.Contains(addr, ":") {
						s.SetAddr(addr)
					} else {
						s.SetAddr(fmt.Sprintf("%s:%d", addr, port))
					}
				}
			}
			// Enable OpenAPI & Swagger UI for local development.
			s.SetOpenApiPath("/api.json")
			s.SetSwaggerPath("/swagger")
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
			s.Run()
			return nil
		},
	}
)
