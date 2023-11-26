package slacktestbot

import "errors"

type ErrCode uint16

const (
	ErrCodeNotFoundPath ErrCode = 4004
	ErrCodeUndefined    ErrCode = 5000
)

type ErrEntity struct {
	Code    ErrCode
	Status  int
	Message string
}

var (
	errEntityMapping = map[ErrCode]ErrEntity{
		ErrCodeNotFoundPath: {ErrCodeNotFoundPath, 404, "Not Found Path"},
		ErrCodeUndefined:    {ErrCodeUndefined, 500, "Undefined Error"},
	}
)

func GetErrEntity(code ErrCode) ErrEntity {
	if status, ok := errEntityMapping[code]; ok {
		return status
	}
	return errEntityMapping[ErrCodeUndefined]
}

type ErrWithCode struct {
	Code ErrCode
	Err  error
}

func (e ErrWithCode) Error() string {
	return e.Err.Error()
}

func NewError(code ErrCode, message string) error {
	return ErrWithCode{Code: code, Err: errors.New(message)}
}

func WrapError(code ErrCode, message string, cause error) error {
	return ErrWithCode{Code: code, Err: errors.Join(errors.New(message), cause)}
}
