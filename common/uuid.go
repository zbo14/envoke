package util

import (
	uuid "github.com/nu7hatch/gouuid"
)

func Uuid4() string {
	uuid4, err := uuid.NewV4()
	Check(err)
	return uuid4.String()
}
