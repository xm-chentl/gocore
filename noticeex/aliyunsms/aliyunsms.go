package aliyunsms

import (
	"encoding/json"
	"strings"

	"github.com/GiterLab/aliyun-sms-go-sdk/dysms"
	"github.com/tobyzxj/uuid"
	"github.com/xm-chentl/gocore/noticeex"
)

type smsImpl struct {
	accessID  string
	accessKey string
	signName  string
}

func (s smsImpl) Sendf(format string, args ...interface{}) error {
	return nil
}

func (s smsImpl) Send(templateCode string, templateParam map[string]string, phones ...string) (err error) {
	templateParamByte, _ := json.Marshal(templateParam)
	_, err = dysms.SendSms(
		uuid.New(),
		strings.Join(phones, ","),
		s.signName,
		templateCode,
		string(templateParamByte),
	).DoActionWithException()

	return
}

func New(accessID, accessKey, signName string) noticeex.INotice {
	dysms.HTTPDebugEnable = true // 生成环境可以去掉
	dysms.SetACLClient(accessID, accessKey)
	return &smsImpl{
		accessID:  accessID,
		accessKey: accessKey,
	}
}
