package menu

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	v1 "backend/api/menu/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

type menuAPIEnvelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func TestMenuAPIEndpoints(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	ensureTestTenant(t, ctx, "00000000-0000-0000-0000-000000000000")

	t.Cleanup(func() {
		dao.SysMenu.Ctx(ctx).Unscoped().WhereLike(dao.SysMenu.Columns().Name, "ApiMenu%").Delete()
	})

	s := startMenuAPIServer(t)
	client := g.Client().ContentJson()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

	gtest.C(t, func(t *gtest.T) {
		allContent := client.GetContent(ctx, "/menu/all")
		allEnv := decodeMenuEnvelope(t, allContent)
		t.Assert(allEnv.Code, gcode.CodeOK.Code())

		var allMenus v1.MenuAllRes
		t.AssertNil(json.Unmarshal(allEnv.Data, &allMenus))
		t.AssertGT(len(allMenus), 0)

		createContent := client.PostContent(ctx, "/sys-menu", `{"name":"ApiMenuRoot","path":"/api-menu-root","component":"/api/menu/root","type":"menu","status":1,"order":10}`)
		createEnv := decodeMenuEnvelope(t, createContent)
		t.Assert(createEnv.Code, gcode.CodeOK.Code())

		var createRes v1.CreateMenuRes
		t.AssertNil(json.Unmarshal(createEnv.Data, &createRes))
		t.AssertNE(createRes.Id, "")

		getContent := client.GetContent(ctx, fmt.Sprintf("/sys-menu/%s", createRes.Id))
		getEnv := decodeMenuEnvelope(t, getContent)
		t.Assert(getEnv.Code, gcode.CodeOK.Code())

		var getRes struct {
			Name string `json:"name"`
		}
		t.AssertNil(json.Unmarshal(getEnv.Data, &getRes))
		t.Assert(getRes.Name, "ApiMenuRoot")

		updatePayload := fmt.Sprintf(`{"id":"%s","name":"ApiMenuUpdated","path":"/api-menu-updated","component":"/api/menu/updated","type":"menu","status":1,"order":20}`, createRes.Id)
		updateContent := client.PutContent(ctx, fmt.Sprintf("/sys-menu/%s", createRes.Id), updatePayload)
		updateEnv := decodeMenuEnvelope(t, updateContent)
		t.Assert(updateEnv.Code, gcode.CodeOK.Code())

		listContent := client.GetContent(ctx, "/sys-menu/list?name=ApiMenuUpdated")
		listEnv := decodeMenuEnvelope(t, listContent)
		t.Assert(listEnv.Code, gcode.CodeOK.Code())

		var listRes struct {
			List []struct {
				Name string `json:"name"`
			} `json:"list"`
		}
		t.AssertNil(json.Unmarshal(listEnv.Data, &listRes))
		t.Assert(hasMenuName(listRes.List, "ApiMenuUpdated"), true)

		deleteContent := client.DeleteContent(ctx, fmt.Sprintf("/sys-menu/%s", createRes.Id))
		deleteEnv := decodeMenuEnvelope(t, deleteContent)
		t.Assert(deleteEnv.Code, gcode.CodeOK.Code())

		missingContent := client.GetContent(ctx, fmt.Sprintf("/sys-menu/%s", createRes.Id))
		missingEnv := decodeMenuEnvelope(t, missingContent)
		t.Assert(missingEnv.Code, gcode.CodeNotFound.Code())
	})
}

func startMenuAPIServer(t *testing.T) *ghttp.Server {
	t.Helper()
	s := g.Server(guid.S())
	s.SetAddr(ghttp.FreePortAddress)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(NewV1())
	})
	s.SetDumpRouterMap(false)
	s.Start()
	time.Sleep(100 * time.Millisecond)
	t.Cleanup(func() { s.Shutdown() })
	return s
}

func decodeMenuEnvelope(t *gtest.T, content string) menuAPIEnvelope {
	t.Helper()
	var env menuAPIEnvelope
	t.AssertNil(json.Unmarshal([]byte(content), &env))
	return env
}

func ensureTestTenant(t *testing.T, ctx context.Context, tenantID string) {
	t.Helper()
	existing, err := dao.SysTenant.Ctx(ctx).Where(dao.SysTenant.Columns().Id, tenantID).One()
	if err != nil {
		t.Fatalf("failed to query tenant: %v", err)
	}
	if existing.IsEmpty() {
		_, err := dao.SysTenant.Ctx(ctx).Data(g.Map{
			dao.SysTenant.Columns().Id:   tenantID,
			dao.SysTenant.Columns().Name: "API Test Tenant",
		}).Insert()
		if err != nil {
			t.Fatalf("failed to create tenant: %v", err)
		}
	}
}

func hasMenuName(items []struct {
	Name string `json:"name"`
}, name string) bool {
	for _, item := range items {
		if item.Name == name {
			return true
		}
	}
	return false
}

func TestMenuAPIValidation(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	ensureTestTenant(t, ctx, "00000000-0000-0000-0000-000000000000")

	s := startMenuAPIServer(t)
	client := g.Client().ContentJson()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

	gtest.C(t, func(t *gtest.T) {
		// Test creating a menu with a missing name
		createContent := client.PostContent(ctx, "/sys-menu", `{"path":"/api-menu-validation","type":"menu"}`)
		createEnv := decodeMenuEnvelope(t, createContent)
		t.AssertNE(createEnv.Code, gcode.CodeOK.Code())

		// Test creating a menu with a missing type
		createContent = client.PostContent(ctx, "/sys-menu", `{"name":"ApiMenuValidation","path":"/api-menu-validation"}`)
		createEnv = decodeMenuEnvelope(t, createContent)
		t.AssertNE(createEnv.Code, gcode.CodeOK.Code())

		// Test updating a menu with a missing name
		updatePayload := `{"path":"/api-menu-updated"}`
		updateContent := client.PutContent(ctx, "/sys-menu/some-id", updatePayload)
		updateEnv := decodeMenuEnvelope(t, updateContent)
		t.AssertNE(updateEnv.Code, gcode.CodeOK.Code())
	})
}

