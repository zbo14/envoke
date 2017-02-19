package common

import (
	"os"
)

var Stdout = os.Stdout

type Files []*os.File

func CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}

func MustCreateFile(path string) *os.File {
	file, err := CreateFile(path)
	Check(err)
	return file
}

func OpenWriteFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
}

func MustOpenWriteFile(path string) *os.File {
	file, err := OpenWriteFile(path)
	Check(err)
	return file
}

func OpenFile(path string) (*os.File, error) {
	return os.Open(path)
}

func MustOpenFile(path string) *os.File {
	file, err := OpenFile(path)
	Check(err)
	return file
}

func ReadFile(path string) ([]byte, error) {
	file, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	bytes, err := ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func MustReadFile(path string) []byte {
	bytes, err := ReadFile(path)
	Check(err)
	return bytes
}

func Getenv(key string) string {
	return os.Getenv(key)
}
