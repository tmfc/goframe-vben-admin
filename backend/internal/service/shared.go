package service

import (
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
