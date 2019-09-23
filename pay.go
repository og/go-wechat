package gwechat

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	qs "github.com/google/go-querystring/query"
	ge "github.com/og/go-error"
	gconv "github.com/og/x/conv"
	grand "github.com/og/x/rand"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type PayErrBase struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg string `xml:"return_msg"`
}
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
type PayUnifiedOrderQuery struct {
	xxxxxxxxxxxx string `xml:"-";note:"必传参数开始"`
	APPID 			string `xml:"appid"url:"appid,omitempty";____________note:"不读取config 中的appid，因为一个商户号可以向不同的小程序进行支付"`
	OutTradeNo 		string `xml:"out_trade_no"url:"out_trade_no,omitempty";_____note:"商户订单号"`
	TotalFee 		int    `xml:"total_fee"url:"total_fee,omitempty";________note:"支付金额。注意是以分为单位"`
	Body 			string `xml:"body"url:"body,omitempty";_____________note:"商品描述"`
	SpbillCreateIP 	string `xml:"spbill_create_ip"url:"spbill_create_ip,omitempty";_note:"支付客户端IP"`
	NotifyURL	    string `xml:"notify_url"url:"notify_url,omitempty";_______note:"通知地址"`
	TradeType 		string `xml:"trade_type"url:"trade_type,omitempty";_______note:"交易类型"`
	OpenID			string `xml:"openid"url:"openid,omitempty";___________note:"trade_type=JSAPI 时候必传"`
	xxxxxxxxxxxxxx string `xml:"-";note:"可选参数"`
	DeviceInfo string `xml:"device_info"url:"device_info,omitempty"`
	Attach string `xml:"attach"url:"attach,omitempty"`
	TimeStart string `xml:"time_start"url:"time_start,omitempty"`
	TimeExpire string `xml:"time_expire"url:"time_expire,omitempty"`
	GoodsTag string `xml:"goods_tag"url:"goods_tag,omitempty"`
	ProductID string `xml:"product_id"url:"product_id,omitempty"`
	LimitPay string `xml:"limit_pay"url:"limit_pay,omitempty"`
	Receipt string `xml:"receipt"url:"receipt,omitempty"`
	Detail string `xml:"detail"url:"detail,omitempty"`
	XMLName  xml.Name `xml:"xml"url:"-"`
	FeeType string `xml:"fee_type"url:"fee_type,omitempty"`
	xxxxxxxxxxxxxxxxxxxxx string `xml:"-";note:"内部属性不要设置"`
	MCHID string `xml:"mch_id"url:"mch_id,omitempty"`
	MCHKey string `xml:"-"url:"-"`
	NonceStr string `xml:"nonce_str"url:"nonce_str,omitempty"`
	Sign string `xml:"sign"url:"-"`

}

type UnifedOrderResult struct {
	ResultCode string `xml:"result_code"`
	ResultMsg string `xml:"return_msg"`
	CodeURL string `xml:"code_url"`
	PrepayID string `xml:"prepay_id"`
	TradeType string `xml:"trade_type"`
	ErrCode string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	APPID string `xml:"appid"`
	MCHID string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr string `xml:"nonce_str"`
	Sign string `xml:"sign"`
	ClientApiConfig struct {
		APPID string `url:"appId"json:"appId"`
		TimeStamp string `url:"timeStamp"json:"timeStamp"`
		NonceStr string `url:"nonceStr"json:"nonceStr"`
		Package string `url:"package"json:"package"`
		SignType string `url:"signType"json:"signType"`
		Sign string `url:"sign"json:"sign"`
	}

}
// https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_1
// 统一下单
func (this Wechat) PayUnifiedOrder (query PayUnifiedOrderQuery) (result UnifedOrderResult, payErrRes PayErrRes)  {
	query.MCHID = this.mchID
	query.MCHKey = this.mchKey
	query.NonceStr = grand.StringLetter(32)

	query.Sign = CreatePaySign(query)
	xmlByteList, err := xml.Marshal(query)
	if err !=nil {
		payErrRes.SetError(err)
		return
	}
	client := &http.Client{}
	requestBody :=  bytes.NewReader([]byte{})
	requestBody = bytes.NewReader(xmlByteList)
	request, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/unifiedorder", requestBody)
	if err !=nil {
		payErrRes.SetError(err)
		return
	}
	request.Header.Set("Content-Type", "text/xml")
	response, err := client.Do(request) ; if err !=nil { payErrRes.SetError(err) ; return }
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body) ; if err !=nil { payErrRes.SetError(err) ; return }
	err = xml.Unmarshal(body, &result) ; ge.Check(err)
	if result.ResultCode != "SUCCESS" {
		payErrRes.SetError(errors.New(result.ResultMsg))
		return
	}
	if result.ResultCode != "SUCCESS" {
		payErrRes.SetError(errors.New(result.ErrCodeDes))
		return
	}
	result.ClientApiConfig.APPID = query.APPID
	result.ClientApiConfig.TimeStamp = gconv.Int64String(time.Now().Unix())
	result.ClientApiConfig.NonceStr = grand.StringLetter(32)
	result.ClientApiConfig.Package = "prepay_id=" + result.PrepayID
	result.ClientApiConfig.SignType = "MD5"
	clientAPIConfigQS, err := qs.Values(result.ClientApiConfig) ; ge.Check(err)
	md5Byte := md5.Sum([]byte(clientAPIConfigQS.Encode()))
	result.ClientApiConfig.Sign = strings.ToUpper(fmt.Sprintf("%x", md5Byte))
	return
}

func CreatePaySign(query PayUnifiedOrderQuery) (sign string) {
	// fuck wechat
	tempNotifyURL := query.NotifyURL
	query.NotifyURL = "xxNotifyURLxx"
	queryString, err := qs.Values(query) ; ge.Check(err)
	urlString := queryString.Encode()
	urlString = strings.ReplaceAll(urlString, query.NotifyURL, tempNotifyURL)
	log.Print(urlString)
	md5Byte := md5.Sum([]byte(urlString + "&key=" + query.MCHKey))
	sign = strings.ToUpper(fmt.Sprintf("%x", md5Byte))
	return
}