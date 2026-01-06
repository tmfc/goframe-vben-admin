package service

import (
	"context"
	"encoding/json"
	"strings"

	"backend/internal/consts"

	"github.com/golang-jwt/jwt/v4"
)

func parseRoles(raw any) []string {
	roles := make([]string, 0)
	switch v := raw.(type) {
	case nil:
		return roles
	case []string:
		return v
	case []any:
		for _, r := range v {
			if s, ok := r.(string); ok && strings.TrimSpace(s) != "" {
				roles = append(roles, strings.TrimSpace(s))
			}
		}
		return roles
	case string:
		raw = v
	default:
		return roles
	}

	rawStr := raw.(string)
	if rawStr == "" {
		return roles
	}
	if err := json.Unmarshal([]byte(rawStr), &roles); err == nil {
		return roles
	}
	if strings.Contains(rawStr, ",") {
		for _, role := range strings.Split(rawStr, ",") {
			if trimmed := strings.TrimSpace(role); trimmed != "" {
				roles = append(roles, trimmed)
			}
		}
		return roles
	}
	roles = append(roles, rawStr)
	return roles
}

// ParseRoles exposes role parsing for middleware usage.
func ParseRoles(raw any) []string {
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
