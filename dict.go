package gwechat

import "github.com/og/go-dict"
type dictStruct struct {
	PayOrderQuery struct{
		TradeState struct {
			SUCCESS    string `dict:"SUCCESS"note:"支付成功"`
			REFUND     string `dict:"REFUND"note:"转入退款"`
			NOTPAY     string `dict:"NOTPAY"note:"未支付"`
			CLOSED     string `dict:"CLOSED"note:"已关闭"`
			REVOKED    string `dict:"REVOKED"note:"已撤销（刷卡支付）"`
			USERPAYING string `dict:"USERPAYING"note:"用户支付中"`
			PAYERROR   string `dict:"PAYERROR"note:"支付失败(其他原因，如银行返回失败)"`
		}
	}
	PayUnifiedOrder struct{
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
