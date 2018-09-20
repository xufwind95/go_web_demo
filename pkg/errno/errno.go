package errno

import "fmt"

type Errno struct {
	Code    int
	Message string
}

// 实现error接口
func (err Errno) Error() string {
	return err.Message
}

type Err struct {
	Code    int
	Message string
	Err     error
}

// 创建Err对象
func New(errno *Errno, err error) *Err {
	return &Err{
		Code:    errno.Code,
		Message: errno.Message,
		Err:     err,
	}
}

// Err对象添加错误信息
func (err *Err) Add(message string) error {
	err.Message += message
	return err
}

// Err对象按格式添加错误信息
func (err *Err) Addf(format string, args ...interface{}) error {
	err.Message += " " + fmt.Sprintf(format, args...)
	return err
}

// 实现error接口
func (err *Err) Error() string {
	return fmt.Sprintf(
		"Err - code: %d, message: %s, error: %s",
		err.Code,    // 错误码
		err.Message, // 错误信息
		err.Err,     // 系统返回的真实错误
	)
}

// 判断是否为某些错误
func IsErrUserNotFound(err error) bool {
	code, _ := DecodeErr(err)
	return code == ErrUserNotFound.Code
}

// 解析错误码和错误信息
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
	}
	return InternalServerError.Code, err.Error()
}
