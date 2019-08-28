package gwechat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


var wechat = New(Config{
	APPID: EnvAPPID,
	APPSecret: EnvAPPSecret,
})

// func TestCentralControlServerGetAccessToken (t *testing.T) {
// 	wechatError := New(Config{
// 		APPID: "",
// 		APPSecret: "",
// 	})
// 	_, errRes := wechatError.CentralControlServerGetAccessToken()
// 	assert.Equal(t,	len(errRes.ErrMsg) != 0, true)
// 	assert.Equal(t, true, errRes.Fail)
// 	assert.Equal(t, 41002, errRes.ErrCode)
// }
// func TestUnsafeGetAccessToken (t *testing.T) {
// 	firstAccessToken, errRes := wechat.UnsafeCentralControlServerGetAccessToken()
// 	assert.Equal(t, ErrResponse{Fail: false, ErrCode:0, ErrMsg:"",}, errRes)
// 	tokenLen := len(firstAccessToken)
// 	assert.EqualValues(t, 136<= tokenLen && tokenLen <= 157,true)
// 	// check cache
// 	secondAccessToken, errRes:= wechat.UnsafeCentralControlServerGetAccessToken()
// 	assert.Equal(t, firstAccessToken, secondAccessToken)
// 	assert.Equal(t, ErrResponse{Fail: false, ErrCode:0, ErrMsg:"",}, errRes)
// }

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