package gwechat

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	gjson "github.com/og/go-json"
	"github.com/pkg/errors"
	"log"
	"regexp"
)

type WeappUserInfoDataWatermark struct {
	APPID     string `json:"appid"`
	Timestamp int `json:"timestamp"`
}
type WeappUserInfoData struct {
	AvatarURL string `json:"avatarUrl"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Gender    int `json:"gender"`
	Language  string `json:"language"`
	NickName  string `json:"nickName"`
	Province  string `json:"province"`
	OpenID    string `json:"openId"`
	UnionID   string `json:"unionId"`
	Watermark WeappUserInfoDataWatermark `json:"watermark"`
}
func (self Wechat) DecodeWeappUserInfo(sessionKey string, encryptedData string, iv string) (data WeappUserInfoData, err error) {
	if len(sessionKey) != 24 {
		err = errors.New("sessionKey len must be 24")
		return
	}
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		err = errors.New("sessionKey decodeBase64Error")
		return
	}
	if len(iv) != 24 {
		err = errors.New("iv len must be 24")
		return
	}
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		err = errors.New("iv decodeBase64Error")
		return
	}
	if encryptedData == "" {
		err = errors.New("encryptedData can not be a empty string")
		return
	}
	aesCipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		err = errors.New("encryptedData decodeBase64Error")
		return
	}
	aesPlantText := make([]byte, len(aesCipherText))
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		err = errors.New("IllegalBuffer")
		return
	}

	mode := cipher.NewCBCDecrypter(aesBlock, aesIV)
	mode.CryptBlocks(aesPlantText, aesCipherText)
	aesPlantText = protectPKCS7UnPadding(aesPlantText)
	re := regexp.MustCompile(`[^\{]*(\{.*\})[^\}]*`)
	byteList := []byte(re.ReplaceAllString(string(aesPlantText), "$1"))
	err = gjson.ParseByteWithErr(byteList, &data)
	if err != nil {
		log.Print("DecodeWeappUserInfo json parser error ", string(byteList))
		return
	}
	if data.Watermark.APPID != self.appID {
		err = errors.New("appID is not match")
		return WeappUserInfoData{}, err
	}
	return
}
func protectPKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	if length > 0 {
		unPadding := int(plantText[length-1])
		return plantText[:(length - unPadding)]
	}
	return plantText
}