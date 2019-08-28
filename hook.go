package wechat

import (
	"sync"
	"time"
)

type defaultHook struct {
	sync.RWMutex
	m map[string]string
}
func (hook defaultHook) ShortURLReadStorage (longURL string) (shortURL string, has bool) {
	return "", false
}
func (hook defaultHook) ShortURLWriteStorage (longURL string, shortURL string) {

}
func (hook defaultHook) AccessTokenReadStorage(key string) (value string, has bool) {
	hook.RLock()
	value, has = hook.m[key]
	hook.RUnlock()
	return
}
func (hook defaultHook) AccessTokenWriteStorage(key string, value string, expiration time.Duration) {
	hook.Lock()
	hook.m[key] = value
	hook.Unlock()
	if expiration != 0 {
		time.AfterFunc(expiration, func() {
			hook.Lock()
			delete(hook.RLock.m, key)
			hook.Unlock()
		})
	}
}
func DefaultHook() defaultHook {
	return defaultHook{}
}


