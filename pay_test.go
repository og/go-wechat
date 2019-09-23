package gwechat_test

import (
	gjson "github.com/og/go-json"
	grand "github.com/og/go-rand"
	gwechat "github.com/og/go-wechat"
	l "github.com/og/x/log"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var someMCH = gwechat.New(gwechat.Config{
	MCHID: gwechat.TestMCHID,
	MCHKey: gwechat.TestMCHKey,
})

func TestWechat_PayUnifiedorder(t *testing.T) {
	testOutTradeNo := grand.StringLetter(32)
	payResult, payErrRes := someMCH.PayUnifiedOrder(gwechat.PayUnifiedOrderData{
		APPID:  gwechat.TestEnvAPPID,
		OutTradeNo: testOutTradeNo,
		TotalFee: 1,
		Body: "test",
		SpbillCreateIP: "218.81.205.86", // 真实环境务必使用用户客户端ip
		NotifyURL: "https://github.com/og",
		TradeType: gwechat.Dict().PayUnifiedOrder.TradeType.JSAPI,
		OpenID: "ovv4N5I8YJ82fTpBf_JOz-3-ljnE",
	})
	if payErrRes.Fail {
		log.Print(payResult)
		panic(payErrRes.Error)
	}
	log.Print("testOutTradeNo ", testOutTradeNo)
	log.Print(gjson.StringUnfold(payResult))
}
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=20_1
func TestCreatePaySign(t *testing.T) {

	assert.Equal(t, `40DE1DF376357F4F657DBCCACF73223F`, gwechat.CreatePaySign(gwechat.PayUnifiedOrderData{
		APPID: gwechat.TestEnvAPPID,
		Body: "body",
		DeviceInfo: "info",
		MCHID: "1112225871",
		NonceStr: "SrwzbDPKqTOsHkVxzjtTXOnrMrPedtcO",
	}))
	assert.Equal(t, `048DEBED0793A789B6D3AA4D002170D5`, gwechat.CreatePaySign(gwechat.PayUnifiedOrderData{
		APPID: gwechat.TestEnvAPPID,
		Body: "body",
		DeviceInfo: "info",
		MCHID: "1112225871",
		NonceStr: "SrwzbDPKqTOsHkVxzjtTXOnrMrPedtcO",
		NotifyURL: "http://github.com/og/go-wechat",
	}))

	assert.Equal(t, `687EC6E0AFCDBFB825E05B4EA83103BF`, gwechat.CreatePaySign(gwechat.PayUnifiedOrderData{
		APPID: gwechat.TestEnvAPPID,
		MCHID: "1112225871",
		MCHKey: "wxda123461113e2357wxdac12521113e",
		Body: "body",
		NonceStr: "SrwzbDPKqTOsHkVxzjtTXOnrMrPedtcO",
	}))
}
func TestWechat_PayOrderQuery(t *testing.T) {
	result, payErrRes := someMCH.PayOrderQuery(gwechat.PayOrderQueryData{
		APPID: gwechat.TestEnvAPPID,
		OutTradeNo: "siJXeNJakGeauZOvGYgusVZBirpBAxOb",
	})
	if payErrRes.Fail {
		panic(payErrRes.Error)
	}
	l.V(gjson.StringUnfold(result))
}