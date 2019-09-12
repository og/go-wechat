package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	grand "github.com/og/go-rand"
	gwechat "github.com/og/go-wechat"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getMD5(v string) string {
	md5Byte := md5.Sum([]byte(v))
	return fmt.Sprintf("%x", md5Byte)
}

type wechatHook struct {}
func (self wechatHook) GetAccessToken(appID string, appSecret string) (accessToken string , err error) {
	url := "http://localhost:6136/api/wechat/project_name/get_key?type=access_token&hash=" + getMD5(appID + appSecret)
	ressource, err := http.Get(url)
	if err != nil { return "", err }
	if ressource != nil {
		defer ressource.Body.Close()
	}
	body, err := ioutil.ReadAll(ressource.Body) ; if err != nil { return "", err }
	var res gwechat.WecahtCentralControlServiceRes
	err = json.Unmarshal(body, &res) ; if err != nil { return "", err }
	switch res.Type {
	case "pass":
		return res.Data.AccessToken, nil
	case "fail":
		return "", errors.New(res.Msg)
	}
	return
}
func (self wechatHook) GetJSAPITicket(appID string, appSecret string) (ticket string, err error) {
	url := "http://localhost:6136/api/wechat/project_name/get_key?type=jsapi_ticket&hash=" + getMD5(appID + appSecret)
	ressource, err := http.Get(url)
	if err != nil { return "", err }
	if ressource != nil {
		defer ressource.Body.Close()
	}
	body, err := ioutil.ReadAll(ressource.Body) ; if err != nil { return "", err }
	var res gwechat.WecahtCentralControlServiceRes
	err = json.Unmarshal(body, &res) ; if err != nil { return "", err }
	switch res.Type {
	case "pass":
		return res.Data.JSAPITicket, nil
	case "fail":
		return "", errors.New(res.Msg)
	}
	return
}
type JSAPIConfig struct {
	APPID string `json:"appId"`
	Timestamp int64 `json:"timestamp"`
	NonceStr string `json:"nonceStr"`
	Signature string `json:"signature"`
}
func GetJSAPIConfig(url string, jsapiTicket string, appID string) JSAPIConfig {
	jsAPIConfig := JSAPIConfig{
		APPID: appID,
		Timestamp: time.Now().UTC().Unix(),
		NonceStr: grand.StringLetter(20),
	}
	query := "jsapi_ticket=" + jsapiTicket + "&noncestr=" + jsAPIConfig.NonceStr + "&timestamp=" +  strconv.FormatInt(jsAPIConfig.Timestamp, 10) + "&url=" + url
	shaByte := sha1.Sum([]byte(query))
	jsAPIConfig.Signature = fmt.Sprintf("%x", shaByte[:])
	return jsAPIConfig
}
func main () {
	appID := gwechat.TestEnvAPPID // 这里换成你自己的 appid
	appSecret := gwechat.TestEnvAPPSecret // 这里换成你自己的 appSecret
	var wechat = gwechat.New(gwechat.Config{
		APPID: appID,
		APPSecret: appSecret,
		Hook: wechatHook{},
	})
	accessToken, has := wechat.GetAccessToken()
	log.Print("accessToken ", accessToken, has)

	jsapiTicket, has := wechat.GetJSAPITicket()
	log.Print("jsapiTicket ", jsapiTicket, has)
	jsapiConfig := GetJSAPIConfig("https://github.com/og", jsapiTicket, appID)
	log.Print(jsapiConfig)
}
