package main

import (
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


var accessTokenMemoryCache = &gwechat.AccessTokenMemoryCache{}
func GetAccessTokenCtrl (w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil { http.Error(w, fmt.Sprint(err), http.StatusInternalServerError); return }
	}()
	var res gwechat.GetAccessTokenRes
	res.Type = "pass"
	// 这里换成你自己的 appid 和 appSecret
	accessToken, errRes := gwechat.UnsafeGetAccessToken(gwechat.EnvAPPID, gwechat.EnvAPPSecret, accessTokenMemoryCache)
	if errRes.Fail {
		res.Type = "fail"
		jsonb, err :=json.Marshal(errRes) ; if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
		res.Msg = string(jsonb)
	} else {
		res.Type = "pass"
		res.Data.AccessToken = accessToken
	}
	jsonb, err := json.Marshal(&res) ; if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonb)
}