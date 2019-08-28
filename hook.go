package gwechat

import (
	"sync"
	"time"
)

type HookInterface interface {
	ShortURLReadStorage (longURL string) (shortURL string, has bool)
	ShortURLWriteStorage (longURL string, shortURL string)
	ReadStorageCentralControlServerGetAccessToken(key string) (value string, has bool)
	WriteStorageCentralControlServerGetAccessToken(key string, value string, expiration time.Duration)
	RequestCentralControlServerAccessToken(appID string) (accessToken string , err error)
}

type defaultHook struct {

}

type memoryMapT struct {
	sync.RWMutex
	m map[string]string
}
var memoryMap memoryMapT
func (hook defaultHook) ReadStorageCentralControlServerGetAccessToken(key string) (value string, has bool) {
	memoryMap.RLock()
	value, has = memoryMap.m[key]
	memoryMap.RUnlock()
	return
}
func (hook defaultHook) WriteStorageCentralControlServerGetAccessToken(key string, value string, expiration time.Duration) {
	memoryMap.Lock()
	if memoryMap.m == nil {
		memoryMap.m = map[string]string{}
	}
	memoryMap.m[key] = value
	memoryMap.Unlock()
	if expiration != 0 {
		time.AfterFunc(expiration, func() {
			memoryMap.Lock()
			delete(memoryMap.m, key)
			memoryMap.Unlock()
		})
	}
}
func DefaultHook() defaultHook {
	return defaultHook{}
}


