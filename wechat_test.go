package wechat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var wechat = Wechat{
	APPID: EnvAPPID,
	APPSecret: EnvAPPSecret,
	Cache: DefaultCache(),
	Hook: DefaultHook(),
}



func TestGetAccessToken (t *testing.T) {
	firstAccessToken := wechat.GetAccessToken()
	tokenLen := len(firstAccessToken)
	assert.EqualValues(t, 136<= tokenLen && tokenLen <= 157,true)
	// check cache
	// secondAccessToken := wechat.GetAccessToken()
	// assert.Equal(t, firstAccessToken, secondAccessToken)
}

func TestGetShortURL (t *testing.T) {
	// https://w.url.cn/s/A7b7sXQ
	firstShortURL := wechat.GetShortURL("https://github.com/og")
	assert.Regexp(t, "^https://w\\.url\\.cn/.*", firstShortURL)
	// check cache
	// secondShortURL := wechat.GetShortURL("https://github.com/og")
	// assert.Regexp(t, firstShortURL, secondShortURL)
}

func TestWechat_WebRedirectAuthorize(t *testing.T) {
	{
		url := wechat.WebRedirectAuthorize(
			Dict().WebRedirectAuthorize.Scope.SnsapiBase,
			"https://github.com/og/gowecaht",
			"WECHAT_AUTH",
		)
		assert.Equal(t, "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx25076a02429daba9&redirect_uri=https%3A%2F%2Fgithub.com%2Fog%2Fgowecaht&response_type=code&scope=snsapi_base&state=WECHAT_AUTH#wechat_redirect", url)
	}
	{
		url := wechat.WebRedirectAuthorize(
			Dict().WebRedirectAuthorize.Scope.SnsapiBase,
			"https://github.com/og/gowecaht?a=1&b=2",
			"WECHAT_AUTH",
		)
		assert.Equal(t, "https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx25076a02429daba9&redirect_uri=https%3A%2F%2Fgithub.com%2Fog%2Fgowecaht%3Fa%3D1%26b%3D2&response_type=code&scope=snsapi_base&state=WECHAT_AUTH#wechat_redirect", url)
	}
}