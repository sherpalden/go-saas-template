package app_error

import (
	"fmt"

	"github.com/pkg/errors"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
type AppError struct {
	Code          uint         `json:"code"`
	Message       string       `json:"message"`
	FieldErrors   []FieldError `json:"fieldapp_error"`
	originalError error
}

func (err AppError) Error() string {
	return err.originalError.Error()
}

func New(err error, code uint) *AppError {
	return &AppError{
		Code:          code,
		Message:       err.Error(),
		originalError: err,
	}
}

func (appErr *AppError) Wrap(msg string) *AppError {
	appErr.originalError = errors.Wrap(appErr, msg)
	return appErr
}

func (appErr *AppError) Wrapf(msg string, args ...interface{}) *AppError {
	appErr.originalError = errors.Wrapf(appErr, msg, args...)
	return appErr
}

func (appErr *AppError) SetMessage(msg string) *AppError {
	appErr.Message = msg
	return appErr
}

func (appErr *AppError) SetMessagef(msg string, args ...interface{}) *AppError {
	appErr.Message = fmt.Sprint(msg, args)
	return appErr
}

func (appErr *AppError) SetFieldErrors(errList []FieldError) *AppError {
	appErr.FieldErrors = append(appErr.FieldErrors, errList...)
	return appErr
}

func (appErr *AppError) SetCode(code uint) *AppError {
	appErr.Code = code
	return appErr
}
