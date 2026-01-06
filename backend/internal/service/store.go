package service

import "sync"

var refreshTokenCache RefreshTokenStore = newRefreshTokenStore()

type RefreshTokenStore interface {
	Add(token, userID string)
	Remove(token string)
	Replace(oldToken, newToken, userID string)
	Valid(token, userID string) bool
}

type refreshTokenStore struct {
	sync.RWMutex
	tokens map[string]string
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
