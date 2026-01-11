package service

import (
	"context"
	"os"
	"sort"
	"strings"
	"time"

	v1 "backend/api/auth/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenTTL  = time.Hour * 72
	RefreshTokenTTL = time.Hour * 24 * 30
)

const (
	refreshTokenCookieName = "jwt"
	envJWTSecret           = "JWT_SECRET"
	envRefreshTokenSecret  = "JWT_REFRESH_SECRET"
	cfgJWTSecretKey        = "auth.jwt_secret"
	cfgRefreshSecretKey    = "auth.refresh_secret"
)

var localAuth IAuth

var roleAccessCodes = map[string][]string{
	consts.RoleSuper: {
		"System:Menu:List",
		"System:Menu:Create",
		"System:Menu:Edit",
		"System:Menu:Delete",
		"System:Dept:List",
		"System:Dept:Create",
		"System:Dept:Edit",
		"System:Dept:Delete",
	},
	consts.RoleAdmin: {
		"System:Menu:List",
		"System:Menu:Edit",
		"System:Dept:List",
		"System:Dept:Edit",
	},
	consts.RoleUser: {
		"System:Menu:List",
		"System:Dept:List",
	},
	consts.RoleGuest: {
		"System:Menu:List",
	},
}

func Auth() IAuth {
	return localAuth
}

func RegisterAuth(i IAuth) {
	localAuth = i
}

type sAuth struct {
	tokens TokenService
	store  RefreshTokenStore
}

func init() {
	RegisterAuth(NewAuth())
}

func NewAuth() *sAuth {
	if tokenSvc == nil {
		tokenSvc = NewJWTTokenService()
	}
	if refreshTokenCache == nil {
		refreshTokenCache = newRefreshTokenStore()
	}
	return &sAuth{
		tokens: tokenSvc,
		store:  refreshTokenCache,
	}
}

func RegisterRefreshTokenStore(store RefreshTokenStore) {
	if store != nil {
		refreshTokenCache = store
	}
}

func loadSecretFromConfigOrEnv(envName, cfgKey string) ([]byte, error) {
	if cfgValue, err := g.Cfg().Get(context.Background(), cfgKey); err == nil && cfgValue != nil {
		if value := strings.TrimSpace(cfgValue.String()); value != "" {
			return []byte(value), nil
		}
	}
	value := strings.TrimSpace(os.Getenv(envName))
	if value == "" {
		return nil, gerror.New("missing required secret in config or environment: " + envName)
	}
	return []byte(value), nil
}

// IAuth defines the interface for authentication service.
type IAuth interface {
	Login(ctx context.Context, in v1.LoginReq) (out *v1.LoginRes, err error)
	RefreshToken(ctx context.Context, in v1.RefreshTokenReq) (out *v1.RefreshTokenRes, err error)
	Logout(ctx context.Context, in v1.LogoutReq) (out *v1.LogoutRes, err error)
	GetAccessCodes(ctx context.Context, in v1.GetAccessCodesReq) (out *v1.GetAccessCodesRes, err error)
	// Temp method for creating a user for testing
	CreateUserForTest(ctx context.Context, username, password string) error
}

// Login implements interface IAuth.Login.
func (s *sAuth) Login(ctx context.Context, in v1.LoginReq) (out *v1.LoginRes, err error) {
	var user *entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Username, in.Username).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, gerror.NewCode(consts.ErrorCodeUserNotFound, "User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password))
	if err != nil {
		return nil, gerror.NewCode(consts.ErrorCodeIncorrectPassword, "Incorrect password")
	}

	roles := parseRoles(user.Roles)
	if len(roles) == 0 {
		roles = []string{consts.DefaultRole()}
	}

	accessToken, err := s.tokens.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.tokens.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	s.store.Add(refreshToken, gconv.String(user.Id))
	SetRefreshTokenCookie(ctx, refreshToken)

	homePath := user.HomePath
	if homePath == "" {
		homePath = "/dashboard"
	}

	out = &v1.LoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserInfo: struct {
			ID       int64    `json:"id"`
			Username string   `json:"username"`
			RealName string   `json:"realName"`
			Avatar   string   `json:"avatar"`
			HomePath string   `json:"homePath"`
			Roles    []string `json:"roles"`
		}{
			ID:       user.Id,
			Username: user.Username,
			RealName: user.RealName,
			Avatar:   user.Avatar,
			HomePath: homePath,
			Roles:    roles,
		},
	}
	return
}

// CreateUserForTest creates a user for testing purposes.
func (s *sAuth) CreateUserForTest(ctx context.Context, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = dao.SysUser.Ctx(ctx).Data(g.Map{
		dao.SysUser.Columns().Username: username,
		dao.SysUser.Columns().Password: string(hashedPassword),
		dao.SysUser.Columns().TenantId: consts.DefaultTenantID,
	}).Insert()
	return err
}

// RefreshToken implements interface IAuth.RefreshToken.
func (s *sAuth) RefreshToken(ctx context.Context, in v1.RefreshTokenReq) (out *v1.RefreshTokenRes, err error) {
	tokenStr := resolveRefreshToken(ctx, in.RefreshToken)
	if tokenStr == "" {
		return nil, gerror.NewCode(consts.ErrorCodeRefreshTokenRequired, "refresh token is required")
	}

	claims, err := s.tokens.ParseRefreshToken(tokenStr)
	if err != nil {
		return nil, err
	}

	userID := gconv.Int64(claims["id"])
	if userID == 0 {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "invalid refresh token")
	}
	userIDStr := gconv.String(userID)
	if !s.store.Valid(tokenStr, userIDStr) {
		return nil, gerror.NewCode(consts.ErrorCodeRefreshTokenInvalid, "invalid refresh token")
	}

	var user entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Id, userID).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, gerror.NewCode(consts.ErrorCodeUserNotFound, "user not found")
	}

	accessToken, err := s.tokens.GenerateAccessToken(&user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.tokens.GenerateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	s.store.Replace(tokenStr, refreshToken, userIDStr)
	SetRefreshTokenCookie(ctx, refreshToken)

	out = &v1.RefreshTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return
}

// Logout implements interface IAuth.Logout.
func (s *sAuth) Logout(ctx context.Context, in v1.LogoutReq) (out *v1.LogoutRes, err error) {
	tokenStr := resolveRefreshToken(ctx, in.RefreshToken)
	if tokenStr != "" {
		s.store.Remove(tokenStr)
	}
	ClearRefreshTokenCookie(ctx)
	return &v1.LogoutRes{}, nil
}

// GetAccessCodes implements interface IAuth.GetAccessCodes.
func (s *sAuth) GetAccessCodes(ctx context.Context, in v1.GetAccessCodesReq) (out *v1.GetAccessCodesRes, err error) {
	token, err := ResolveAccessToken(ctx, in.Token)
	if err != nil {
		return nil, err
	}

	claims, err := s.tokens.ParseAccessToken(token)
	if err != nil {
		return nil, err
	}

	userID := gconv.Int64(claims["id"])
	if userID == 0 {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "invalid access token")
	}

	var user entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Id, userID).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, gerror.NewCode(consts.ErrorCodeUserNotFound, "user not found")
	}

	roles := parseRoles(user.Roles)
	codes, err := accessCodesFromCasbin(ctx, gconv.String(user.TenantId), roles)
	if err != nil || len(codes) == 0 {
		codes = buildAccessCodes(roles)
	}

	out = &v1.GetAccessCodesRes{
		Codes: codes,
	}
	return
}

func resolveRefreshToken(ctx context.Context, provided string) string {
	return ResolveRefreshToken(ctx, provided)
}

func buildAccessCodes(roles []string) []string {
	if len(roles) == 0 {
		roles = []string{consts.DefaultRole()}
	}
	set := make(map[string]struct{})
	for _, role := range roles {
		role = strings.TrimSpace(role)
		codes, ok := roleAccessCodes[role]
		if !ok {
			codes = roleAccessCodes[consts.DefaultRole()]
		}
		for _, code := range codes {
			set[code] = struct{}{}
		}
	}
	result := make([]string, 0, len(set))
	for code := range set {
		result = append(result, code)
	}
	sort.Strings(result)
	return result
}
