package gwechat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


var wechat = New(Config{
	APPID: TestEnvAPPID,
	APPSecret: TestEnvAPPSecret,
	CenterService: namePublicAccountCenterService{},
})

type wechatHook struct {}
var memoryCache = &MemoryCache{}
var wechatMemeberCache = MemoryCache{}
type namePublicAccountCenterService struct {}
func (self namePublicAccountCenterService) GetAccessToken (appID string, appSecret string) (accessToken string , errRes ErrResponse){
	return UnsafeGetAccessToken(appID, appSecret, &wechatMemeberCache)
}

func (self namePublicAccountCenterService) GetJSAPITicket(appID string, appSecret string) (ticket string, errRes ErrResponse){
	accessToken , errRes := self.GetAccessToken(appID, appSecret)
	if errRes.Fail {return "", errRes}
	return UnsafeGetJSAPITicket(appID, accessToken, &wechatMemeberCache)
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
		)
		assert.Equal(t, "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + TestEnvAPPID + "&redirect_uri=https%3A%2F%2Fgithub.com%2Fog%2Fgowecaht&response_type=code&scope=snsapi_base&state=WECHAT_AUTH#wechat_redirect", url)
	}
	{
		url := wechat.WebRedirectAuthorize(
			Dict().WebRedirectAuthorize.Scope.SnsapiBase,
			"https://github.com/og/gowecaht?a=1&b=2",
			"WECHAT_AUTH",
		)
		assert.Equal(t, "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + TestEnvAPPID + "&redirect_uri=https%3A%2F%2Fgithub.com%2Fog%2Fgowecaht%3Fa%3D1%26b%3D2&response_type=code&scope=snsapi_base&state=WECHAT_AUTH#wechat_redirect", url)
	}
}

func TestWechat_CodeAccessToken(t *testing.T) {
	code := "0712lV2W1Eyt1Y0cc51W1NwR2W12lV2n"
	wechat.WebAccessToken(code)
}
