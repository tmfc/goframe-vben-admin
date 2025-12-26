package service

import (
	"context"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"backend/api/auth/v1"
	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v4"
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

var (
	jwtSecret          []byte
	refreshTokenSecret []byte
	refreshTokenCache  RefreshTokenStore = newRefreshTokenStore()
	localAuth          IAuth
)

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

type refreshTokenStore struct {
	sync.RWMutex
	tokens map[string]string
}

type RefreshTokenStore interface {
	Add(token, userID string)
	Remove(token string)
	Replace(oldToken, newToken, userID string)
	Valid(token, userID string) bool
}

func newRefreshTokenStore() *refreshTokenStore {
	return &refreshTokenStore{
		tokens: make(map[string]string),
	}
}

func (s *refreshTokenStore) Add(token, userID string) {
	s.Lock()
	defer s.Unlock()
	s.tokens[token] = userID
}

func (s *refreshTokenStore) Remove(token string) {
	s.Lock()
	defer s.Unlock()
	delete(s.tokens, token)
}

func (s *refreshTokenStore) Replace(oldToken, newToken, userID string) {
	s.Lock()
	defer s.Unlock()
	delete(s.tokens, oldToken)
	s.tokens[newToken] = userID
}

func (s *refreshTokenStore) Valid(token, userID string) bool {
	s.RLock()
	defer s.RUnlock()
	if stored, ok := s.tokens[token]; ok {
		return stored == userID
	}
	return false
}

type sAuth struct{}

func init() {
	RegisterAuth(NewAuth())
}

func NewAuth() *sAuth {
	return &sAuth{}
}

func RegisterRefreshTokenStore(store RefreshTokenStore) {
	if store != nil {
		refreshTokenCache = store
	}
}

func loadJWTSecrets() {
	if jwtSecret == nil {
		if value, err := loadSecretFromConfigOrEnv(envJWTSecret, cfgJWTSecretKey); err == nil {
			jwtSecret = value
		}
	}
	if refreshTokenSecret == nil {
		if value, err := loadSecretFromConfigOrEnv(envRefreshTokenSecret, cfgRefreshSecretKey); err == nil {
			refreshTokenSecret = value
		}
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

func ensureJWTSecret() error {
	if len(jwtSecret) > 0 {
		return nil
	}
	if value, err := loadSecretFromConfigOrEnv(envJWTSecret, cfgJWTSecretKey); err == nil {
		jwtSecret = value
		return nil
	} else {
		return err
	}
}

func ensureRefreshSecret() error {
	if len(refreshTokenSecret) > 0 {
		return nil
	}
	if value, err := loadSecretFromConfigOrEnv(envRefreshTokenSecret, cfgRefreshSecretKey); err == nil {
		refreshTokenSecret = value
		return nil
	} else {
		return err
	}
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

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	refreshTokenCache.Add(refreshToken, user.Id)
	setRefreshTokenCookie(ctx, refreshToken)

	homePath := user.HomePath
	if homePath == "" {
		homePath = "/dashboard"
	}

	out = &v1.LoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserInfo: struct {
			ID       string   `json:"id"`
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
		dao.SysUser.Columns().TenantId: "00000000-0000-0000-0000-000000000000",
	}).Insert()
	return err
}

// RefreshToken implements interface IAuth.RefreshToken.
func (s *sAuth) RefreshToken(ctx context.Context, in v1.RefreshTokenReq) (out *v1.RefreshTokenRes, err error) {
	tokenStr := resolveRefreshToken(ctx, in.RefreshToken)
	if tokenStr == "" {
		return nil, gerror.NewCode(consts.ErrorCodeRefreshTokenRequired, "refresh token is required")
	}

	claims, err := parseRefreshToken(tokenStr)
	if err != nil {
		return nil, err
	}

	userID, _ := claims["id"].(string)
	if userID == "" {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "invalid refresh token")
	}
	if !refreshTokenCache.Valid(tokenStr, userID) {
		return nil, gerror.NewCode(consts.ErrorCodeRefreshTokenInvalid, "invalid refresh token")
	}

	var user entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Id, userID).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user.Id == "" {
		return nil, gerror.NewCode(consts.ErrorCodeUserNotFound, "user not found")
	}

	accessToken, err := s.generateAccessToken(&user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.generateRefreshToken(&user)
	if err != nil {
		return nil, err
	}

	refreshTokenCache.Replace(tokenStr, refreshToken, userID)
	setRefreshTokenCookie(ctx, refreshToken)

	out = &v1.RefreshTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	writeAccessTokenResponse(ctx, accessToken)
	return
}

// Logout implements interface IAuth.Logout.
func (s *sAuth) Logout(ctx context.Context, in v1.LogoutReq) (out *v1.LogoutRes, err error) {
	tokenStr := resolveRefreshToken(ctx, in.RefreshToken)
	if tokenStr != "" {
		refreshTokenCache.Remove(tokenStr)
	}
	clearRefreshTokenCookie(ctx)
	return &v1.LogoutRes{}, nil
}

// GetAccessCodes implements interface IAuth.GetAccessCodes.
func (s *sAuth) GetAccessCodes(ctx context.Context, in v1.GetAccessCodesReq) (out *v1.GetAccessCodesRes, err error) {
	token, err := resolveAccessToken(ctx, in.Token)
	if err != nil {
		return nil, err
	}

	claims, err := parseToken(token)
	if err != nil {
		return nil, err
	}

	userID, _ := claims["id"].(string)
	if userID == "" {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "invalid access token")
	}

	var user entity.SysUser
	err = dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Id, userID).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user.Id == "" {
		return nil, gerror.NewCode(consts.ErrorCodeUserNotFound, "user not found")
	}

	roles := parseRoles(user.Roles)
	codes, err := accessCodesFromCasbin(ctx, user.TenantId, roles)
	if err != nil || len(codes) == 0 {
		codes = buildAccessCodes(roles)
	}

	out = &v1.GetAccessCodesRes{
		Codes: codes,
	}
	return
}

func (s *sAuth) generateAccessToken(user *entity.SysUser) (string, error) {
	if err := ensureJWTSecret(); err != nil {
		return "", gerror.NewCode(consts.ErrorCodeUnauthorized, err.Error())
	}
	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"tenantId": user.TenantId,
		"exp":      time.Now().Add(AccessTokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *sAuth) generateRefreshToken(user *entity.SysUser) (string, error) {
	if err := ensureRefreshSecret(); err != nil {
		return "", gerror.NewCode(consts.ErrorCodeUnauthorized, err.Error())
	}
	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"exp":      time.Now().Add(RefreshTokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshTokenSecret)
}

func resolveRefreshToken(ctx context.Context, provided string) string {
	if provided != "" {
		return provided
	}
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return ""
	}
	return req.Cookie.Get(refreshTokenCookieName).String()
}

func resolveAccessToken(ctx context.Context, provided string) (string, error) {
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

func setRefreshTokenCookie(ctx context.Context, token string) {
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

func clearRefreshTokenCookie(ctx context.Context) {
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return
	}
	req.Cookie.Remove(refreshTokenCookieName)
}

func parseRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	if tokenStr == "" {
		return nil, gerror.NewCode(consts.ErrorCodeRefreshTokenRequired, "refresh token is empty")
	}
	if err := ensureRefreshSecret(); err != nil {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, err.Error())
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "unexpected signing method")
		}
		return refreshTokenSecret, nil
	})
	if err != nil {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "invalid refresh token")
}

func writeAccessTokenResponse(ctx context.Context, token string) {
	if token == "" {
		return
	}
	req := g.RequestFromCtx(ctx)
	if req == nil {
		return
	}
	req.Response.Write(token)
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
