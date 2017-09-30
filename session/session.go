package session

import (
	"time"
	"sync"
)

type Session struct {
	id      string
	values  map[string]interface{}
	expires time.Time
	lock    sync.RWMutex
}

func (s *Session) updateExpires(expireTime time.Duration) {
	s.expires = time.Now().Add(expireTime)
}

func (s *Session) Set(k string, v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values[k] = v
}

func (s *Session) Get(k string) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if v, ok := s.values[k]; ok {
		return v
	}
	return nil
}

func (s *Session) Delete(k string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.values, k)
}
