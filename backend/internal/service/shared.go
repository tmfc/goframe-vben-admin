package service

import (
	"context"
	"encoding/json"
	"strings"

	"backend/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
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
	if tokenStr == "" {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "token is empty")
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "invalid token")
}

// ParseAccessToken exposes JWT parsing for middleware usage.
func ParseAccessToken(tokenStr string) (jwt.MapClaims, error) {
	return parseToken(tokenStr)
}

// ResolveAccessToken exposes access token extraction for middleware usage.
func ResolveAccessToken(ctx context.Context, provided string) (string, error) {
	return resolveAccessToken(ctx, provided)
}
