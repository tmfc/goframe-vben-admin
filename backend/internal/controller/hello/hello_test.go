package hello

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// TestHelloController_HTTP verifies that the /hello endpoint returns the expected content.
func TestHelloController_HTTP(t *testing.T) {
	s := g.Server(guid.S())
	s.SetAddr(ghttp.FreePortAddress)

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(NewV1())
	})

	s.SetDumpRouterMap(false)
	s.Start()
	t.Cleanup(func() { s.Shutdown() })

	client := g.Client()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

	gtest.C(t, func(t *gtest.T) {
		content := client.GetContent(context.Background(), "/hello")
		t.AssertNE(content, "")
		t.AssertIN("Hello World!", content)
	})
}
