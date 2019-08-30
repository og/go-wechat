package gwechat

import (
	"encoding/json"
	qs "github.com/google/go-querystring/query"
	ge "github.com/og/go-error"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type MemoryCache struct {
	sync.RWMutex
	m map[string]string
}
func (self *MemoryCache) ReadCache (key string) (value string, has bool) {
	self.RLock()
	value, has = self.m[key]
	self.RUnlock()
	return
}

func (self *MemoryCache) WriteCache (key string, value string, expiration time.Duration) {
	self.Lock()
	if self.m == nil {
		self.m = map[string]string{}
	}
	self.m[key] = value
	self.Unlock()
	if expiration != 0 {
		time.AfterFunc(expiration, func() {
			self.Lock()
			delete(self.m, key)
			self.Unlock()
		})
	}
	return
}

type WecahtCentralControlServiceRes struct {
	Type string `json:"type"`
	Msg string `json:"msg"`
	Data struct{
		Type string `json:"type"`
		AccessToken string `json:"access_token"`
		JSAPITicket string `json:"jsapi_ticket"`
	}`json:"data"`
}

// 中央控制服务器获取access_token
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
// 多个服务端使用同一个 appid 的情况下应该只有一个服务端向微信获取  access_token，否则会导致 access_token 冲突
// 所以公众号开发者一定要使用中控服务器统一获取和刷新access_token，其他业务逻辑服务器所使用的access_token均来自于该中控服务器，不应该各自去刷新，否则容易造成冲突，导致access_token覆盖而影响业务；
type Cache interface {
	ReadCache (key string) (value string, has bool)
	WriteCache (key string, value string, expiration time.Duration)
}
func UnsafeGetAccessToken (appID string, appSecret string, cache Cache) (accessToken string, errRes ErrResponse) {
	cacheKey := "access_token:"+appID
	accessToken, hasCache := cache.ReadCache(cacheKey)
	if hasCache {
		return accessToken, errRes
	}
	apiPath := "/cgi-bin/token"
	type apiQuery struct {
		GrantType string `url:"grant_type"`
		APPID string `url:"appid"`
		Secret string `url:"secret"`
	}
	type apiResponse struct {
		wechatErrorJSON
		AccessToken string `json:"access_token"`
		ExpiresIn int `json:"expires_in"`
	}
	requestPATH := apiDomain + apiPath
	query := apiQuery{
		GrantType: "client_credential", // 获取access_token填写client_credential
		APPID: appID,
		Secret: appSecret,
	}
	queryValues, err := qs.Values(query); ge.Check(err)
	requestURL := requestPATH +  "?" + queryValues.Encode()
	res, err := http.Get(requestURL); ge.Check(err)
	if res != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body); ge.Check(err)
	var resData apiResponse
	err = json.Unmarshal(body, &resData); ge.Check(err)
	if resData.ErrCode != 0 {
		errRes.Fail = true
		errRes.ErrCode = resData.ErrCode
		errRes.ErrMsg = resData.ErrMsg
		return "", errRes
	}
	accessToken = resData.AccessToken
	// -120秒防止过期
	cache.WriteCache(cacheKey, accessToken, time.Duration(resData.ExpiresIn-120)*time.Second)
	return accessToken, errRes
}

// 获取授权页ticket
// https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Vendor_API_List.html#1
func UnsafeGetJSAPITicket(appID string, accessToken string, cache Cache) ( ticket string, errRes ErrResponse) {
	cacheKey := "jsapi_ticket:"+appID
	ticket, hasTicket := cache.ReadCache(cacheKey)
	if hasTicket {
		return
	}
	apiPath := "/cgi-bin/ticket/getticket"
	type apiQuery struct {
		AccessToken string `url:"access_token"`
		Type string `url:"type"`
	}
	type apiResponse struct {
		wechatErrorJSON
		Ticket string `json:"ticket"`
		ExpiresIn int `json:"expires_in"`
	}
	requestPATH := apiDomain + apiPath
	query := apiQuery{
		AccessToken: accessToken,
		Type: "jsapi",
	}
	queryValues, err := qs.Values(query); ge.Check(err)
	requestURL := requestPATH +  "?" + queryValues.Encode()
	res, err := http.Get(requestURL); ge.Check(err)
	if res != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body); ge.Check(err)
	var resData apiResponse
	err = json.Unmarshal(body, &resData); ge.Check(err)
	if resData.ErrCode != 0 {
		errRes.Fail = true
		errRes.ErrCode = resData.ErrCode
		errRes.ErrMsg = resData.ErrMsg
		return "", errRes
	}
	ticket = resData.Ticket
	// -120秒防止过期
	cache.WriteCache(cacheKey, ticket, time.Duration(resData.ExpiresIn-120)*time.Second)
	return
}