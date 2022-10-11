package noticeex

// INotice 通知接口
type INotice interface {
	Sendf(format string, args ...interface{}) error
}
