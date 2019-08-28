package gwechat

import "github.com/og/go-dict"
type dictStruct struct {
	WebRedirectAuthorize struct{
		Scope struct{
			SnsapiBase string `dict:"snsapi_base"`
			SnsapiUserinfo string `dict:"snsapi_userinfo"`
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
