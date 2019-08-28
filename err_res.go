package gwechat
import "github.com/og/go-json"
type ErrResponse struct {
	Fail bool  `json:"fail"`
	ErrCode int `json:"errcode"`
	ErrMsg string `json:"errmsg"`
}
func (self ErrResponse) String() string {
	return gjson.String(self)
}