package service

import (
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

const (
	casbinDefaultModelPath = "resource/casbin/model.conf"
	casbinDefaultDomain    = "default"
)

var (
	casbinOnce      sync.Once
	casbinEnforcer  *casbin.Enforcer
	casbinInitError error
)

// Casbin returns a singleton enforcer initialized with the local adapter.
func Casbin(ctx context.Context) (*casbin.Enforcer, error) {
	casbinOnce.Do(func() {
		modelPath, err := resolveCasbinModelPath(ctx)
		if err != nil {
			casbinInitError = err
			return
		}
		enforcer, err := casbin.NewEnforcer(modelPath, NewCasbinAdapter(ctx))
		if err != nil {
			casbinInitError = err
			return
		}
		// 确保策略/分组变更自动通过 Adapter 持久化到 casbin_rule
		enforcer.EnableAutoSave(true)
		casbinEnforcer = enforcer
	})
	return casbinEnforcer, casbinInitError
}

func resolveCasbinModelPath(ctx context.Context) (string, error) {
	modelPath := casbinDefaultModelPath
	if cfgValue, err := g.Cfg().Get(ctx, "casbin.model"); err == nil && cfgValue != nil {
		if value := strings.TrimSpace(cfgValue.String()); value != "" {
			modelPath = value
		}
	}
	if gfile.Exists(modelPath) {
		return gfile.RealPath(modelPath), nil
	}
	if found, _ := gfile.Search(modelPath, gfile.Pwd()); found != "" {
		return found, nil
	}
	return "", gerror.Newf("casbin model file not found: %s", modelPath)
}

// NormalizeDomain ensures a non-empty Casbin domain.
func NormalizeDomain(domain string) string {
	if strings.TrimSpace(domain) == "" {
		return casbinDefaultDomain
	}
	return domain
}

func accessCodesFromCasbin(ctx context.Context, domain string, roles []string) ([]string, error) {
	enforcer, err := Casbin(ctx)
	if err != nil || enforcer == nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, nil
	}
	domain = NormalizeDomain(domain)
	set := make(map[string]struct{})
	for _, role := range roles {
		role = strings.TrimSpace(role)
		if role == "" {
			continue
		}
		permissions := enforcer.GetPermissionsForUserInDomain(role, domain)
		for _, perm := range permissions {
			if len(perm) < 3 {
				continue
			}
			code := strings.TrimSpace(perm[2])
			if code == "" {
				continue
			}
			set[code] = struct{}{}
		}
	}
	if len(set) == 0 {
		return nil, nil
	}
	codes := make([]string, 0, len(set))
	for code := range set {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes, nil
}
