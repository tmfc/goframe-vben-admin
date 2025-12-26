package middleware

import (
	"strings"

	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

var publicPaths = map[string]struct{}{
	"/auth/login":   {},
	"/auth/refresh": {},
}

// CasbinAuthz enforces interface-level permission checks using Casbin.
func CasbinAuthz() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		if r.Method == "OPTIONS" {
			r.Middleware.Next()
			return
		}
		if _, ok := publicPaths[r.URL.Path]; ok {
			r.Middleware.Next()
			return
		}

		token, err := service.ResolveAccessToken(r.Context(), "")
		if err != nil {
			r.SetError(err)
			r.Exit()
			return
		}
		claims, err := service.ParseAccessToken(token)
		if err != nil {
			r.SetError(err)
			r.Exit()
			return
		}

		userID, _ := claims["id"].(string)
		if userID == "" {
			r.SetError(gerror.NewCode(consts.ErrorCodeUnauthorized, "invalid token subject"))
			r.Exit()
			return
		}

		var user entity.SysUser
		if err := dao.SysUser.Ctx(r.Context()).Where(dao.SysUser.Columns().Id, userID).Scan(&user); err != nil {
			r.SetError(err)
			r.Exit()
			return
		}
		if user.Id == "" {
			r.SetError(gerror.NewCode(consts.ErrorCodeUserNotFound, "user not found"))
			r.Exit()
			return
		}

		roles := service.ParseRoles(user.Roles)
		if len(roles) == 0 {
			roles = []string{consts.DefaultRole()}
		}

		tenantID := user.TenantId
		if claimTenant, ok := claims["tenantId"].(string); ok && strings.TrimSpace(claimTenant) != "" {
			tenantID = claimTenant
		}
		domain := service.NormalizeDomain(tenantID)
		obj := r.URL.Path
		act := strings.ToLower(r.Method)

		enforcer, err := service.Casbin(r.Context())
		if err != nil {
			r.SetError(err)
			r.Exit()
			return
		}

		for _, role := range roles {
			if strings.TrimSpace(role) == "" {
				continue
			}
			allowed, err := enforcer.Enforce(role, domain, obj, act)
			if err != nil {
				r.SetError(err)
				r.Exit()
				return
			}
			if allowed {
				r.Middleware.Next()
				return
			}
		}

		r.SetError(gerror.NewCode(consts.ErrorCodeUnauthorized, "permission denied"))
		r.Exit()
	}
}
