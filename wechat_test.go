package gwechat

import (
	gjson "github.com/og/go-json"
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

func TestWechat_WeappCode2Session(t *testing.T) {
	{
		// 不传code
		res, errRes := wechat.WeappCode2Session("")
		assert.Equal(t, errRes.ErrCode, res.ErrCode)
		assert.Equal(t, errRes.ErrMsg, res.ErrMsg)
		assert.Equal(t, "", res.OpenID)
		assert.Equal(t, "", res.SessionKey)
		assert.Equal(t, true, errRes.Fail)
		assert.Equal(t, 41008, errRes.ErrCode)
		assert.Equal(t, true, len(errRes.ErrMsg)!=0)
	}
	{
		// 传新获取到的code
		res, errRes := wechat.WeappCode2Session(TestWeappClientCode)
		assert.Equal(t, errRes.ErrCode, res.ErrCode)
		assert.Equal(t, errRes.ErrMsg, res.ErrMsg)
		assert.Equal(t, 28, len(res.OpenID))
		assert.Equal(t, true, len(res.SessionKey)!=0)
		assert.Equal(t, false, errRes.Fail)
		assert.Equal(t, 0, errRes.ErrCode)
		assert.Equal(t, true, len(errRes.ErrMsg)==0)
	}

}
func TestWechat_DecodeWeappUserInfo(t *testing.T) {
	data, err := wechat.DecodeWeappUserInfo(
		"GpI3eix4bY79DNaV3si5Lw==",
		"6EuA3YRNWjv1ufVDCm+EzN6408N6odLHkkhVlGqBggX6oHyLrY/8j1zRKmeSlbRsMBLV0cgXlrdvo5izOIa2/Xzn4EmjD193FMxS/Umw0MSfrwQwLf07oo2qy3XSZAzdYIXf1morQqeYf1ZRrQ/bCQYPPQYoWFzvkG7NZ44VoZtbc9y5vzOei78xD0joILVwRgG0Ksg27SQZGZwryX+DOc3C6dDIW0UXNE3rwKsw0lYK1y5I6J/oJI20ZPOaPal2tKPBgjzTs6C1aJhC+rguOup2O3365sds3kVM9gIIUOLbQBySFHS12zcbmzWdbzledXVsKiu3YGOGB9crempyx9b2mD2wXQ5V/MtkxF9ssVguI2VJ1mSJUPQ+A70VJ4IPXWtko8zYnH4TaXyBGT6mOGMI8swSHeYJN9uE5sdZtAnSvjZQpVa4mEI+HxE17jXFF45T8W8m6ZMZDHAPHOIhiJzo2QIMsJINowoxpppKS8JkVsb/GAluVgfzbEsVxbFg+6p8J8QG+e/XYf6DxskjDw==",
		"lUlNEr2sUs5K8+vopobWMQ==",
		)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", data.OpenID)
	assert.Equal(t, TestEnvAPPID, data.Watermark.APPID)
	assert.NotEqual(t, 0, data.Watermark.Timestamp)
	assert.Equal(t, `{
  "avatarUrl": "https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTK6IxCxniax4vfhVw0d0fZdbqdicn59FH5Mqw0uj1AJiaOQ7BYcTnFqCtEhcpYGXeTrN0Bey6fA92VYQ/132",
  "city": "Hongkou",
  "country": "China",
  "gender": 1,
  "language": "zh_CN",
  "nickName": "储国柱",
  "province": "Shanghai",
  "openId": "ovv4N5M-IGBKOD1PXArPWw7HlvPU",
  "unionId": "os17EwbVY3o2tonH85y3Rs_417wQ",
  "watermark": {
    "appid": "wx9f01246e31fd5cae",
    "timestamp": 1568278735
  }
}`, gjson.StringUnfold(data))
	expectdData := WeappUserInfoData{
		AvatarURL: "https://wx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTK6IxCxniax4vfhVw0d0fZdbqdicn59FH5Mqw0uj1AJiaOQ7BYcTnFqCtEhcpYGXeTrN0Bey6fA92VYQ/132",
		City: "Hongkou",
		Country: "China",
		Gender: 1,
		Language: "zh_CN",
		NickName: "储国柱",
		Province: "Shanghai",
		OpenID: "ovv4N5M-IGBKOD1PXArPWw7HlvPU",
		UnionID: "os17EwbVY3o2tonH85y3Rs_417wQ",
		Watermark: WeappUserInfoDataWatermark{APPID: "wx9f01246e31fd5cae", Timestamp: 1568278735,},
	}
	assert.Equal(t, expectdData, data)

}
