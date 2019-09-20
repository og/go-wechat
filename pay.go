package gwechat

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	grand "github.com/og/go-rand"
	"log"
	"net/url"
	"strings"
)

type PayErrRes struct {
	Fail bool
	Msg string
	Error error
}
func (payErrRes *PayErrRes) SetError(err error) {
	payErrRes.Fail = true
	payErrRes.Msg = err.Error()
	payErrRes.Error = err
}
type PayUnifiedorderMustData struct {
	APPID 			string `xml:"appid";____________note:"不读取config 中的appid，因为一个商户号可以向不同的小程序进行支付"`
	OutTradeNo 		string `xml:"out_trade_no";_____note:"商户订单号"`
	TotalFee 		int    `xml:"total_fee";________note:"支付金额。注意是以分为单位"`
	Body 			string `xml:"body";_____________note:"商品描述"`
	SpbillCreateIP 	string `xml:"spbill_create_ip";_note:"支付客户端IP"`
	NotifyURL	    string `xml:"notify_url";_______note:"通知地址"`
	TradeType 		string `xml:"trade_type";_______note:"交易类型"`
	OpenID			string `xml:"openid";___________note:"trade_type=JSAPI 时候必传"`
}
type PayUnifiedorderOptionalData struct {
	DeviceInfo string `xml:"device_info"`
	Attach string `xml:"attach"`
	TimeStart string `xml:"time_start"`
	TimeExpire string `xml:"time_expire"`
	GoodsTag string `xml:"goods_tag"`
	ProductID string `xml:"product_id"`
	LimitPay string `xml:"limit_pay"`
	Receipt string `xml:"receipt"`
	Detail string `xml:"detail"`
	SceneInfo struct{
		ID string `xml:"id"`
		Name string `xml:"name"`
		AreaCode string `xml:"area_code"`
		Address string `xml:"address"`
	}
}
// https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_1
// 统一下单
func (this Wechat) PayUnifiedorder (mustData PayUnifiedorderMustData, optionalData PayUnifiedorderOptionalData) (payErrRes PayErrRes)  {

	query := struct {
		XMLName  xml.Name `xml:"xml"`
		PayUnifiedorderMustData
		PayUnifiedorderOptionalData
		MCHID string `xml:"mch_id"`
		NonceStr string `xml:"nonce_str"`
		SignType string `xml:"sign_type"`
		Sign string `xml:"sign"`
		FeeType string `xml:"fee_type"`
	}{
		PayUnifiedorderMustData: mustData,
		PayUnifiedorderOptionalData: optionalData,
		MCHID: this.mchID,
		NonceStr: grand.StringLetter(32),
		SignType: dict.PayUnifiedorder.SignType.MD5,
		FeeType: dict.PayUnifiedorder.FeeType.CNY,
	}
	query.Sign = CreatePaySign(PaySignSource{
		APPID: query.APPID,
		MCHID: this.mchID,
		MCHKey: this.mchKey,
		Body: query.Body,
		NonceStr: query.NonceStr,
		DeviceInfo: query.DeviceInfo,
	})

	xmlByteList, err := xml.Marshal(query)
	if err !=nil {
		payErrRes.SetError(err)
		return
	}
	xmlString := string(xmlByteList[:])
	log.Print("query.Sign ", query.Sign)
	_=xmlString
	return
}

type PaySignSource struct {
	APPID string
	MCHID string
	MCHKey string
	Body string
	NonceStr string
	DeviceInfo string
}
func CreatePaySign(data PaySignSource) (sign string) {
	signURLValues := url.Values{
		"appid": {data.APPID},
		"mch_id": {data.MCHID},
		"body": {data.Body},
		"nonce_str": {data.NonceStr},
	}
	if data.DeviceInfo != "" {
		signURLValues.Set("device_info", data.DeviceInfo)
	}
	md5Byte := md5.Sum([]byte(signURLValues.Encode() + "&key=" + data.MCHKey))
	sign = strings.ToUpper(fmt.Sprintf("%x", md5Byte))
	return
}