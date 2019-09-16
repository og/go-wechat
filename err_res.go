package gwechat

import (
	"github.com/og/go-json"
	"github.com/pkg/errors"
)
type ErrResponse struct {
	Fail bool  `json:"fail"`
	ErrCode int `json:"errcode"`
	ErrMsg string `json:"errmsg"`
	Error error
}
func (self *ErrResponse) SetFail(code int, msg string) {
	self.Fail = true
	self.ErrCode = code
	self.ErrMsg = msg
	self.Error = errors.New(gjson.String(self))
}
func (self *ErrResponse) SetSystemError(err error) {
	self.SetFail(61450, err.Error())
}