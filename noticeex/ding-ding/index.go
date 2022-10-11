package dingding

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/xm-chentl/gocore/noticeex"
)

type dingdingImpl struct {
	token     string
	secretKey string
}

func (d dingdingImpl) Sendf(format string, args ...interface{}) (err error) {
	return d.send(fmt.Sprintf(format, args...))
}

func (d *dingdingImpl) genSign() (timestamp int64, sign string) {
	timestamp = time.Now().UnixNano() / 1e6
	str := fmt.Sprintf("%d\n%s", timestamp, d.secretKey)
	h := hmac.New(sha256.New, []byte(d.secretKey))
	h.Write([]byte(str))
	sign = base64.StdEncoding.EncodeToString(h.Sum(nil))

	return
}

func (d dingdingImpl) getURL() string {
	timestamp, sign := d.genSign()
	return fmt.Sprintf(
		"https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%d&sign=%s",
		d.token, timestamp, sign,
	)
}

func (d dingdingImpl) send(msg string) (err error) {
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": msg,
		},
	}

	bytesData, _ := json.Marshal(data)
	resp, err := http.Post(
		d.getURL(),
		"application/json",
		bytes.NewReader(bytesData),
	)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request dingding webhook fail err(code:%d, msg:%v)", resp.StatusCode, err)
		return
	}

	type respResult struct {
		ErrCode int      `json:"errcode"`
		ErrMsg  string   `json:"errmsg"`
		Help    []string `json:"more"`
	}
	result := &respResult{}
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, result)
	if result.ErrCode != 0 {
		err = fmt.Errorf(result.ErrMsg)
		return
	}

	return
}

// New 实例一个消息实现
func New(token, secretKey string) noticeex.INotice {
	return &dingdingImpl{
		token:     token,
		secretKey: secretKey,
	}
}
