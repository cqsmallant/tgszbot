package constant

var Errno = map[int]string{
	400:   "系统错误",
	401:   "签名认证错误",
	10009: "无法解析请求参数",
}

var (
	SystemErr    = Err(400)
	SignatureErr = Err(401)
)

type RspError struct {
	Code int
	Msg  string
}

func (re *RspError) Error() string {
	return re.Msg
}

func Err(code int) (err error) {
	err = &RspError{
		Code: code,
		Msg:  Errno[code],
	}
	return err
}

func (re *RspError) Render() (code int, msg string) {
	return re.Code, re.Msg
}
