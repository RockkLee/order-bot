package apperr

import "errors"

type Err struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func (e Err) Error() string { return e.Msg }

func (e Err) Is(target error) bool {
	var t Err
	ok := errors.As(target, &t)
	return ok && e.Code == t.Code
}
