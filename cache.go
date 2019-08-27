package wechat

import (
	"time"
)

type defaultCache struct {
}
func DefaultCache () defaultCache {
	return defaultCache{}
}
func (cache defaultCache) Read(key string) (value string, has bool) {
	return "", false
}
func (cache defaultCache) Write(key string, value string, expiration time.Duration) {

}
