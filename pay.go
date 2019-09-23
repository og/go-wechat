package gwechat

import (
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	ge "github.com/og/go-error"
	gconv "github.com/og/x/conv"
	glist "github.com/og/x/list"
	gmap "github.com/og/x/map"
	grand "github.com/og/x/rand"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
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
	x string `xml:"-";note:"可选参数开始"`
	APPID 			string `xml:"appid";____________note:"不读取config 中的appid，因为一个商户号可以向不同的小程序进行支付"`
	OutTradeNo 		string `xml:"out_trade_no";_____note:"商户订单号"`
	TotalFee 		int    `xml:"total_fee";________note:"支付金额。注意是以分为单位"`
	Body 			string `xml:"body";_____________note:"商品描述"`
	SpbillCreateIP 	string `xml:"spbill_create_ip";_note:"支付客户端IP"`
	NotifyURL	    string `xml:"notify_url";_______note:"通知地址"`
	TradeType 		string `xml:"trade_type";_______note:"交易类型"`
	OpenID			string `xml:"openid";___________note:"trade_type=JSAPI 时候必传"`
	xx string `xml:"-";note:"可选参数结束"`
	DeviceInfo string `xml:"device_info"`
	Attach string `xml:"attach"`
	TimeStart string `xml:"time_start"`
	TimeExpire string `xml:"time_expire"`
	GoodsTag string `xml:"goods_tag"`
	ProductID string `xml:"product_id"`
	LimitPay string `xml:"limit_pay"`
	Receipt string `xml:"receipt"`
	Detail string `xml:"detail"`
	XMLName  xml.Name `xml:"xml"`
	FeeType string `xml:"fee_type"`
	xxx string `xml:"-";note:"内部属性不要配置 Start"`
	MCHID string `xml:"mch_id"`
	MCHKey string `xml:"-"`
	NonceStr string `xml:"nonce_str"`
	Sign string `xml:"sign"`
	xxxx string `xml:"-";note:"内部属性不要配置 End"`

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
	return
}

func CreatePaySign(query PayUnifiedOrderQuery) (sign string) {
	queryValues := url.Values{}
	value := reflect.ValueOf(query)
	valueType := reflect.TypeOf(query)
	itemLen := value.NumField()
	for i:=0;i<itemLen;i++ {
		queryKey := reflect.StructTag(valueType.Field(i).Tag).Get("xml")
		switch queryKey{
			case "-":
				continue
			case "":
				panic("need tag name:" + valueType.Field(i).Name)
		}
		valueItem := value.Field(i)
		queryValue := ""
		switch valueItem.Type().String() {
		case "string":
			queryValue = valueItem.String()
		case "int":
			queryValue = gconv.Int64String(valueItem.Int())
		case "float64":
			queryValue = fmt.Sprintf("%f", valueItem.Float())
		case "xml.Name":
			continue
		default:
			panic("error type:" + valueItem.Type().String())
		}
		if queryValue != "" {
			queryValues.Set(queryKey, queryValue)
		}
	}
	keyList := gmap.Keys(queryValues).String()
	queryStringList := glist.StringList{}
	for i:=0;i<len(keyList);i++ {
		key := keyList[i]
		queryStringList.Push(key + "=" + queryValues.Get(key))
	}
	md5Byte := md5.Sum([]byte(queryStringList.Join("&") + "&key=" + query.MCHKey))
	sign = strings.ToUpper(fmt.Sprintf("%x", md5Byte))
	return
}