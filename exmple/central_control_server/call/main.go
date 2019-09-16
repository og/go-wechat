package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	gjson "github.com/og/go-json"
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

type projectNameWechatcenterService struct {}
func (self projectNameWechatcenterService) GetAccessToken(appID string, appSecret string) (accessToken string , errRes gwechat.ErrResponse) {
	url := "http://localhost:6136/api/wechat/project_name/get_key?type=access_token&hash=" + getMD5(appID + appSecret)
	ressource, err := http.Get(url)
	if err != nil { errRes.SetSystemError(err); return "",errRes  }
	defer ressource.Body.Close()
	body, err := ioutil.ReadAll(ressource.Body) ; if err != nil { errRes.SetSystemError(err); return "", errRes }
	var res gwechat.WecahtCentralControlServiceRes
	err = json.Unmarshal(body, &res) ; if err != nil { errRes.SetSystemError(err); return "", errRes  }
	switch res.Type {
	case "pass":
		return res.Data.AccessToken, errRes
	case "fail":
		errRes.SetSystemError(errors.New(res.Msg))
		return "", errRes
	}
	return
}
func (self projectNameWechatcenterService) GetJSAPITicket(appID string, appSecret string) (ticket string, errRes gwechat.ErrResponse) {
	url := "http://localhost:6136/api/wechat/project_name/get_key?type=jsapi_ticket&hash=" + getMD5(appID + appSecret)
	ressource, err := http.Get(url)
	if err != nil { errRes.SetSystemError(err) ; return "", errRes }
	if ressource != nil {
		defer ressource.Body.Close()
	}
	body, err := ioutil.ReadAll(ressource.Body) ; if err != nil { errRes.SetSystemError(err); return "", errRes  }
	var res gwechat.WecahtCentralControlServiceRes
	err = json.Unmarshal(body, &res) ; if err != nil { errRes.SetSystemError(err) ; return "", errRes }
	switch res.Type {
	case "pass":
		return res.Data.JSAPITicket, errRes
	case "fail":
		errRes.SetSystemError(errors.New(res.Msg))
		return "", errRes
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
		CenterService: projectNameWechatcenterService{},
	})
	accessToken, errRes := wechat.GetAccessToken()
	log.Print("GetAccessToken ", accessToken, errRes)

	jsapiTicket, errRes := wechat.GetJSAPITicket()
	log.Print("GetJSAPITicket ", jsapiTicket, errRes)
	jsapiConfig := GetJSAPIConfig("https://github.com/og", jsapiTicket, appID)
	log.Print("GetJSAPIConfig ", gjson.StringUnfold(jsapiConfig))
}
