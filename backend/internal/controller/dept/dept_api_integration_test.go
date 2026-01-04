package dept

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"backend/api/dept/v1"
	"backend/internal/dao"
	"backend/internal/testutil"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
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

func TestDeptAPIEndpoints(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.TODO()
	ensureTestTenant(t, ctx, "00000000-0000-0000-0000-000000000000")

	t.Cleanup(func() {
		dao.SysDept.Ctx(ctx).Unscoped().WhereLike(dao.SysDept.Columns().Name, "ApiDept%").Delete()
	})

	s := startDeptAPIServer(t)
	client := g.Client().ContentJson()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

	gtest.C(t, func(t *gtest.T) {
		createContent := client.PostContent(ctx, "/sys-dept", `{"name":"ApiDeptRoot","parentId":"0","status":1,"order":1,"creatorId":1}`)
		createEnv := decodeEnvelope(t, createContent)
		t.Assert(createEnv.Code, gcode.CodeOK.Code())

		var createRes v1.CreateDeptRes
		t.AssertNil(json.Unmarshal(createEnv.Data, &createRes))
		t.AssertNE(createRes.Id, "")

		getContent := client.GetContent(ctx, fmt.Sprintf("/sys-dept/%s", createRes.Id))
		getEnv := decodeEnvelope(t, getContent)
		t.Assert(getEnv.Code, gcode.CodeOK.Code())

		var getRes struct {
			Name string `json:"name"`
		}
		t.AssertNil(json.Unmarshal(getEnv.Data, &getRes))
		t.Assert(getRes.Name, "ApiDeptRoot")

		updateContent := client.PutContent(ctx, fmt.Sprintf("/sys-dept/%s", createRes.Id), `{"name":"ApiDeptUpdated","parentId":"0","status":1,"order":2,"modifierId":1}`)
		updateEnv := decodeEnvelope(t, updateContent)
		t.Assert(updateEnv.Code, gcode.CodeOK.Code())

		listContent := client.GetContent(ctx, "/sys-dept/list?name=ApiDeptUpdated")
		listEnv := decodeEnvelope(t, listContent)
		t.Assert(listEnv.Code, gcode.CodeOK.Code())

		var listRes struct {
			List []struct {
				Name string `json:"name"`
			} `json:"list"`
		}
		t.AssertNil(json.Unmarshal(listEnv.Data, &listRes))
		t.Assert(hasDeptName(listRes.List, "ApiDeptUpdated"), true)

		treeContent := client.GetContent(ctx, "/sys-dept/tree")
		treeEnv := decodeEnvelope(t, treeContent)
		t.Assert(treeEnv.Code, gcode.CodeOK.Code())

		var treeRes struct {
			List []struct {
				Name     string `json:"name"`
				Children []struct {
					Name string `json:"name"`
				} `json:"children"`
			} `json:"list"`
		}
		t.AssertNil(json.Unmarshal(treeEnv.Data, &treeRes))
		t.Assert(containsDept(treeRes.List, "ApiDeptUpdated"), true)

		deleteContent := client.DeleteContent(ctx, fmt.Sprintf("/sys-dept/%s", createRes.Id))
		deleteEnv := decodeEnvelope(t, deleteContent)
		t.Assert(deleteEnv.Code, gcode.CodeOK.Code())

		missingContent := client.GetContent(ctx, fmt.Sprintf("/sys-dept/%s", createRes.Id))
		missingEnv := decodeEnvelope(t, missingContent)
		t.Assert(missingEnv.Code, gcode.CodeNotFound.Code())
	})
}

func startDeptAPIServer(t *testing.T) *ghttp.Server {
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

func decodeEnvelope(t *gtest.T, content string) apiEnvelope {
	t.Helper()
	var env apiEnvelope
	t.AssertNil(json.Unmarshal([]byte(content), &env))
	return env
}

func hasDeptName(items []struct {
	Name string `json:"name"`
}, name string) bool {
	for _, item := range items {
		if item.Name == name {
			return true
		}
	}
	return false
}

func containsDept(items []struct {
	Name     string `json:"name"`
	Children []struct {
		Name string `json:"name"`
	} `json:"children"`
}, name string) bool {
	for _, item := range items {
		if item.Name == name {
			return true
		}
		if containsChild(item.Children, name) {
			return true
		}
	}
	return false
}

func containsChild(items []struct {
	Name string `json:"name"`
}, name string) bool {
	for _, item := range items {
		if item.Name == name {
			return true
		}
	}
	return false
}
