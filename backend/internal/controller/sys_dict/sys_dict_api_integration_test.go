package sys_dict

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	v1 "backend/api/sys_dict/v1"
	v1data "backend/api/sys_dict_data/v1"
	v1type "backend/api/sys_dict_type/v1"
	"backend/internal/consts"
	"backend/internal/controller/sys_dict_data"
	"backend/internal/controller/sys_dict_type"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/service"
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

func TestSysDictAuditAndOptionsE2E(t *testing.T) {
	testutil.RequireDatabase(t)

	ctx := context.WithValue(context.TODO(), consts.CtxKeyTenantID, consts.DefaultTenantID)
	ensureTestTenant(t, ctx, consts.DefaultTenantID)

	t.Cleanup(func() {
		dao.SysDictData.Ctx(ctx).Unscoped().WhereLike(dao.SysDictData.Columns().Value, "ApiDictData%").Delete()
		dao.SysDictType.Ctx(ctx).Unscoped().WhereLike(dao.SysDictType.Columns().TypeCode, "ApiDictAudit%").Delete()
	})

	token := mustGenerateToken(t, 999, 1, "audit-user")
	s := startDictAPIServer(t)
	client := g.Client().ContentJson()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
	client.SetHeader("Authorization", "Bearer "+token)

	gtest.C(t, func(t *gtest.T) {
		createTypeContent := client.PostContent(ctx, "/sys-dict-type", `{"typeCode":"ApiDictAudit","typeName":"Api Dict Audit","status":1,"sort":1}`)
		createTypeEnv := decodeEnvelope(t, createTypeContent)
		t.Assert(createTypeEnv.Code, gcode.CodeOK.Code())

		var createTypeRes v1type.CreateDictTypeRes
		t.AssertNil(json.Unmarshal(createTypeEnv.Data, &createTypeRes))
		t.AssertNE(createTypeRes.Id, 0)

		var dictType entity.SysDictType
		err := dao.SysDictType.Ctx(ctx).Where(dao.SysDictType.Columns().Id, createTypeRes.Id).Scan(&dictType)
		t.AssertNil(err)
		t.Assert(dictType.CreatorId, int64(999))
		t.Assert(dictType.ModifierId, int64(999))

		updateTypeContent := client.PutContent(ctx, fmt.Sprintf("/sys-dict-type/%d", createTypeRes.Id), `{"typeCode":"ApiDictAudit","typeName":"Api Dict Audit Updated","status":1,"sort":2}`)
		updateTypeEnv := decodeEnvelope(t, updateTypeContent)
		t.Assert(updateTypeEnv.Code, gcode.CodeOK.Code())

		var updatedType entity.SysDictType
		err = dao.SysDictType.Ctx(ctx).Where(dao.SysDictType.Columns().Id, createTypeRes.Id).Scan(&updatedType)
		t.AssertNil(err)
		t.Assert(updatedType.ModifierId, int64(999))

		createDataContent := client.PostContent(ctx, "/sys-dict-data", fmt.Sprintf(`{"dictTypeId":%d,"label":"Default","labelI18n":{"en":"Default EN","zh-tw":"Default TW"},"value":"ApiDictDataValue","status":1,"sort":1}`, createTypeRes.Id))
		createDataEnv := decodeEnvelope(t, createDataContent)
		t.Assert(createDataEnv.Code, gcode.CodeOK.Code())

		var createDataRes v1data.CreateDictDataRes
		t.AssertNil(json.Unmarshal(createDataEnv.Data, &createDataRes))
		t.AssertNE(createDataRes.Id, 0)

		var dictData entity.SysDictData
		err = dao.SysDictData.Ctx(ctx).Where(dao.SysDictData.Columns().Id, createDataRes.Id).Scan(&dictData)
		t.AssertNil(err)
		t.Assert(dictData.CreatorId, int64(999))

		client.SetHeader("Accept-Language", "en")
		optionsContent := client.GetContent(ctx, "/sys-dict/options/ApiDictAudit")
		optionsEnv := decodeEnvelope(t, optionsContent)
		t.Assert(optionsEnv.Code, gcode.CodeOK.Code())

		updateDataContent := client.PutContent(ctx, fmt.Sprintf("/sys-dict-data/%d", createDataRes.Id), fmt.Sprintf(`{"dictTypeId":%d,"label":"Default Updated","labelI18n":{"en":"Default Updated EN"},"value":"ApiDictDataValue","status":1,"sort":2}`, createTypeRes.Id))
		updateDataEnv := decodeEnvelope(t, updateDataContent)
		t.Assert(updateDataEnv.Code, gcode.CodeOK.Code())

		var updatedData entity.SysDictData
		err = dao.SysDictData.Ctx(ctx).Where(dao.SysDictData.Columns().Id, createDataRes.Id).Scan(&updatedData)
		t.AssertNil(err)
		t.Assert(updatedData.ModifierId, int64(999))

		client.SetHeader("Accept-Language", "en")
		optionsContent = client.GetContent(ctx, "/sys-dict/options/ApiDictAudit")
		optionsEnv = decodeEnvelope(t, optionsContent)
		t.Assert(optionsEnv.Code, gcode.CodeOK.Code())

		var optionsRes v1.GetDictOptionsRes
		t.AssertNil(json.Unmarshal(optionsEnv.Data, &optionsRes))
		t.Assert(optionsRes.List[0].Label, "Default Updated EN")

		client.SetHeader("Accept-Language", "fr")
		labelContent := client.GetContent(ctx, "/sys-dict/label/ApiDictAudit/ApiDictDataValue")
		labelEnv := decodeEnvelope(t, labelContent)
		t.Assert(labelEnv.Code, gcode.CodeOK.Code())

		var labelRes v1.GetDictLabelRes
		t.AssertNil(json.Unmarshal(labelEnv.Data, &labelRes))
		t.Assert(labelRes.Label, "Default Updated")
	})
}

func startDictAPIServer(t *testing.T) *ghttp.Server {
	t.Helper()
	s := g.Server(guid.S())
	s.SetAddr(ghttp.FreePortAddress)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(NewV1())
		group.Bind(sys_dict_data.NewV1())
		group.Bind(sys_dict_type.NewV1())
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

func mustGenerateToken(t *testing.T, userID int64, tenantID int64, username string) string {
	t.Helper()
	user := &entity.SysUser{
		Id:       userID,
		TenantId: tenantID,
		Username: username,
		Roles:    "[\"super\"]",
	}
	tokenSvc := service.NewJWTTokenService()
	token, err := tokenSvc.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return token
}

func decodeEnvelope(t *gtest.T, content string) apiEnvelope {
	t.Helper()
	var env apiEnvelope
	t.AssertNil(json.Unmarshal([]byte(content), &env))
	return env
}
