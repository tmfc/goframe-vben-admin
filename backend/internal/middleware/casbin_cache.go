package middleware

import (
	"context"
	"sync"
	"time"

	"backend/internal/consts"
	"backend/internal/dao"
	"backend/internal/model/entity"
	"backend/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
)

type cachedUser struct {
	user      entity.SysUser
	roles     []string
	expiresAt time.Time
}

var (
	authUserCacheTTL = 2 * time.Minute
	authUserCacheMu  sync.RWMutex
	authUserCache    = make(map[string]cachedUser)
)

func getCachedUser(userID string) (cachedUser, bool) {
	authUserCacheMu.RLock()
	entry, ok := authUserCache[userID]
	authUserCacheMu.RUnlock()
	if !ok {
		return cachedUser{}, false
	}
	if time.Now().After(entry.expiresAt) {
		authUserCacheMu.Lock()
		delete(authUserCache, userID)
		authUserCacheMu.Unlock()
		return cachedUser{}, false
	}
	return entry, true
}

func setCachedUser(userID string, entry cachedUser) {
	authUserCacheMu.Lock()
	authUserCache[userID] = entry
	authUserCacheMu.Unlock()
}

func loadCachedUser(ctx context.Context, userID string) (cachedUser, error) {
	if entry, ok := getCachedUser(userID); ok {
		return entry, nil
	}

	var user entity.SysUser
	if err := dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns().Id, userID).Scan(&user); err != nil {
		return cachedUser{}, err
	}
	if user.Id == "" {
		return cachedUser{}, gerror.NewCode(consts.ErrorCodeUserNotFound, "user not found")
	}

	roles := service.ParseRoles(user.Roles)
	if len(roles) == 0 {
		roles = []string{consts.DefaultRole()}
	}

	entry := cachedUser{
		user:      user,
		roles:     roles,
		expiresAt: time.Now().Add(authUserCacheTTL),
	}
	setCachedUser(userID, entry)
	return entry, nil
}

func resetAuthCache() {
	authUserCacheMu.Lock()
	authUserCache = make(map[string]cachedUser)
	authUserCacheMu.Unlock()
}
