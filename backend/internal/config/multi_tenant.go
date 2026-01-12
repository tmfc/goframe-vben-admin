package config

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

const envMultiTenant = "APP_MULTI_TENANT"

// IsMultiTenantEnabled returns whether multi-tenant mode is enabled.
func IsMultiTenantEnabled(ctx context.Context) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	if raw := strings.TrimSpace(os.Getenv(envMultiTenant)); raw != "" {
		if val, err := strconv.ParseBool(raw); err == nil {
			return val
		}
	}
	cfgValue, err := g.Cfg().Get(ctx, "app.multiTenant")
	if err != nil || cfgValue == nil {
		return false
	}
	return cfgValue.Bool()
}
