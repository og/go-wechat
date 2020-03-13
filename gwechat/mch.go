package gwechat



type MCH struct {
	PublicAPPID string
	MCHID string
}
type MCHConfig struct {
	PublicAPPID string
	MCHID string
}
func NewMCH(config MCHConfig) (mch MCH) {
	mch.PublicAPPID = config.PublicAPPID
	mch.MCHID = config.MCHID
	return
}
