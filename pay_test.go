package gwechat_test

import (
	gjson "github.com/og/go-json"
	grand "github.com/og/go-rand"
	gwechat "github.com/og/go-wechat"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var someMCH = gwechat.New(gwechat.Config{
	MCHID: gwechat.TestMCHID,
	MCHKey: gwechat.TestMCHKey,
})

func TestWechat_PayUnifiedorder(t *testing.T) {
	payResult, payErrRes := someMCH.PayUnifiedOrder(gwechat.PayUnifiedOrderQuery{
		APPID: "wx9f01246e31fd5cae",
		OutTradeNo: grand.StringLetter(32),
		TotalFee: 1,
		Body: "test",
		SpbillCreateIP: "218.81.205.86", // 真实环境务必使用用户客户端ip
		NotifyURL: "https://github.com/og/go-wechat",
		TradeType: gwechat.Dict().PayUnifiedorder.TradeType.JSAPI,
		OpenID: "otU685kfo1D-tAbgaNVqbPiiRd8k",
	})
	if payErrRes.Fail {
		panic(payErrRes.Error)
	}
	log.Print(gjson.StringUnfold(payResult))
}
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=20_1
func TestCreatePaySign(t *testing.T) {

	assert.Equal(t, `40DE1DF376357F4F657DBCCACF73223F`, gwechat.CreatePaySign(gwechat.PayUnifiedOrderQuery{
		APPID: "wx9111246221fd5cae",
		Body: "body",
		DeviceInfo: "info",
		MCHID: "1112225871",
		NonceStr: "SrwzbDPKqTOsHkVxzjtTXOnrMrPedtcO",
	}))
	assert.Equal(t, `048DEBED0793A789B6D3AA4D002170D5`, gwechat.CreatePaySign(gwechat.PayUnifiedOrderQuery{
		APPID: "wx9111246221fd5cae",
		Body: "body",
		DeviceInfo: "info",
		MCHID: "1112225871",
		NonceStr: "SrwzbDPKqTOsHkVxzjtTXOnrMrPedtcO",
		NotifyURL: "http://github.com/og/go-wechat",
	}))

	assert.Equal(t, `687EC6E0AFCDBFB825E05B4EA83103BF`, gwechat.CreatePaySign(gwechat.PayUnifiedOrderQuery{
		APPID: "wx9111246221fd5cae",
		MCHID: "1112225871",
		MCHKey: "wxda123461113e2357wxdac12521113e",
		Body: "body",
		NonceStr: "SrwzbDPKqTOsHkVxzjtTXOnrMrPedtcO",
	}))
}