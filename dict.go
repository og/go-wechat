package gwechat

import "github.com/og/go-dict"
type dictStruct struct {
	WebRedirectAuthorize struct{
		Scope struct{
			SnsapiBase string `dict:"snsapi_base"`
			SnsapiUserinfo string `dict:"snsapi_userinfo"`
		}
	}
	WebUserInfo struct{
		Lang struct{
			ZHCN string `dict:"zh_CN"`
			ZHTW string `dict:"zh_TW"`
			EN string `dict:"en"`
		}
	}
}
var dict = dictStruct{}
func init () {
	gdict.Fill(&dict)
}
func Dict() dictStruct {
	return dict
}
