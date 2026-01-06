package dao

import (
	"context"
	"strings"

	"backend/internal/consts"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// WithoutTenant marks context to bypass tenant scoping (for explicit cross-tenant ops).
func WithoutTenant(ctx context.Context) context.Context {
	return context.WithValue(ctx, consts.CtxKeySkipTenant, true)
}

func tenantIDFromCtx(ctx context.Context) string {
	if v := ctx.Value(consts.CtxKeyTenantID); v != nil {
		if tenantID, ok := v.(string); ok && strings.TrimSpace(tenantID) != "" {
			return tenantID
		}
	}
	return consts.DefaultTenantID
}

func withTenant(ctx context.Context, m *gdb.Model) *gdb.Model {
	if skip, ok := ctx.Value(consts.CtxKeySkipTenant).(bool); ok && skip {
		return m
	}
	tenantID := tenantIDFromCtx(ctx)
	return m.
		Where("tenant_id", tenantID).
		Data(g.Map{"tenant_id": tenantID})
}
