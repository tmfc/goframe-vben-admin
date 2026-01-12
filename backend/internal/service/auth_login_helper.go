package service

import (
	"strings"

	"backend/internal/consts"
)

type loginIdentity struct {
	Username   string
	TenantCode string
	TenantID   string
	HasSuffix  bool
}

func parseLoginIdentity(raw string, multiTenant bool) loginIdentity {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return loginIdentity{
			Username: strings.TrimSpace(raw),
			TenantID: consts.DefaultTenantID,
		}
	}

	user := trimmed
	tenantCode := ""
	hasSuffix := false
	if at := strings.LastIndex(trimmed, "@"); at > 0 && at < len(trimmed)-1 {
		user = strings.TrimSpace(trimmed[:at])
		tenantCode = strings.TrimSpace(trimmed[at+1:])
		hasSuffix = tenantCode != ""
	}

	if !multiTenant || !hasSuffix {
		return loginIdentity{
			Username:  user,
			TenantID:  consts.DefaultTenantID,
			HasSuffix: hasSuffix,
		}
	}

	return loginIdentity{
		Username:   user,
		TenantCode: tenantCode,
		HasSuffix:  true,
	}
}

func hasSuperRole(roles []string) bool {
	for _, role := range roles {
		if strings.TrimSpace(role) == consts.RoleSuper {
			return true
		}
	}
	return false
}
