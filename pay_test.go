package gwechat_test

import (
	grand "github.com/og/go-rand"
	gwechat "github.com/og/go-wechat"
	"github.com/stretchr/testify/assert"
	"testing"
)

var someMCH = gwechat.New(gwechat.Config{
	MCHID: gwechat.TestMCHID,
	MCHKey: gwechat.TestMCHKey,
})

func TestWechat_PayUnifiedorder(t *testing.T) {

	payErrRes := someMCH.PayUnifiedorder(gwechat.PayUnifiedorderMustData{
		APPID: "wx9f01246e31fd5cae",
		OutTradeNo: grand.StringLetter(32),
		TotalFee: 1,
		Body: "test",
		SpbillCreateIP: "218.81.205.86", // 真实环境务必使用用户客户端ip
		NotifyURL: "https://github.com/og/go-wechat",
		TradeType: gwechat.Dict().PayUnifiedorder.TradeType.JSAPI,
		OpenID: "otU685kfo1D-tAbgaNVqbPiiRd8k",
	}, gwechat.PayUnifiedorderOptionalData{})
	if payErrRes.Fail {
		panic(payErrRes.Error)
	}
}
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=20_1
func TestCreatePaySign(t *testing.T) {
	assert.Equal(t, `F7F63040EE22E4D84905EB3CF9B35C7E`, gwechat.CreatePaySign(gwechat.PaySignSource{
		APPID: "wx9111246221fd5cae",
		MCHID: "1112225871",
		MCHKey: "wxda123461113e2357wxdac12521113e",
		Body: "body",
		NonceStr: "SrwzbDPKqTOsHkVxzjtTXOnrMrPedtcO",
		DeviceInfo: "info",
	}))

	assert.Equal(t, `316C621C950CC4222D905C463185DDF7`, gwechat.CreatePaySign(gwechat.PaySignSource{
		APPID: "wx9111246221fd5cae",
		MCHID: "1112225871",
		MCHKey: "wxda123461113e2357wxdac12521113e",
		Body: "body",
		NonceStr: "SrwzbDPKqTOsHkVxzjtTXOnrMrPedtcO",
	}))
}