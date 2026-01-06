package service

import (
	"context"
	"net/http"
	"strings"

	"backend/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func ResolveAccessToken(ctx context.Context, provided string) (string, error) {
	if provided != "" {
		return provided, nil
	}
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return "", gerror.NewCode(consts.ErrorCodeUnauthorized, "missing authorization token")
	}
	header := req.Header.Get("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer"))
	if token == "" {
		return "", gerror.NewCode(consts.ErrorCodeUnauthorized, "missing authorization token")
	}
	return token, nil
}

func ResolveRefreshToken(ctx context.Context, provided string) string {
	if provided != "" {
		return provided
	}
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return ""
	}
	return req.Cookie.Get(refreshTokenCookieName).String()
}

func SetRefreshTokenCookie(ctx context.Context, token string) {
	if token == "" {
		return
	}
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return
	}
	req.Cookie.SetCookie(
		refreshTokenCookieName,
		token,
		"",
		"/",
		RefreshTokenTTL,
		ghttp.CookieOptions{
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			HttpOnly: true,
		},
	)
}

func ClearRefreshTokenCookie(ctx context.Context) {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return
	}
	req.Cookie.Remove(refreshTokenCookieName)
}

func WriteAccessTokenResponse(ctx context.Context, token string) {
	if token == "" {
		return
	}
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return
	}
	req.Response.Write(token)
}
