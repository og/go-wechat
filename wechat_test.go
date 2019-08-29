package gwechat

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)


var wechat = New(Config{
	APPID: EnvAPPID,
	APPSecret: EnvAPPSecret,
	Hook: wechatHook{},
})

type wechatHook struct {}
var accessTokenMemoryCache = &AccessTokenMemoryCache{}
func (self wechatHook) GetAccessToken(appID string, appSecret string) (accessToken string , err error) {
	accessToken, errRes := UnsafeGetAccessToken(appID, appSecret, accessTokenMemoryCache)
	if errRes.Fail { return  "", errors.New(errRes.ErrMsg) }
	return
}

func TestGetShortURL (t *testing.T) {
	// https://w.url.cn/s/A7b7sXQ
	firstShortURL, errRes := wechat.GetShortURL("https://github.com/og")
	assert.Regexp(t, "^https://w\\.url\\.cn/.*", firstShortURL)
	assert.Equal(t, ErrResponse{Fail: false, ErrCode:0, ErrMsg:"",}, errRes)
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

func TestWechat_CodeAccessToken(t *testing.T) {
	code := "0712lV2W1Eyt1Y0cc51W1NwR2W12lV2n"
	wechat.CodeAccessToken(code)
}

//