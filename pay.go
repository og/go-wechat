package gwechat

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	qs "github.com/google/go-querystring/query"
	ge "github.com/og/go-error"
	gjson "github.com/og/go-json"
	gconv "github.com/og/x/conv"
	grand "github.com/og/x/rand"
	"github.com/pkg/errors"
	"io/ioutil"
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
type PayUnifiedOrderData struct {
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
	PayJSAPIConfig PayJSAPIConfig

}
// https://pay.weixin.qq.com/wiki/doc/api/wxa/wxa_api.php?chapter=9_1
// 统一下单
func (this Wechat) PayUnifiedOrder (data PayUnifiedOrderData) (result UnifedOrderResult, payErrRes PayErrRes)  {
	data.MCHID = this.mchID
	data.MCHKey = this.mchKey
	data.NonceStr = grand.StringLetter(32)

	data.Sign = CreatePaySign(data)
	xmlByteList, err := xml.Marshal(data)
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
	result.PayJSAPIConfig = CreatePayClientAPIConfig(data.APPID, result.PrepayID)
	return
}
type PayJSAPIConfig struct {
	APPID string `url:"appId"json:"appId"`
	TimeStamp string `url:"timeStamp"json:"timeStamp"`
	NonceStr string `url:"nonceStr"json:"nonceStr"`
	Package string `url:"package"json:"package"`
	SignType string `url:"signType"json:"signType"`
	PaySign string `url:"paySign"json:"paySign"`

}
// 根据 appID prepayID  创建客户端支付API配置参数
func CreatePayClientAPIConfig(appID string, prepayID string) (config PayJSAPIConfig) {
	config.APPID = appID
	config.TimeStamp = gconv.Int64String(time.Now().Unix())
	config.NonceStr = grand.StringLetter(32)
	config.Package = "prepay_id=" + prepayID
	config.SignType = "MD5"
	clientAPIConfigQS, err := qs.Values(config) ; ge.Check(err)
	md5Byte := md5.Sum([]byte(clientAPIConfigQS.Encode()))
	config.PaySign = strings.ToUpper(fmt.Sprintf("%x", md5Byte))
	return
}
func CreatePaySign(query PayUnifiedOrderData) (sign string) {
	// fuck wechat
	tempNotifyURL := query.NotifyURL
	query.NotifyURL = "xxNotifyURLxx"
	queryString, err := qs.Values(query) ; ge.Check(err)
	urlString := queryString.Encode()
	urlString = strings.ReplaceAll(urlString, query.NotifyURL, tempNotifyURL)
	md5Byte := md5.Sum([]byte(urlString + "&key=" + query.MCHKey))
	sign = strings.ToUpper(fmt.Sprintf("%x", md5Byte))
	return
}

type PayOrderQueryData struct {
	APPID string `xml:"appid,omitempty"url:"appid,omitempty"`
	OutTradeNo 		string `xml:"out_trade_no,omitempty"url:"out_trade_no,omitempty,omitempty";_____note:"商户订单号"`
}
type PayOrderQueryResult struct {
	OutTradeNo string `xml:"out_trade_no"`
	ResultCode string `xml:"result_code"`
	ResultMsg string `xml:"return_msg"`
	ErrCode string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	APPID string `xml:"APPID"`
	MCHID string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign string `xml:"sign"`
	DeviceInfo string `xml:"device_info"`
	OpenID string `xml:"openid"`
	IsSubscribe string `xml:"is_subscribe"`
	TradeType string `xml:"trade_type"`
	TradeState string `xml:"trade_state"`
	BankType string `xml:"bank_type"`
	TotalFee int `xml:"total_fee"`
	SettlementTotalFee int `xml:"settlement_total_fee"`
	CashFee int `xml:"cash_fee"`
	CashFeeType string `xml:"cash_fee_type"`
	CouponFee string `xml:"coupon_fee"`
	CouponCount string `xml:"coupon_count"`
	TransactionID string `xml:"transaction_id"`
	Attach string `xml:"attach"`
	TimeEnd string `xml:"time_end"`
	TradeStateDesc string `xml:"trade_state_desc"`
}
func (self Wechat) PayOrderQuery (data PayOrderQueryData) (result PayOrderQueryResult, payErrRes PayErrRes) {
	query := struct {
		PayOrderQueryData
		XMLName  xml.Name `xml:"xml"url:"-"`
		MCHID string `xml:"mch_id,omitempty"url:"mch_id,omitempty"`
		TransactionID string `xml:"transaction_id,omitempty"url:"transaction_id,omitempty"`
		NonceStr string `xml:"nonce_str,omitempty"url:"nonce_str,omitempty"`
		SignType string `xml:"sign_type,omitempty"url:"sign_type,omitempty"`
		Sign string `xml:"sign,omitempty"url:"sign,omitempty"`
	}{
		PayOrderQueryData: data,
		MCHID: self.mchID,
		NonceStr: grand.StringLetter(32),
		SignType: "MD5",
	}
	urlString, err := qs.Values(query) ;ge.Check(err)
	md5Byte := md5.Sum([]byte(urlString.Encode() + "&key=" + self.mchKey))
	sign := strings.ToUpper(fmt.Sprintf("%x", md5Byte))
	query.Sign = sign
	xmlByteList ,err := xml.Marshal(query); if err !=nil { payErrRes.SetError(err) ; return }

	client := &http.Client{}
	requestBody :=  bytes.NewReader([]byte{})
	requestBody = bytes.NewReader(xmlByteList)
	request, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/pay/orderquery", requestBody)
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
		payErrRes.SetError(errors.New(gjson.StringUnfold(result)))
		return
	}
	if result.ResultCode != "SUCCESS" {
		payErrRes.SetError(errors.New(gjson.StringUnfold(result)))
		return
	}
	return
}