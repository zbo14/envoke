package common

import "github.com/pkg/errors"

var (
	ErrExpectedPost       = Error("Expected POST request")
	ErrInvalidCondition   = Error("Invalid condition")
	ErrInvalidFulfillment = Error("Invalid fulfillment")
	ErrInvalidId          = Error("Invalid id")
	ErrInvalidKey         = Error("Invalid key")
	ErrInvalidRegex       = Error("Invalid regex")
	ErrInvalidRequest     = Error("Invalid request")
	ErrInvalidSignature   = Error("Invalid signature")
	ErrInvalidSize        = Error("Invalid size")
	ErrInvalidType        = Error("Invalid type")

	ErrServerReset = Error("Server failed to reset")
	ErrServerStart = Error("Server failed to start")
	ErrServerStop  = Error("Server failed to stop")
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
