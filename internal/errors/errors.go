package errors

import "errors"

type wrappedErr struct {
	msg string
	err error
}

func (e wrappedErr) Error() string {
	if e.err == nil {
		return e.msg
	}
	return e.msg + "\n\tcaused by: " + e.err.Error()
}

func (e wrappedErr) Unwrap() error {
	return e.err
}

func Wrap(err error, msg string) error {
	return wrappedErr{
		msg: msg,
		err: err,
	}
}

func New(msg string) error {
	return wrappedErr{msg: msg}
}

func AppErrMsg(err error) (msg string, ok bool) {
	for ; err != nil; err = errors.Unwrap(err) {
		if wrapped, isWrapped := err.(wrappedErr); isWrapped {
			msg, ok = wrapped.msg, true
		}
	}

	return msg, ok
}
