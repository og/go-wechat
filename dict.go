package gwechat

import "github.com/og/go-dict"
type dictStruct struct {
	PayUnifiedorder struct{
		DeviceInfo struct{
			WEB string `dict:"WEB"`
		}
		SignType struct{
			MD5 string `dict:"MD5"`
			HMAC string `dict:"HMAC"`
			SHA256 string `dict:"SHA256"`
		}
		FeeType struct{
			CNY string `dict:"CNY"`
		}
		TradeType struct{
			JSAPI string `dict:"JSAPI";note:"jsapi支付或小程序支付"`
			NATIVE string `dict:"NATIVE"`
			APP string `dict:"APP"`
			MWEB string `dict:"MWEB";note:"H5支付"`
			MICROPAY string `dict:"MICROPAY";note:"付款码支付"`
		}
	}
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
