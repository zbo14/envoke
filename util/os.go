package util

import "os"

func CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}

func OpenFile(path string) (*os.File, error) {
	return os.Open(path)
}

func ReadFile(path string) ([]byte, error) {
	file, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	bytes := ReadAll(file)
	return bytes, nil
}
