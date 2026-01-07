package service

import (
	"strings"
	"sync"
	"time"

	"backend/internal/consts"
	"backend/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/golang-jwt/jwt/v4"
)

var tokenSvc TokenService = NewJWTTokenService()

type TokenService interface {
	GenerateAccessToken(user *entity.SysUser) (string, error)
	GenerateRefreshToken(user *entity.SysUser) (string, error)
	ParseAccessToken(tokenStr string) (jwt.MapClaims, error)
	ParseRefreshToken(tokenStr string) (jwt.MapClaims, error)
}

type jwtTokenService struct {
	accessSecret  []byte
	refreshSecret []byte
	once          sync.Once
	initErr       error
}

func NewJWTTokenService() TokenService {
	return &jwtTokenService{}
}

func RegisterTokenService(ts TokenService) {
	if ts != nil {
		tokenSvc = ts
	}
}

func (s *jwtTokenService) ensureSecrets() error {
	s.once.Do(func() {
		access, err := loadSecretFromConfigOrEnv(envJWTSecret, cfgJWTSecretKey)
		if err != nil {
			s.initErr = err
			return
		}
		refresh, err := loadSecretFromConfigOrEnv(envRefreshTokenSecret, cfgRefreshSecretKey)
		if err != nil {
			s.initErr = err
			return
		}
		s.accessSecret = access
		s.refreshSecret = refresh
	})
	return s.initErr
}

func (s *jwtTokenService) GenerateAccessToken(user *entity.SysUser) (string, error) {
	if err := s.ensureSecrets(); err != nil {
		return "", gerror.NewCode(consts.ErrorCodeUnauthorized, err.Error())
	}
	roles := parseRoles(user.Roles)
	isSuper := false
	for _, role := range roles {
		if strings.TrimSpace(role) == consts.RoleSuper {
			isSuper = true
			break
		}
	}
	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"tenantId": user.TenantId,
		"isSuper":  isSuper,
		"exp":      time.Now().Add(AccessTokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.accessSecret)
}

func (s *jwtTokenService) GenerateRefreshToken(user *entity.SysUser) (string, error) {
	if err := s.ensureSecrets(); err != nil {
		return "", gerror.NewCode(consts.ErrorCodeUnauthorized, err.Error())
	}
	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"exp":      time.Now().Add(RefreshTokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.refreshSecret)
}

func (s *jwtTokenService) ParseAccessToken(tokenStr string) (jwt.MapClaims, error) {
	if tokenStr == "" {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "token is empty")
	}
	if err := s.ensureSecrets(); err != nil {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, err.Error())
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "unexpected signing method")
		}
		return s.accessSecret, nil
	})
	if err != nil {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "invalid token")
}

func (s *jwtTokenService) ParseRefreshToken(tokenStr string) (jwt.MapClaims, error) {
	if tokenStr == "" {
		return nil, gerror.NewCode(consts.ErrorCodeRefreshTokenRequired, "refresh token is empty")
	}
	if err := s.ensureSecrets(); err != nil {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, err.Error())
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "unexpected signing method")
		}
		return s.refreshSecret, nil
	})
	if err != nil {
		return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, gerror.NewCode(consts.ErrorCodeUnauthorized, "invalid refresh token")
}
