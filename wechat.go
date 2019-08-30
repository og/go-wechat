package gwechat

import (
	"bytes"
	"encoding/json"
	qs "github.com/google/go-querystring/query"
	ge "github.com/og/go-error"
	"io/ioutil"
	"log"
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

type wechatErrorJSON struct {
	ErrCode int `json:"errcode"`
	ErrMsg string `json:"errmsg"`
}
type Hook interface {
	GetAccessToken (appID string, appSecret string) (accessToken string , err error)
	GetJSAPITicket(appID string, appSecret string) (ticket string, err error)
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

// 获取中控制平台的 jsapi_ticket
func (this Wechat) GetJSAPITicket () (accessToken string, err error) {
	return this.hook.GetJSAPITicket(this.appID, this.appSecret)
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


// 微信网页授权(第一步：用户同意授权，获取code)
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html#0
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


// 微信网页授权(第二步：通过code换取网页授权access_token)
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html#1
type WebAccessTokenResponse  struct {
	wechatErrorJSON
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID string `json:"openid"`
	Scope string `json:"scope"`
}
func (this Wechat) WebAccessToken(code string) (webAccessTokenResponse WebAccessTokenResponse, errRes ErrResponse) {
	type queryT struct {
		Code string `url:"code"`
		APPID string `url:"appid"`
		Secret string `url:"secret"`
		GrantType string `url:"grant_type"`
	}
	query := queryT{
		Code:      code,
		APPID:     this.appID,
		Secret:    this.appSecret,
		GrantType: "authorization_code",
	}
	querystring, err := qs.Values(query); ge.Check(err)
	requestURL := apiDomain + "/sns/oauth2/access_token" + "?" + querystring.Encode()
	resp, err := http.Get(requestURL); ge.Check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body); ge.Check(err)
	err = json.Unmarshal(body, &webAccessTokenResponse); ge.Check(err)
	if webAccessTokenResponse.ErrCode != 0 {
		errRes.Fail = true
		errRes.ErrCode = webAccessTokenResponse.ErrCode
		errRes.ErrMsg = webAccessTokenResponse.ErrMsg
		return
	}
	return
}

type WebUserInfoResponse struct {
	wechatErrorJSON
	OpenID string `json:"openid"`
	Nickname string `json:"nickname"`
	Sex int `json:"sex"`
	Province string `json:"province"`
	City string `json:"city"`
	Country string `json:"country"`
	HeadIMGURL string `json:"headimgurl"`
	Privilege []string `json:"privilege"`
	Unionid string `json:"unionid";comment:"只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。https://open.weixin.qq.com/cgi-bin/index?t=home/index&lang=zh_CN&token=3910897fc2d64d5279f371701325a78824caac9b"`
}
// 微信网页授权(第四步：拉取用户信息(需scope为 snsapi_userinfo))
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html#3
func (this Wechat) WebGetUserInfo(accessToken string, openID string, lang string) (wechatRes WebUserInfoResponse, errRes ErrResponse) {
	type request struct {
		AccessToken string `url:"access_token"`
		OpenID string `url:"openid"`
		Lang string `url:"lang"`
	}
	query := request{
		AccessToken: accessToken,
		OpenID: openID,
		Lang: lang,
	}
	querystring, err := qs.Values(query); ge.Check(err)
	requestURL := apiDomain + "/sns/userinfo" + "?" + querystring.Encode()
	resp, err := http.Get(requestURL); ge.Check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body); ge.Check(err)
	log.Print(string(body))
	err = json.Unmarshal(body, &wechatRes); ge.Check(err)
	if wechatRes.ErrCode !=0 {
		errRes.Fail = true
		errRes.ErrCode = wechatRes.ErrCode
		errRes.ErrMsg = wechatRes.ErrMsg
		return
	}
	return
}
