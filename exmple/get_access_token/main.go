package main

import (
	"encoding/json"
	gwechat "github.com/og/go-wechat"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

type wechatHook struct {}
func (self wechatHook) GetAccessToken(appID string) (accessToken string , err error) {
	ressource, err := http.Get("http://localhost:6136/api/wechat/get_access_token?appid=" + appID)
	if err != nil { return "", err }
	if ressource != nil {
		defer ressource.Body.Close()
	}
	body, err := ioutil.ReadAll(ressource.Body) ; if err != nil { return "", err }
	var res gwechat.GetAccessTokenRes
	err = json.Unmarshal(body, &res) ; if err != nil { return "", err }
	switch res.Type {
	case "pass":
		return res.Data.AccessToken, nil
	case "fail":
		return "", errors.New(res.Msg)
	}
	return
}
func main () {
	var wechat = gwechat.New(gwechat.Config{
		APPID: gwechat.EnvAPPID, // 这里换成你自己的 appid
		APPSecret: gwechat.EnvAPPSecret, // 这里换成你自己的 appSecret
		Hook: wechatHook{},
	})
	accessToken, has := wechat.GetAccessToken()
	log.Print(accessToken, has)
}
