package wechat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var wechat = Wechat{
	APPID: EnvAPPID,
	APPSecret: EnvAPPSecret,
}
func TestGetAccessToken (t *testing.T) {
	firstAccessToken := wechat.GetAccessToken()
	tokenLen := len(firstAccessToken)
	assert.EqualValues(t, 136<= tokenLen && tokenLen <= 157,true)
	// check cache
	secondAccessToken := wechat.GetAccessToken()
	assert.Equal(t, firstAccessToken, secondAccessToken)
}

func TestGetShortURL (t *testing.T) {
	// https://w.url.cn/s/A7b7sXQ
	firstShortURL := wechat.GetShortURL("https://github.com/og")
	assert.Regexp(t, "^https://w\\.url\\.cn/.*", firstShortURL)
	// check cache
	secondShortURL := wechat.GetShortURL("https://github.com/og")
	assert.Regexp(t, firstShortURL, secondShortURL)
}
