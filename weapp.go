package gwechat

import (
	"encoding/json"
	qs "github.com/google/go-querystring/query"
	ge "github.com/og/go-error"
	"io/ioutil"
	"net/http"
)

// 微信小程序登录
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html
// 本函数没有返回 unionid 因为根据场景不同微信不一定会返回 unionid
func (this Wechat) WeappCode2Session(code string) (res WeappCode2SessionResponse, errRes ErrResponse) {
	query := struct {
		Code string `url:"js_code"`
		APPID string `url:"appid"`
		Secret string `url:"secret"`
		GrantType string `url:"grant_type"`
	}{
		Code:      code,
		APPID:     this.appID,
		Secret:    this.appSecret,
		GrantType: "authorization_code",
	}
	querystring, err := qs.Values(query); ge.Check(err)
	requestURL := apiDomain + "/sns/jscode2session" + "?" + querystring.Encode()
	resp, err := http.Get(requestURL); ge.Check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body); ge.Check(err)
	err = json.Unmarshal(body, &res); ge.Check(err)
	if res.ErrCode != 0 {
		errRes.SetFail(res.ErrCode, res.ErrMsg)
		return
	}
	return
}
