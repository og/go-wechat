package gwechat

// 接口名称：发放红包接口
// 官方文档: https://pay.weixin.qq.com/wiki/doc/api/tools/cash_coupon.php?chapter=13_4&index=3
// 接口地址: https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack
func (mch MCH) Sendredpack() {

}

type SendredpackReq struct {
	SendName string `note:"红包发送者名称"`
	ReOpenID string `note:"接受红包的用户openid"`
	TotalAmount int `note:"付款金额，单位分"`
	Wishing string  `note:"红包祝福语	"`
	ClientIP string `note:"调用接口的机器IP地址"`
	ActName string  `note:"活动名称"`
	Remark string   `note:"备注信息"`
	SceneID string  `note:"发放红包使用场景，红包金额大于200或者小于1元时必传"`
	RiskInfo string `note:"信息"`
}
func (SendredpackReq) Dict() (dict struct {
	PRODUCT_1_Promotion string `note:"商品促销"`
	PRODUCT_2_Draw string `note:"抽奖"`
	PRODUCT_3_VirtualGoodsWinPrizes string `note:"虚拟物品兑奖"`
	PRODUCT_4_InternalWelfare string `note:"企业内部福利"`
	PRODUCT_5_ChannelCommission string `note:"渠道分润"`
	PRODUCT_6_InsuranceCommission string `note:"保险回馈"`
	PRODUCT_7_Lottery string `note:"彩票派奖"`
	PRODUCT_8_TaxLaTombola string `note:"税务刮奖"`
}) {
	dict.PRODUCT_1_Promotion = "PRODUCT_1"
	dict.PRODUCT_2_Draw = "PRODUCT_2"
	dict.PRODUCT_3_VirtualGoodsWinPrizes = "PRODUCT_3"
	dict.PRODUCT_4_InternalWelfare = "PRODUCT_4"
	dict.PRODUCT_5_ChannelCommission = "PRODUCT_5"
	dict.PRODUCT_6_InsuranceCommission = "PRODUCT_6"
	dict.PRODUCT_7_Lottery = "PRODUCT_7"
	dict.PRODUCT_8_TaxLaTombola = "PRODUCT_8"
	return
}
type sendredpackReqXML struct {
	nonce_str string `xml:"nonce_str"`
	sign string `xml:"sign"`
	mch_billno string `xml:"mch_billno"`
	mch_id string `xml:"mch_id"`
	wxappid string `xml:"wxappid"`
	send_name string `xml:"send_name"`
	re_openid string `xml:"re_openid"`
	total_amount int `xml:"total_amount"`
	total_num int `xml:"total_num"`
	wishing string `xml:"wishing"`
	client_ip string `xml:"client_ip"`
	act_name string `xml:"act_name"`
	remark string `xml:"remark"`
	scene_id string `xml:"scene_id"`
	risk_info string `xml:"risk_info"`
}