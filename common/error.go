package common

import "github.com/pkg/errors"

var (
	ErrCriteriaNotMet     = Error("Criteria not met")
	ErrEmptyStr           = Error("Empty string")
	ErrExpectedPost       = Error("Expected POST request")
	ErrExpectedGet        = Error("Expected GET request")
	ErrInvalidCondition   = Error("Invalid condition")
	ErrInvalidEmail       = Error("Invalid email")
	ErrInvalidField       = Error("Invalid field")
	ErrInvalidFingerprint = Error("Invalid fingerprint")
	ErrInvalidFulfillment = Error("Invalid fulfillment")
	ErrInvalidId          = Error("Invalid id")
	ErrInvalidKey         = Error("Invalid key")
	ErrInvalidLogin       = Error("Invalid login")
	ErrInvalidModel       = Error("Invalid model")
	ErrInvalidRequest     = Error("Invalid request")
	ErrInvalidSignature   = Error("Invalid signature")
	ErrInvalidSize        = Error("Invalid size")
	ErrInvalidTerritory   = Error("Invalid territory")
	ErrInvalidTime        = Error("Invalid time")
	ErrInvalidType        = Error("Invalid type")
	ErrInvalidUrl         = Error("Invalid url")
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
