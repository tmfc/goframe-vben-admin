package service

import (
	"context"
	"encoding/json"
	"strings"

	"backend/internal/consts"

	"github.com/golang-jwt/jwt/v4"
)

func parseRoles(raw string) []string {
	roles := make([]string, 0)
	if raw == "" {
		return roles
	}
	if err := json.Unmarshal([]byte(raw), &roles); err == nil {
		return roles
	}
	if strings.Contains(raw, ",") {
		for _, role := range strings.Split(raw, ",") {
			if trimmed := strings.TrimSpace(role); trimmed != "" {
				roles = append(roles, trimmed)
			}
		}
		return roles
	}
	roles = append(roles, raw)
	return roles
}

// ParseRoles exposes role parsing for middleware usage.
func ParseRoles(raw string) []string {
	return parseRoles(raw)
}

func parseToken(tokenStr string) (jwt.MapClaims, error) {
	return tokenSvc.ParseAccessToken(tokenStr)
}

// ParseAccessToken exposes JWT parsing for middleware usage.
func ParseAccessToken(tokenStr string) (jwt.MapClaims, error) {
	return parseToken(tokenStr)
}

// TenantID returns resolved tenant ID with default fallback.
func TenantID(ctx context.Context) string {
	return resolveTenantID(ctx)
}

func resolveTenantID(ctx context.Context) string {
	const defaultTenantID = consts.DefaultTenantID

	if v := ctx.Value(consts.CtxKeyTenantID); v != nil {
		if tenantID, ok := v.(string); ok && strings.TrimSpace(tenantID) != "" {
			return tenantID
		}
	}

	token, err := ResolveAccessToken(ctx, "")
	if err != nil {
		return defaultTenantID
	}
	claims, err := parseToken(token)
	if err != nil {
		return defaultTenantID
	}
	tenantID, _ := claims["tenantId"].(string)
	if strings.TrimSpace(tenantID) == "" {
		return defaultTenantID
	}
	return tenantID
}
