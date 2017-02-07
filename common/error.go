package common

import "github.com/pkg/errors"

var (
	ErrInvalidCondition   = Error("Invalid condition")
	ErrInvalidFulfillment = Error("Invalid fulfillment")
	ErrInvalidKey         = Error("Invalid key")
	ErrInvalidRegex       = Error("Invalid regex")
	ErrInvalidSize        = Error("Invalid size")
	ErrInvalidType        = Error("Invalid type")
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Error(msg string) error {
	return errors.New(msg)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	panic(Sprintf(format, args...))
}
