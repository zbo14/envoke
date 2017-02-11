package common

import "github.com/pkg/errors"

var (
	ErrExpectedPost       = Error("Expected POST request")
	ErrExpectedGet        = Error("Expected GET request")
	ErrInvalidCondition   = Error("Invalid condition")
	ErrInvalidId          = Error("Invalid id")
	ErrInvalidFulfillment = Error("Invalid fulfillment")
	ErrInvalidKey         = Error("Invalid key")
	ErrInvalidModel       = Error("Invalid model")
	ErrInvalidRegex       = Error("Invalid regex")
	ErrInvalidRequest     = Error("Invalid request")
	ErrInvalidSignature   = Error("Invalid signature")
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

func ErrorAppend(err error, msg string) error {
	return Error(err.Error() + ": " + msg)
}
