package gerr

import "fmt"

type GoError struct {
	code int64
	msg  string
	desc string

	wrapper *GoError
}

// 所有错误码注册到这里
var errors map[int64]string

func init() {
	errors = make(map[int64]string)
	errors[defaultErrorCodeNotFoundCode] = defaultErrorCodeNotFoundMsg
}

// 错误码不存在
var defaultErrorCodeNotFoundCode int64
var defaultErrorCodeNotFoundMsg = "ErrorCodeNotFound"

// 错误码不存在设置
func SetDefaultErrorNotFoundCode(code int64) {
	defaultErrorCodeNotFoundCode = code
}

// 注册错误码
func RegisterError(code int64, format string) {
	if code == defaultErrorCodeNotFoundCode {
		panic(fmt.Sprintf("Code %d conflicts to ErrorCodeNotFound, Use func SetDefaultErrorNotFoundCode()", code))
	}

	errors[code] = format
}

func RegisterErrors(errs map[int64]string) {
	for k, v := range errs {
		RegisterError(k, v)
	}
}

func (e *GoError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var msg = fmt.Sprintf("%d-%s", e.code, e.msg)
	if e.desc != "" {
		msg += "(" + e.desc + ")"
	}

	if e.wrapper == nil {
		return msg
	}

	return msg + " >> " + e.wrapper.Error()
}

func New(code int64, desc ...string) *GoError {
	var e GoError
	var ok bool
	e.code = code
	e.msg, ok = errors[code]
	if !ok {
		e.msg = defaultErrorCodeNotFoundMsg
		e.code = defaultErrorCodeNotFoundCode
		return &e
	}

	if len(desc) == 0 {
		return &e
	}
	e.desc = desc[0]
	return &e
}

func (e *GoError) Wrap(code int64, desc ...string) *GoError {
	if e.wrapper == nil {
		e.wrapper = New(code, desc...)
		return e
	}

	e.wrapper.Wrap(code, desc...)
	return e
}
