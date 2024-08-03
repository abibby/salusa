package errors

import (
	goerrors "errors"
	"runtime/debug"
)

type SentinelError string

func (e SentinelError) Error() string {
	return string(e)
}

type Stacker interface {
	Stack() []byte
}

type StackError struct {
	err   error
	stack []byte
}

func WithStack(err error) error {
	return &StackError{
		err:   err,
		stack: debug.Stack(),
	}
}
func (e *StackError) Error() string {
	return e.err.Error()
}
func (e *StackError) Stack() []byte {
	return e.stack
}
func (e *StackError) Unwrap() error {
	return e.err
}

type Error struct {
	message string
	stack   []byte
}

func New(message string) error {
	return &Error{
		message: message,
		stack:   debug.Stack(),
	}
}
func (e *Error) Error() string {
	return e.message
}
func (e *Error) Stack() []byte {
	return e.stack
}

func As(err error, target any) bool {
	return goerrors.As(err, target)
}
func Is(err, target error) bool {
	return goerrors.Is(err, target)
}
func Join(errs ...error) error {
	return goerrors.Join(errs...)
}
func Unwrap(err error) error {
	return goerrors.Unwrap(err)
}
