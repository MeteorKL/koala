package session

import (
	"time"
	"strconv"
	"github.com/MeteorKL/koala/util"
	"math/rand"
	"net/http"
	"sync"
)

type SessionStore struct {
	sessions   map[string]*Session
	expireTime time.Duration
}

const DEFAULT_EXPIRE_TIME = 3600

func NewSessionStore(expireTime time.Duration) *SessionStore {
	return &SessionStore{
		sessions:   make(map[string]*Session),
		expireTime: expireTime,
	}
}

func generateSessionID() string {
	return util.HashString(strconv.Itoa(rand.Int()) + time.Now().Format("2006-01-02 15:04:05"))
}

func (store SessionStore) newSession(w http.ResponseWriter, cookieName string) *Session {
	id := generateSessionID()
	store.sessions[id] = &Session{
		id:      id,
		values:  make(map[string]interface{}),
		expires: time.Now().Add(store.expireTime),
		lock:    sync.RWMutex{},
	}
	newCookie(w, cookieName, id)
	return store.sessions[id]
}

// ExistSession 判断是否登录
func (store SessionStore) ExistSession(r *http.Request, cookieName string) bool {
	c, err := r.Cookie(cookieName)
	if err != nil { // 不存在 cookie
		return false
	}
	return store.getSession(c.Value) != nil
}

func (store SessionStore) getSession(id string) *Session {
	if s, ok := store.sessions[id]; ok {
		return s
	}
	return nil
}

// GetSession 创建一个Session，如果存在直接返回
func (store SessionStore) GetSession(r *http.Request, w http.ResponseWriter, cookieName string) *Session {
	c, _ := r.Cookie(cookieName)
	var s *Session
	if c == nil {
		s = store.newSession(w, cookieName)
		return s
	}
	s = store.getSession(c.Value)
	if s == nil {
		s = store.newSession(w, cookieName)
	}
	return s
}

// PeekSession 如果有Session则直接返回,否则返回nil
func (store SessionStore) PeekSession(r *http.Request, cookieName string) *Session {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return nil
	}
	return store.getSession(c.Value)
}

// PeekSessionValue 获取session中的某个值，如果没有则返回nil
func (store SessionStore) PeekSessionValue(r *http.Request, cookieName string, key string) interface{} {
	s := store.PeekSession(r, cookieName)
	if s == nil {
		return nil
	}
	return s.Get(key)
}

func (store SessionStore) DelSession(r *http.Request, w http.ResponseWriter, cookieName string) error {
	s := store.PeekSession(r, cookieName)
	if s != nil {
		delete(store.sessions, s.id)
	}
	delCookie(w, cookieName)
	return nil
}

func (store SessionStore) UpdateExpires(r *http.Request, w http.ResponseWriter, cookieName string, s *Session) error {
	expireTime := store.expireTime
	s.updateExpires(expireTime)
	return nil
}
