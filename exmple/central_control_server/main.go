package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	ge "github.com/og/go-error"
	gwechat "github.com/og/go-wechat"
	"log"
	"net/http"
)




func main () {
	getKeyURL := "/api/wechat/project_name/get_key"
	http.HandleFunc(getKeyURL, GetKeyCtrl)
	port := "6136"
	log.Print("open http://localhost:"+ port + getKeyURL + "?hash=md5(appID+appSecret)&type=access_token")
	log.Print("open http://localhost:"+ port + getKeyURL + "?hash=md5(appID+appSecret)&type=jsapi_ticket")
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
var wechatMemoryCache = &gwechat.MemoryCache{}
type GetKeyQuery struct {
	Hash string `url:"hash"`
	Type string `url:"type";enum:"[]string{access_token,jsapi_ticket}"`
}
func GetKeyCtrl (w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil { http.Error(w, fmt.Sprint(err), http.StatusInternalServerError); return }
	}()
	// 这里换成你自己的 appid 和 appSecret
	appID := gwechat.EnvAPPID
	appSecret := gwechat.EnvAPPSecret
	correctHash := getMD5(appID + appSecret)
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm() ; ge.Check(err)
	query := GetKeyQuery{
		Hash: r.Form.Get("hash"),
		Type: r.Form.Get("type"),
	}

	res := gwechat.WecahtCentralControlServiceRes{}
	res.Type = "pass"
	res.Data.Type = query.Type
	if query.Hash != correctHash {
		res.Type = "fail"
		res.Msg = "error hash (hash = md5(appID + appSecret))"
		jsonb, err := json.Marshal(&res) ; if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
		_, _ = w.Write(jsonb)
		return
	}
	accessToken, errRes := gwechat.UnsafeGetAccessToken(appID, appSecret, wechatMemoryCache)
	if errRes.Fail {
		res.Type = "fail"
		jsonb, err :=json.Marshal(errRes) ; if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
		res.Msg = string(jsonb)
	} else {
		res.Data.AccessToken = accessToken
	}
	switch query.Type {
	case "access_token":
	case "jsapi_ticket":
		jsapiTicket, errRes := gwechat.UnsafeGetJSAPITicket(appID, res.Data.AccessToken, wechatMemoryCache)
		if errRes.Fail {
			http.Error(w, errRes.ErrMsg, http.StatusInternalServerError)
			return
		}
		res.Data.JSAPITicket = jsapiTicket
	default:
		http.Error(w, "query type error", http.StatusInternalServerError)
		return
	}
	res.Type = "pass"
	jsonb, err := json.Marshal(&res) ; if err != nil { http.Error(w, err.Error(), http.StatusInternalServerError); return }
	_, _ = w.Write(jsonb)
}