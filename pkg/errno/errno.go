package errno

import (
	"errors"
	"fmt"
)

type ErrNo struct {
	ErrorCode int64
	ErrorMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("[%d] %s", e.ErrorCode, e.ErrorMsg)
}

func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{
		ErrorCode: code,
		ErrorMsg:  msg,
	}
}

func Errorf(code int64, template string, args ...interface{}) ErrNo {
	return ErrNo{
		ErrorCode: code,
		ErrorMsg:  fmt.Sprintf(template, args...),
	}
}

// WithMessage will replace default msg to new
func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrorMsg = msg
	return e
}

// WithError will add error msg after Message
func (e ErrNo) WithError(err error) ErrNo {
	e.ErrorMsg = e.ErrorMsg + ", " + err.Error()
	return e
}

// ConvertErr convert error to ErrNo
// in Default user ServiceErrorCode
func ConvertErr(err error) ErrNo {
	if err == nil {
		return Success
	}
	errno := ErrNo{}
	if errors.As(err, &errno) {
		return errno
	}

	s := InternalServiceError
	s.ErrorMsg = err.Error()
	return s
}
