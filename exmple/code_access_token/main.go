package main

import (
	ge "github.com/og/go-error"
	gjson "github.com/og/go-json"
	gwechat "github.com/og/go-wechat"
	"github.com/pkg/errors"
	qrcode "github.com/skip2/go-qrcode"
	"log"
	"net/http"
)

var wechat = gwechat.New(gwechat.Config{
	APPID: gwechat.EnvAPPID,
	APPSecret: gwechat.EnvAPPSecret,
	Hook: wechatHook{},
})

type wechatHook struct {}
var accessTokenMemoryCache = &gwechat.AccessTokenMemoryCache{}
func (self wechatHook) GetAccessToken(appID string, appSecret string) (accessToken string , err error) {
	accessToken, errRes := gwechat.UnsafeGetAccessToken(appID, appSecret, accessTokenMemoryCache)
	if errRes.Fail { return  "", errors.New(errRes.ErrMsg) }
	return
}

const port = "7315"
// 这里可以替换成你自己的域名
const WechatAuthDomain = "http://www.admpv.com"
func main () {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm() ; ge.Check(err)
		scope := r.Form.Get("scope")
		if scope == "" {
			_, _ = w.Write([]byte("请带上 scope 参数"))
			return
		}
		redirectURL := wechat.WebRedirectAuthorize(scope, WechatAuthDomain + "/get_access_token", "CustomState:"+ scope)
		// var photo []byte
		// photo, err := qrcode.Encode(redirectURL, qrcode.Medium, 256) ; ge.Check(err)
		// _, _ = w.Write(photo)
		_, _ = w.Write([]byte(redirectURL + "\r\n 复制上面的链接在微信开发者工具中打开跳转后拿到code \r\n 然后将域名替换为 http://localhost:" + port))
		// 正式环境请使用 redirect
		// http.Redirect(w, r, redirectURL ,http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/get_access_token", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm() ; ge.Check(err)
		code := r.Form.Get("code")
		state := r.Form.Get("state")
		_ = state
		accessTokenWechatRes, errRes := wechat.WebAccessToken(code)
		if errRes.Fail {
			_, _ = w.Write(gjson.Byte(gjson.FailMsg(errRes.ErrMsg)))
			return
		}
		var webGetUserInfoRes gwechat.WebUserInfoResponse
		// 进一步授权获取信息
		if accessTokenWechatRes.Scope == gwechat.Dict().WebRedirectAuthorize.Scope.SnsapiUserinfo {
			webGetUserInfoRes, errRes = wechat.WebGetUserInfo(accessTokenWechatRes.AccessToken, accessTokenWechatRes.OpenID, gwechat.Dict().WebUserInfo.Lang.ZHCN)
			if errRes.Fail {
				log.Print("第四步：拉取用户信息(需scope为 snsapi_userinfo) 错误")
				_, _ = w.Write(gjson.Byte(gjson.FailMsg(errRes.ErrMsg)))
				return
			}
		}
		_, _ = w.Write(gjson.Byte(map[string]interface{}{
			"WebAccessToken": accessTokenWechatRes,
			"WebGetUserInfo": webGetUserInfoRes,
		}))
	})

	localDomain := "http://localhost:" + port
	log.Print("打开 " + localDomain + "/?scope=snsapi_base")
	log.Print("或者")
	log.Print("打开 " + localDomain + "/?scope=snsapi_userinfo")
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Print(err)
	}
}
