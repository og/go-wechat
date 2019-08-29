package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	gwechat "github.com/og/go-wechat"
	"log"
	"net/http"
)

func main () {
	apiWechatGetAccessTokenURL := "/api/wechat/get_access_token"
	http.HandleFunc(apiWechatGetAccessTokenURL, GetAccessTokenCtrl)
	port := "6136"
	log.Print("Listen: http://localhost:"+ port + apiWechatGetAccessTokenURL)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Print(err)
	}
}

func getMD5(v string) string {
	md5Byte := md5.Sum([]byte(v))
	return fmt.Sprintf("%x", md5Byte)
}

// 内存缓存必须放在包级别变量中
var accessTokenMemoryCache = &gwechat.AccessTokenMemoryCache{}
func GetAccessTokenCtrl (w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil { http.Error(w, fmt.Sprint(err), http.StatusInternalServerError); return }
	}()
	// 这里换成你自己的 appid 和 appSecret
	appID := gwechat.EnvAPPID
	appSecret := gwechat.EnvAPPSecret
	correctHash := getMD5(appID + appSecret)
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	queryHash := r.Form.Get("hash")
	res := gwechat.GetAccessTokenRes{
		Type: "pass",
	}
	if queryHash != correctHash {
		res.Type = "fail"
		res.Msg = "error hash (hash = md5(appID + appSecret))"
		jsonb, err := json.Marshal(&res) ; if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
		_, _ = w.Write(jsonb)
		return
	}
	accessToken, errRes := gwechat.UnsafeGetAccessToken(appID, appSecret, accessTokenMemoryCache)
	if errRes.Fail {
		res.Type = "fail"
		jsonb, err :=json.Marshal(errRes) ; if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
		res.Msg = string(jsonb)
	} else {
		res.Type = "pass"
		res.Data.AccessToken = accessToken
	}
	jsonb, err := json.Marshal(&res) ; if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
	_, _ = w.Write(jsonb)
}