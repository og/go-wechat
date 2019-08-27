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
	"strings"
	"time"
)
// 接口域名说明
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1465199793_BqlKA
var APIDomain = "https://api.weixin.qq.com"
var AlternativeAPIDomain = "https://api2.weixin.qq.com"

type CacheInterface interface {
	Read(key string) (value string, has bool)
	Write(key string, value string, expiration time.Duration)
}
type HookInterface interface {
	ShortURLReadStorage (longURL string) (shortURL string, has bool)
	ShortURLWriteStorage (longURL string, shortURL string)
}
type Wechat struct {
	APPID string
	APPSecret string
	Cache CacheInterface
	Hook HookInterface
}

func (this Wechat) getCacheKeyGetAccessToken () string {
	return strings.Join([]string{
		"og_wechat",
		this.APPID,
		"get_access_token",
	}, ":")
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
	cacheValue, hasCacheValue := this.Cache.Read(this.getCacheKeyGetAccessToken())
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
	this.Cache.Write(this.getCacheKeyGetAccessToken(), accessToken, time.Duration(resData.ExpiresIn - 30) * time.Second)
	return accessToken
}
// 长链接转短链接接口
// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1443433600
func (this Wechat) GetShortURL (longURL string) (shortURL string)  {
	storageValue, has := this.Hook.ShortURLReadStorage(longURL)
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
	this.Hook.ShortURLWriteStorage(longURL, shortURL)
	return
}