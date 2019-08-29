package gwechat

import (
	"bytes"
	"encoding/json"
	qs "github.com/google/go-querystring/query"
	ge "github.com/og/go-error"
	"io/ioutil"
	"net/http"
)
// 接口域名说明
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1465199793_BqlKA
const apiDomain = "https://api.weixin.qq.com"
const alternativeAPIDomain = "https://api2.weixin.qq.com"

type Wechat struct {
	appID string
	appSecret string
	hook Hook
}

type Hook interface {
	GetAccessToken (appID string, appSecret string) (accessToken string , err error)
}
type Config struct {
	APPID string
	APPSecret string
	Hook Hook
}
func New (config Config) Wechat {
	wechat := Wechat{
		appID: config.APPID,
		appSecret: config.APPSecret,
		hook: config.Hook,
	}
	return wechat
}

// 获取中控制平台的 access_token
func (this Wechat) GetAccessToken () (accessToken string, err error) {
	return this.hook.GetAccessToken(this.appID, this.appSecret)
}


// 长链接转短链接接口
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1443433600
func (this Wechat) GetShortURL (longURL string) (shortURL string, errRes ErrResponse)  {
	apiPath := "/cgi-bin/shorturl"
	type apiQuery struct {
		AccessToken string `url:"access_token"`
	}
	type apiParam struct {
		Action string `json:"action"`
		LongURL string `json:"long_url"`
	}
	type apiResponse struct {
		ErrCode int `json:"errcode"`
		ErrMsg string `json:"errmsg"`
		ShortURL string `json:"short_url"`
	}
	requestPATH := apiDomain + apiPath
	accessToken , err := this.GetAccessToken()
	if err != nil {
		errRes.ErrMsg = err.Error()
		return "", errRes
	}
	query := apiQuery{
		AccessToken: accessToken,
	}
	param := apiParam {
		Action: "long2short", // 此处填long2short，代表长链接转短链接
		LongURL: longURL,
	}
	paramJSON, err := json.Marshal(param); ge.Check(err)
	queryValues, err := qs.Values(query); ge.Check(err)
	requestURL := requestPATH +  "?" + queryValues.Encode()
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(paramJSON)); ge.Check(err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req); ge.Check(err)
	if res != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body); ge.Check(err)
	var respData apiResponse
	err = json.Unmarshal(body, &respData); ge.Check(err)
	if respData.ErrCode != 0 {
		errRes.Fail = true
		errRes.ErrCode = respData.ErrCode
		errRes.ErrMsg = respData.ErrMsg
		return "", errRes
	}
	shortURL = respData.ShortURL
	return
}


// 微信网页授权第一步跳转
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140842
// scope 参数使用 wecaht.Dict().WebRedirectAuthorize.Scope 传递
// redirectURI 授权后重定向的回调链接地址, 函数内部已进行 urlEncode 操作，调用方无需 urlEncode
// state 重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节

// 成功后如果用户同意授权，页面将跳转至 redirect_uri/?code=CODE&state=STATE。
func (this Wechat) WebRedirectAuthorize(scope string, redirectURI string, state string) string {
	type queryT struct {
		AppID string `url:"appid"`
		RedirectURI string `url:"redirect_uri"`
		ResponseType string `url:"response_type"`
		Scope string `url:"scope"`
		State string `url:"state"`
	}
	query := queryT {
		AppID: this.appID,
		RedirectURI: redirectURI,
		ResponseType: "code",
		Scope: scope,
		State: state,
	}
	querystring, err := qs.Values(query) ; ge.Check(err)
	return "https://open.weixin.qq.com/connect/oauth2/authorize?" + querystring.Encode() + "#wechat_redirect"
}
