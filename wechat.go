package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	qs "github.com/google/go-querystring/query"
	"github.com/og/so"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)
// 接口域名说明
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1465199793_BqlKA
var APIDomain = "https://api.weixin.qq.com"
var AlternativeAPIDomain = "https://api2.weixin.qq.com"

type Wechat struct {
	APPID string
	APPSecret string
}
const cacheKey_GetAccessToken = "og_wechat_go_get_access_token"
var memoryCache = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func (this Wechat) GetCache(key string) (value string, has bool) {
	memoryCache.RLock()
	value, has = memoryCache.m[key]
	memoryCache.RUnlock()
	return
}
func (this Wechat) SetCache(key string, value string, expiration time.Duration) {
	memoryCache.Lock()
	memoryCache.m[key] = value
	memoryCache.Unlock()
	if expiration != 0 {
		time.AfterFunc(expiration, func() {
			memoryCache.RLock()
			delete(memoryCache.m, key)
			memoryCache.RUnlock()
		})
	}
}

// 获取access_token
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140183
// 在内部实现了内存缓存 accessToken 目前(2019年08月)微信接口返回的 凭证有效时间 是 7200 秒,内存缓存保存 (7200 - 30) 秒
// 119 分钟内不会再次请求微信接口，这样可以提高响应速度
func (this Wechat) GetAccessToken () (accessToken string) {
	apiPath := "/cgi-bin/token"
	type apiQuery struct {
		GrantType string `url:"grant_type"`
		APPID string `url:"appid"`
		Secret string `url:"secret"`
	}
	type apiResponse struct {
		ErrCode int `json:"errcode"`
		ErrMsg string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn int `json:"expires_in"`
	}
	cacheValue, hasCacheValue := this.GetCache(cacheKey_GetAccessToken)
	if hasCacheValue {
		return cacheValue
	}
	requestPATH := APIDomain + apiPath
	query := apiQuery{
		GrantType: "client_credential", // 获取access_token填写client_credential
		APPID: this.APPID,
		Secret: this.APPSecret,
	}
	queryValues, err := qs.Values(query); so.C(err)
	requestURL := requestPATH +  "?" + queryValues.Encode()
	res, err := http.Get(requestURL); so.C(err)
	if res != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body); so.C(err)
	var resData apiResponse
	err = json.Unmarshal(body, &resData); so.C(err)
	if resData.ErrCode != 0 {
		panic(errors.New(fmt.Sprintf("%#v", resData)))
	}
	accessToken = resData.AccessToken
	// 减去30秒是防御措施，防止access token 失效前一秒进行请求网络返回 token 失效
	this.SetCache(cacheKey_GetAccessToken, accessToken, time.Duration(resData.ExpiresIn - 30) * time.Second)
	return accessToken
}
// 将查询结果缓存再内存中或 redis mysql mongodb 中（默认在内存中）
var shortURLMemoryStorage = make(map[string]string)
func (this Wechat) ShortURLHookReadStorage (longURL string) (shortURL string, has bool) {
	shortURL, has = shortURLMemoryStorage[longURL]
	return
}
// 只设置 5 秒的缓存，防止高并发大量url导致内存爆掉，如需更长的过期时间调用者最好自己存入持久化存储中
func (this Wechat) ShortURLHookWriteStorage (longURL string, shortURL string) () {
	time.AfterFunc(time.Duration(1 * time.Second), func() {
		delete(shortURLMemoryStorage, longURL)
	})
	shortURLMemoryStorage[longURL] = shortURL
	return
}
// 长链接转短链接接口
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1443433600
func (this Wechat) GetShortURL (longURL string) (shortURL string)  {
	storageValue, has := this.ShortURLHookReadStorage(longURL)
	if has {
		shortURL = storageValue
		return
	}
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
	requestPATH := APIDomain + apiPath
	query := apiQuery{
		AccessToken: this.GetAccessToken(),
	}
	param := apiParam {
		Action: "long2short", // 此处填long2short，代表长链接转短链接
		LongURL: longURL,
	}
	paramJSON, err := json.Marshal(param); so.C(err)
	queryValues, err := qs.Values(query); so.C(err)
	requestURL := requestPATH +  "?" + queryValues.Encode()
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(paramJSON)); so.C(err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req); so.C(err)
	if res != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body); so.C(err)
	var respData apiResponse
	err = json.Unmarshal(body, &respData); so.C(err)
	if respData.ErrCode != 0 {
		panic(errors.New(fmt.Sprintf("%#v", respData)))
	}
	shortURL = respData.ShortURL
	this.ShortURLHookWriteStorage(longURL, shortURL)
	return
}