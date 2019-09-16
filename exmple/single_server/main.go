package main

import (
	gwechat "github.com/og/go-wechat"
	"log"
)


var namePublicAccount = gwechat.New(gwechat.Config{
	APPID: gwechat.TestEnvAPPID, // 此处替换你的 appid
	APPSecret: gwechat.TestEnvAPPSecret, // 此处替换你的 app secret
	CenterService: namePublicAccountCenterService{},
})
var wechatMemeberCache = gwechat.MemoryCache{}
type namePublicAccountCenterService struct {}
func (self namePublicAccountCenterService) GetAccessToken (appID string, appSecret string) (accessToken string , errRes gwechat.ErrResponse){
	return gwechat.UnsafeGetAccessToken(appID, appSecret, &wechatMemeberCache)
}

func (self namePublicAccountCenterService) GetJSAPITicket(appID string, appSecret string) (ticket string, errRes gwechat.ErrResponse){
	accessToken , errRes := self.GetAccessToken(appID, appSecret)
	if errRes.Fail {return "", errRes}
	return gwechat.UnsafeGetJSAPITicket(appID, accessToken, &wechatMemeberCache)
}

func main () {
	shortURL, errRes := namePublicAccount.GetShortURL("https://github.com/og")
	if errRes.Fail {
		panic(errRes.Error)
	}
	log.Print(shortURL)
}
