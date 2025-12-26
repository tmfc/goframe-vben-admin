package testutil

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
)

// RequireDatabase skips the test when no database configuration is available.
func RequireDatabase(t *testing.T) {
	t.Helper()
	ctx := context.TODO()
	cfg, err := g.Cfg().Get(ctx, "database.default.link")
	if err != nil {
		t.Skipf("skip: failed to read database config: %v", err)
		return
	}
	if cfg == nil || strings.TrimSpace(cfg.String()) == "" {
		t.Skip("skip: database.default.link not configured")
	}
}
