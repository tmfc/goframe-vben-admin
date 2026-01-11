package middleware

import (
	"context"
	"strings"

	"backend/internal/consts"
	"backend/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
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
			writeAuthError(r, err, authLogContext{path: r.URL.Path, method: r.Method})
			return
		}
		claims, err := service.ParseAccessToken(token)
		if err != nil {
			writeAuthError(r, err, authLogContext{path: r.URL.Path, method: r.Method})
			return
		}

		userID := gconv.String(claims["id"])
		if strings.TrimSpace(userID) == "" {
			writeAuthError(r, gerror.NewCode(consts.ErrorCodeUnauthorized, "invalid token subject"), authLogContext{path: r.URL.Path, method: r.Method})
			return
		}

		isSuper := false
		if rawSuper, ok := claims["isSuper"]; ok {
			switch v := rawSuper.(type) {
			case bool:
				isSuper = v
			case string:
				isSuper = strings.EqualFold(strings.TrimSpace(v), "true")
			case float64:
				isSuper = v != 0
			}
		} else {
			roles := service.ParseRoles(claims["roles"])
			for _, rname := range roles {
				if rname == consts.RoleSuper {
					isSuper = true
					break
				}
			}
		}

		tenantID := gconv.String(claims["tenantId"])
		if isSuper {
			headerTenant := strings.TrimSpace(r.Header.Get("X-TENANT-ID"))
			if headerTenant != "" {
				tenantID = headerTenant
			}
		}
		if strings.TrimSpace(tenantID) == "" {
			tenantID = consts.DefaultTenantID
		}

		// Inject tenant into context for downstream usage.
		ctxWithTenant := context.WithValue(r.Context(), consts.CtxKeyTenantID, tenantID)
		r.SetCtx(ctxWithTenant)
		err = authorizeCasbin(ctxWithTenant, authRequest{
			userID:   userID,
			tenantID: tenantID,
			path:     r.URL.Path,
			method:   r.Method,
		})
		if err != nil {
			writeAuthError(r, err, authLogContext{
				userID: userID,
				tenant: tenantID,
				path:   r.URL.Path,
				method: r.Method,
			})
			return
		}
		r.Middleware.Next()
	}
}

type authRequest struct {
	userID   string
	tenantID string
	path     string
	method   string
}

type authLogContext struct {
	userID string
	tenant string
	path   string
	method string
}

func authorizeCasbin(ctx context.Context, req authRequest) error {
	if strings.TrimSpace(req.userID) == "" {
		return gerror.NewCode(consts.ErrorCodeUnauthorized, "invalid token subject")
	}
	entry, err := loadCachedUser(ctx, req.userID)
	if err != nil {
		return err
	}

	tenantID := gconv.String(entry.user.TenantId)
	if strings.TrimSpace(req.tenantID) != "" {
		tenantID = req.tenantID
	}
	domain := service.NormalizeDomain(tenantID)
	obj := req.path
	act := strings.ToLower(req.method)

	enforcer, err := service.Casbin(ctx)
	if err != nil {
		return err
	}

	for _, role := range entry.roles {
		role = strings.TrimSpace(role)
		if role == "" {
			continue
		}
		if role == consts.RoleSuper {
			// 超管直接放行，避免策略缺失或未同步导致阻断
			return nil
		}
		allowed, err := enforcer.Enforce(role, domain, obj, act)
		if err != nil {
			return err
		}
		if allowed {
			return nil
		}
	}

	return gerror.NewCode(consts.ErrorCodeUnauthorized, "permission denied")
}

func writeAuthError(r *ghttp.Request, err error, ctxInfo authLogContext) {
	logAuthError(r.Context(), err, ctxInfo)
	r.SetError(err)
	r.Exit()
}

func logAuthError(ctx context.Context, err error, ctxInfo authLogContext) {
	code := gerror.Code(err)
	if code == nil {
		g.Log().Errorf(ctx, "casbin authz failed user=%s tenant=%s path=%s method=%s err=%v", ctxInfo.userID, ctxInfo.tenant, ctxInfo.path, ctxInfo.method, err)
		return
	}
	if code == consts.ErrorCodeUnauthorized || code == consts.ErrorCodeUserNotFound {
		g.Log().Warningf(ctx, "casbin authz denied user=%s tenant=%s path=%s method=%s err=%v", ctxInfo.userID, ctxInfo.tenant, ctxInfo.path, ctxInfo.method, err)
		return
	}
	g.Log().Errorf(ctx, "casbin authz error user=%s tenant=%s path=%s method=%s err=%v", ctxInfo.userID, ctxInfo.tenant, ctxInfo.path, ctxInfo.method, err)
}
